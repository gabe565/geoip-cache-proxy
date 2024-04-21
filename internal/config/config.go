package config

import (
	"time"
)

type Config struct {
	LogLevel  string
	LogFormat string

	RedisAddr     string
	RedisPassword string

	UpdatesAddr string
	UpdatesHost string

	DownloadAddr string
	DownloadHost string

	CacheDuration time.Duration
}
