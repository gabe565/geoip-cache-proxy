package config

import (
	"log/slog"
	"strings"
	"time"
)

type Config struct {
	logLevel  string
	logFormat string

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
		logLevel:  strings.ToLower(slog.LevelInfo.String()),
		logFormat: FormatAuto.String(),

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
