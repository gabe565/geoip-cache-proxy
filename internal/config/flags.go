package config

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	EnvPrefix = "GEOIP_"

	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagRedisAddr     = "redis-addr"
	FlagRedisPassword = "redis-password"

	FlagUpdatesAddr = "updates-addr"
	FlagUpdatesHost = "updates-host"

	FlagDownloadAddr = "download-addr"
	FlagDownloadHost = "download-host"

	FlagCacheDuration = "cache-duration"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.LogLevel, FlagLogLevel, "l", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	cmd.Flags().StringVar(&c.LogFormat, FlagLogFormat, "color", "Log format (auto, color, plain, json)")

	cmd.Flags().StringVar(&c.RedisAddr, FlagRedisAddr, "localhost:6379", "Redis address")
	cmd.Flags().StringVar(&c.RedisPassword, FlagRedisPassword, "", "Redis password")

	cmd.Flags().StringVar(&c.UpdatesAddr, FlagUpdatesAddr, ":8080", "Listen address")
	cmd.Flags().StringVar(&c.UpdatesHost, FlagUpdatesHost, "updates.maxmind.com", "MaxMind updates host")

	cmd.Flags().StringVar(&c.DownloadAddr, FlagDownloadAddr, ":8081", "Listen address")
	cmd.Flags().StringVar(&c.DownloadHost, FlagDownloadHost, "download.maxmind.com", "MaxMind download host")

	cmd.Flags().DurationVar(&c.CacheDuration, FlagCacheDuration, 24*time.Hour, "Length of time to cache MaxMind response")
}