package service

import (
	"encoding/base64"
	"flag"
	viper "github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"github.com/allegro/bigcache"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db *gorm.DB
	localCache *bigcache.BigCache
	listenPath string
)

// StartServing is the main entry of the application
func StartServing() {
	var decryptKey string
	var configLocation string
	var dataToEncrypt string
	flag.StringVar(&decryptKey, "decryptKey", "", "DecryptKey")
	flag.StringVar(&configLocation, "configLocation", "", "ConfigLocation")
	flag.StringVar(&dataToEncrypt, "dataToEncrypt", "", "DataToEncrypt")

	flag.Parse()
	if decryptKey == "" {
		log.Fatalf("decryptKey is not set")
	}
	key := []byte(decryptKey)

	if dataToEncrypt != "" {
		resultData, _ := AesEncrypt([]byte(dataToEncrypt), key)
		log.Printf("Encrypted data %s", base64.StdEncoding.EncodeToString(resultData))
		return
	}

	viper.AddConfigPath(configLocation)
	viper.SetConfigName("podistributor-config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Fatal error config file", err)
	}

	viper.SetDefault("maxDbIdleConns", 16)
	viper.SetDefault("maxDbOpenConns", 16)
	viper.SetDefault("listenPort", 17800)
	viper.SetDefault("listenPath", "/")
	viper.SetDefault("monitorPort", 11800)
	maxDbIdleConns := viper.GetInt("maxDbIdleConns")
	maxDbOpenConns := viper.GetInt("maxDbOpenConns")
	listenPort := viper.GetInt("listenPort")
	monitorPort := viper.GetInt("monitorPort")
	listenPath = viper.GetString("listenPath")

	log.Printf("decryptKey:%s, maxDbIdleConns:%d, maxDbOpenConns:%d, configLocation:%s, listenPort:%d, listenPath:%s, monitorPort:%d", 
		decryptKey, maxDbIdleConns, maxDbOpenConns, configLocation, listenPort, listenPath, monitorPort)

	encryptedDbInfo := viper.GetString("encryptedDbInfo")
	if encryptedDbInfo == "" {
		log.Fatalln("encryptedDbInfo is not found")
	}

	decodedEncryptedDbInfo, err := base64.StdEncoding.DecodeString(encryptedDbInfo)
	if err != nil {
		log.Fatalln("Failed to decode base64 db info", err)
	}
	originDbInfo, err := AesDecrypt(decodedEncryptedDbInfo, key)
	if err != nil {
		log.Fatalln("Failed to decript dbInfo", err)
	}
	dbInfo := string(originDbInfo)

	db, err = gorm.Open(mysql.Open(dbInfo), &gorm.Config{})
	sqlDB, err := db.DB()
	defer sqlDB.Close()
	if err != nil {
		log.Panicln("Failed to create db connection", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Panicln("Failed to connect to db", err)
	}
	sqlDB.SetMaxIdleConns(maxDbIdleConns)
	sqlDB.SetMaxOpenConns(maxDbOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Printf("Connected to MySQL.")

	config := bigcache.Config {
		Shards: 256,
		LifeWindow: 60 * time.Minute,
		CleanWindow: 3 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize: 500,
		Verbose: false,
		HardMaxCacheSize: 512,
		OnRemove: nil,
		OnRemoveWithReason: nil,
	}
	localCache, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal(err)
	}

	go func () {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":" + strconv.Itoa(monitorPort), nil))
	}()

	http.HandleFunc(listenPath, HandleRequest)
	http.ListenAndServe(":" + strconv.Itoa(listenPort), nil)
}
