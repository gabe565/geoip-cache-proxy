package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc(FlagLogLevel,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"trace",
				strings.ToLower(slog.LevelDebug.String()),
				strings.ToLower(slog.LevelInfo.String()),
				strings.ToLower(slog.LevelWarn.String()),
				strings.ToLower(slog.LevelError.String()),
			}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	if err := cmd.RegisterFlagCompletionFunc(FlagLogFormat,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return LogFormatStrings(), cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	npCompFlags := []string{
		FlagRedisHost,
		FlagRedisPort,
		FlagRedisPassword,
		FlagRedisDB,
		FlagUpdatesAddr,
		FlagUpdatesHost,
		FlagDownloadAddr,
		FlagDownloadHost,
		FlagAccountID,
		FlagLicenseKey,
		FlagCacheDuration,
		FlagDebugAddr,
	}
	for _, name := range npCompFlags {
		if err := cmd.RegisterFlagCompletionFunc(name, cobra.NoFileCompletions); err != nil {
			panic(err)
		}
	}
}
