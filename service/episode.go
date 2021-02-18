package service

import (
	"log"
	"net/http"
	"strings"
	"gorm.io/gorm"
	"encoding/json"
)


// HandleRequest is the main function which will handle redirect and analysis.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	requestCounter.Inc()
	path := r.URL.Path[len(listenPath):]
    if path == "" {
		http.Error(w, "Empty path", http.StatusBadRequest)
		badRequestCounter.Inc()
		return
    }

    pathArray := strings.Split(path, "/")
	if len(pathArray) != 3 || strings.ToLower(pathArray[1]) != "ep" {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		badRequestCounter.Inc()
		return
	}

    podcastName := pathArray[0]
	episodeName := pathArray[2]
	if podcastName == "" || episodeName == "" {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		badRequestCounter.Inc()
		return
	}

	cacheKey := "pe:" + podcastName + "-" + episodeName
	cacheBytes, _ := localCache.Get(cacheKey)
	if cacheBytes != nil && len(cacheBytes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		notFoundCounter.Inc()
		return
	}
	var dataFromCache CacheModel
    _ = json.Unmarshal(cacheBytes, &dataFromCache)

	dataFromSourceCounter := 0

	var relatedPodcast Podcast
	if dataFromCache.Exists {
		relatedPodcast = dataFromCache.PodcastModel
	} else {
		relatedPodcast = Podcast{}
		podcastQueryResult := db.Where("name = ?", podcastName).First(&relatedPodcast)
		dbReqCounter.Inc()
		if podcastQueryResult.Error != nil {
			if podcastQueryResult.Error == gorm.ErrRecordNotFound {
				markEmptyCache(cacheKey)
				http.Error(w, "No related podcast found", http.StatusNotFound)
				w.WriteHeader(http.StatusNotFound)
				dbNotFoundCounter.Inc()
				notFoundCounter.Inc()
				return
			}
			log.Printf("Failed to query db for podcast with param podName=%s, episode=%s, err=%s", podcastName, episodeName, podcastQueryResult.Error)
			http.Error(w, "No related podcast found", http.StatusNotFound)
			internalErrorCounter.Inc()
			return
		}
		dataFromSourceCounter++
	}

	var relatedEpisode Episode
	if dataFromCache.Exists {
		relatedEpisode = dataFromCache.EpisodeModel
	} else {
		relatedEpisode = Episode{}
		episodeQueryResult := db.Where("podcast_id = ? AND name = ?", relatedPodcast.ID, episodeName).First(&relatedEpisode)
		dbReqCounter.Inc()
		if episodeQueryResult.Error != nil {
			if episodeQueryResult.Error == gorm.ErrRecordNotFound {
				markEmptyCache(cacheKey)
				http.Error(w, "No related episode found", http.StatusNotFound)
				w.WriteHeader(http.StatusNotFound)
				dbNotFoundCounter.Inc()
				notFoundCounter.Inc()
				return
			}
			log.Printf("Failed to query db for episode with param podName=%s, episode=%s, err=%s", podcastName, episodeName, episodeQueryResult.Error)
			http.Error(w, "Failed to query db for episode", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			internalErrorCounter.Inc()
			return
		}
		dataFromSourceCounter++
	}

	if dataFromSourceCounter > 0 {
		dataToPutIntoCache := CacheModel{}
		dataToPutIntoCache.Exists = true
		dataToPutIntoCache.PodcastModel = relatedPodcast
		dataToPutIntoCache.EpisodeModel = relatedEpisode
		cacheBytes, jsonErr := json.Marshal(dataToPutIntoCache)
		if jsonErr != nil {
			log.Printf("Failed to convert cache model to json, model=%+v, err=%s", dataToPutIntoCache, jsonErr)
			cachePutFailureCounter.Inc()
		} else {
			cachePutCounter.Inc()
			localCacheErr := localCache.Set(cacheKey, cacheBytes)
			if localCacheErr != nil {
				log.Printf("Failed to save cache model, data=%s, err=%s", string(cacheBytes), localCacheErr)
				cachePutFailureCounter.Inc()
			}
		}
	}

	if !relatedPodcast.Enabled {
		log.Printf("Specified podcast is not enabled podName=%s, episode=%s", podcastName, episodeName)
		http.Error(w, "Specified podcast is not enabled", http.StatusForbidden)
		forbiddenErrorCounter.Inc()
		return
	}

	var mainEpisodeURIArr []string
    _ = json.Unmarshal([]byte(relatedEpisode.MainURIList), &mainEpisodeURIArr)
	var backupEpisodeURLArr []string
    _ = json.Unmarshal([]byte(relatedEpisode.BackupURLList), &backupEpisodeURLArr)

	if len(mainEpisodeURIArr) == 0 && len(backupEpisodeURLArr) == 0 {
		log.Printf("No episode url configured, podName=%s, episode=%s", podcastName, episodeName)
		http.Error(w, "No episode url configured", http.StatusInternalServerError)
		internalErrorCounter.Inc()
		return
	}

	var redirectURI string
	if relatedPodcast.EpisodeBakcupURLEnabled && len(backupEpisodeURLArr) != 0 {
		if len(backupEpisodeURLArr) - 1 < relatedPodcast.EpisodeBackupURLLevel{
			redirectURI = backupEpisodeURLArr[0]
		} else {
			redirectURI = backupEpisodeURLArr[relatedPodcast.EpisodeBackupURLLevel]
		}
	} else {
		if len(mainEpisodeURIArr) - 1 < relatedPodcast.EpisodeMainURILevel {
			redirectURI = relatedPodcast.Domain + mainEpisodeURIArr[0]
		} else {
			redirectURI = relatedPodcast.Domain + mainEpisodeURIArr[relatedPodcast.EpisodeMainURILevel]
		}
	}

	if redirectURI == "" {
		log.Printf("RedirectURL is blank, podName=%s, episode=%s", podcastName, episodeName)
		http.Error(w, "RedirectURL is unexpected blank", http.StatusInternalServerError)
		internalErrorCounter.Inc()
		return
	}

	var analysisURLList []string
    _ = json.Unmarshal([]byte(relatedEpisode.AnalysisURLList), &analysisURLList)
	asyncSendAnalysisData(r, analysisURLList)

	succeedCounter.Inc()
	http.Redirect(w, r, redirectURI, http.StatusFound)
}

func markEmptyCache(key string) {
	localCache.Set(key, []byte{})
	emptyCachePutCounter.Inc()
}

// asyncSendAnalysisData is invoked by HandleRequest in order to request analysis services async.
func asyncSendAnalysisData(req *http.Request, analysisURLArr []string) {
	if len(analysisURLArr) == 0 {
		return
	}
	analysisHitCounter.Inc()
	httpClient := &http.Client{}
	for _, url := range analysisURLArr {
		currentURL := url
		go func() {
			proxyReq, err := http.NewRequest(http.MethodGet, currentURL, nil)

			proxyReq.Header = make(http.Header)
			for h, val := range req.Header {
				proxyReq.Header[h] = val
			}

			resp, err := httpClient.Do(proxyReq)
			if err != nil {
				log.Printf("Failed to async send analysis request, url=%s", currentURL)
				analysisFailureCounter.Inc()
				return
			}
			defer resp.Body.Close()
		}()
	}
}