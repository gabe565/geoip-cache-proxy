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

	UpdatesAddr string
	UpdatesHost string

	DownloadAddr string
	DownloadHost string

	CacheDuration time.Duration

	DebugAddr string
}
