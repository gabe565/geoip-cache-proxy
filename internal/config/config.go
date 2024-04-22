package config

import (
	"time"
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

	CacheDuration time.Duration

	DebugAddr string
}

func NewDefault() *Config {
	return &Config{
		LogLevel:  "info",
		LogFormat: "auto",

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
