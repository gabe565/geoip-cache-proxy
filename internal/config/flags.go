package config

import (
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagHTTPTimeout  = "http-timeout"
	FlagUpdatesAddr  = "updates-addr"
	FlagUpdatesHost  = "updates-host"
	FlagDownloadAddr = "download-addr"
	FlagDownloadHost = "download-host"

	FlagAccountID  = "account-id"
	FlagLicenseKey = "license-key"

	FlagCacheDir      = "cache-dir"
	FlagCacheDuration = "cache-duration"
	FlagCleanupEvery  = "cleanup-every"

	FlagDebugAddr = "debug-addr"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&c.LogLevel, FlagLogLevel, "l", c.LogLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	cmd.Flags().StringVar(&c.LogFormat, FlagLogFormat, c.LogFormat, "Log format (auto, color, plain, json)")

	cmd.Flags().DurationVar(&c.HTTPTimeout, FlagHTTPTimeout, c.HTTPTimeout, "HTTP request timeout")
	cmd.Flags().StringVar(&c.UpdatesAddr, FlagUpdatesAddr, c.UpdatesAddr, "Listen address")
	cmd.Flags().StringVar(&c.UpdatesHost, FlagUpdatesHost, c.UpdatesHost, "MaxMind updates host")
	cmd.Flags().StringVar(&c.DownloadAddr, FlagDownloadAddr, c.DownloadAddr, "Listen address")
	cmd.Flags().StringVar(&c.DownloadHost, FlagDownloadHost, c.DownloadHost, "MaxMind download host")

	cmd.Flags().IntVar(&c.AccountID, FlagAccountID, c.AccountID, "MaxMind account ID")
	cmd.Flags().StringVar(&c.LicenseKey, FlagLicenseKey, c.LicenseKey, "MaxMind license key")

	cmd.Flags().StringVar(&c.CacheDir, FlagCacheDir, c.CacheDir, "Cache directory path")
	cmd.Flags().DurationVar(&c.CacheDuration, FlagCacheDuration, c.CacheDuration, "Length of time to cache MaxMind response")
	cmd.Flags().DurationVar(&c.CleanupEvery, FlagCleanupEvery, c.CleanupEvery, "Interval to clean up expired cache entries")

	cmd.Flags().StringVar(&c.DebugAddr, FlagDebugAddr, c.DebugAddr, "Debug pprof listen address")
}
