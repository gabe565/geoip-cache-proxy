package config

import (
	"time"
)

type Config struct {
	LogLevel  string
	LogFormat string

	HTTPTimeout  time.Duration
	UpdatesAddr  string
	UpdatesHost  string
	DownloadAddr string
	DownloadHost string

	AccountID  int
	LicenseKey string

	CacheDir      string
	CacheDuration time.Duration
	CleanupEvery  time.Duration

	DebugAddr string
}

func NewDefault() *Config {
	return &Config{
		LogLevel:  "info",
		LogFormat: "auto",

		HTTPTimeout:  30 * time.Second,
		UpdatesAddr:  ":8080",
		UpdatesHost:  "updates.maxmind.com",
		DownloadAddr: ":8081",
		DownloadHost: "download.maxmind.com",

		CacheDir:      "data",
		CacheDuration: 12 * time.Hour,
		CleanupEvery:  15 * time.Minute,

		DebugAddr: ":6060",
	}
}
