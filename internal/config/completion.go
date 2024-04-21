package config

import "github.com/spf13/cobra"

func RegisterCompletions(cmd *cobra.Command) {
	if err := cmd.RegisterFlagCompletionFunc(FlagLogLevel,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	if err := cmd.RegisterFlagCompletionFunc(FlagLogFormat,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"auto", "color", "plain", "json"}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	npCompFlags := []string{
		FlagRedisAddr,
		FlagRedisPassword,
		FlagRedisDB,
		FlagUpdatesAddr,
		FlagUpdatesHost,
		FlagDownloadAddr,
		FlagDownloadHost,
		FlagCacheDuration,
		FlagDebugAddr,
	}
	for _, name := range npCompFlags {
		if err := cmd.RegisterFlagCompletionFunc(name, cobra.NoFileCompletions); err != nil {
			panic(err)
		}
	}
}
