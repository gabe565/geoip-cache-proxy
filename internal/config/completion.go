package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc(FlagLogLevel,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{
				zerolog.LevelTraceValue,
				zerolog.LevelDebugValue,
				zerolog.LevelInfoValue,
				zerolog.LevelWarnValue,
				zerolog.LevelErrorValue,
				zerolog.LevelFatalValue,
				zerolog.LevelPanicValue,
			}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	if err := cmd.RegisterFlagCompletionFunc(FlagLogFormat,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{FormatAuto, FormatColor, FormatPlain, FormatJSON}, cobra.ShellCompDirectiveNoFileComp
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
