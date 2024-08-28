package config

import (
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel  string
	LogFormat string

	RedisHost     string
	RedisPort     uint16
	RedisPassword string
	RedisDB       int

	HTTPTimeout  time.Duration
	UpdatesAddr  string
	UpdatesHost  string
	DownloadAddr string
	DownloadHost string

	AccountID  int
	LicenseKey string

	CacheDuration time.Duration

	DebugAddr string

	TranslateIngressNginxPaths bool
}

func NewDefault() *Config {
	return &Config{
		LogLevel:  zerolog.LevelInfoValue,
		LogFormat: FormatAuto,

		RedisHost: "localhost",
		RedisPort: 6379,

		HTTPTimeout:  30 * time.Second,
		UpdatesAddr:  ":8080",
		UpdatesHost:  "updates.maxmind.com",
		DownloadAddr: ":8081",
		DownloadHost: "download.maxmind.com",

		CacheDuration: 12 * time.Hour,

		DebugAddr: ":6060",
	}
}
