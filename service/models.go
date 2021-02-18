package service

import (
	"time"
)

// Podcast is the struct of the corresponding db table `podcast` model.
type Podcast struct {
	ID uint `gorm:"primarykey"`
	Name string
	Enabled bool
	Rss string
	Domain string
	EpisodeURLDomain string
	EpisodeMainURILevel int
	EpisodeBakcupURLEnabled bool
	EpisodeBackupURLLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Episode is the struct of the corresponding db table `episode` model.
type Episode struct {
	ID uint `gorm:"primarykey"`
	PodcastID int64
	MainURIList string
	BackupURLList string
	AnalysisURLList string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CacheModel is a struct to combine other models in one cache model.
type CacheModel struct {
	Exists bool
	PodcastModel Podcast
	EpisodeModel Episode
}