package config

import (
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagRedisHost     = "redis-host"
	FlagRedisPort     = "redis-port"
	FlagRedisPassword = "redis-password"
	FlagRedisDB       = "redis-db"

	FlagHTTPTimeout  = "http-timeout"
	FlagUpdatesAddr  = "updates-addr"
	FlagUpdatesHost  = "updates-host"
	FlagDownloadAddr = "download-addr"
	FlagDownloadHost = "download-host"

	FlagAccountID  = "account-id"
	FlagLicenseKey = "license-key"

	FlagCacheDuration = "cache-duration"

	FlagDebugAddr = "debug-addr"

	FlagTranslateIngressNginxUrls = "translate-ingress-nginx-urls"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	cmd.Flags().StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (auto, color, plain, json)")

	cmd.Flags().StringVar(&c.RedisHost, FlagRedisHost, c.RedisHost, "Redis host")
	cmd.Flags().Uint16Var(&c.RedisPort, FlagRedisPort, c.RedisPort, "Redis port")
	cmd.Flags().StringVar(&c.RedisPassword, FlagRedisPassword, c.RedisPassword, "Redis password")
	cmd.Flags().IntVar(&c.RedisDB, FlagRedisDB, c.RedisDB, "Redis database")

	cmd.Flags().DurationVar(&c.HTTPTimeout, FlagHTTPTimeout, c.HTTPTimeout, "HTTP request timeout")
	cmd.Flags().StringVar(&c.UpdatesAddr, FlagUpdatesAddr, c.UpdatesAddr, "Listen address")
	cmd.Flags().StringVar(&c.UpdatesHost, FlagUpdatesHost, c.UpdatesHost, "MaxMind updates host")
	cmd.Flags().StringVar(&c.DownloadAddr, FlagDownloadAddr, c.DownloadAddr, "Listen address")
	cmd.Flags().StringVar(&c.DownloadHost, FlagDownloadHost, c.DownloadHost, "MaxMind download host")

	cmd.Flags().IntVar(&c.AccountID, FlagAccountID, c.AccountID, "MaxMind account ID")
	cmd.Flags().StringVar(&c.LicenseKey, FlagLicenseKey, c.LicenseKey, "MaxMind license key")

	cmd.Flags().DurationVar(&c.CacheDuration, FlagCacheDuration, c.CacheDuration, "Length of time to cache MaxMind response")

	cmd.Flags().StringVar(&c.DebugAddr, FlagDebugAddr, c.DebugAddr, "Debug pprof listen address")

	cmd.Flags().BoolVarP(&c.TranslateIngressNginxPaths, FlagTranslateIngressNginxUrls, "", false, "Whether to translate Ingress-Nginx's expected file names to Maxmind paths.")
}
