package config

import (
	"log/slog"
	"strings"

	"gabe565.com/utils/must"
	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogLevel,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"trace",
				strings.ToLower(slog.LevelDebug.String()),
				strings.ToLower(slog.LevelInfo.String()),
				strings.ToLower(slog.LevelWarn.String()),
				strings.ToLower(slog.LevelError.String()),
			}, cobra.ShellCompDirectiveNoFileComp
		},
	))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogFormat,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return slogx.FormatStrings(), cobra.ShellCompDirectiveNoFileComp
		},
	))

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
		must.Must(cmd.RegisterFlagCompletionFunc(name, cobra.NoFileCompletions))
	}
}
