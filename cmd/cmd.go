package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "geoip-cache-proxy",
		RunE: run,

		SilenceErrors: true,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	return nil
}
