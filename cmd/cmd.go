package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/redis"
	"github.com/gabe565/geoip-cache-proxy/internal/server"
	"github.com/spf13/cobra"
)

var version = "beta"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "geoip-cache-proxy",
		Short:   "A GeoIP database caching proxy",
		RunE:    run,
		Version: buildVersion(version),

		ValidArgsFunction: cobra.NoFileCompletions,
		SilenceErrors:     true,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
	}
	cmd.InitDefaultVersionFlag()
	conf := config.NewDefault()
	conf.RegisterFlags(cmd)
	config.RegisterCompletions(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}

var ErrMissingConfig = errors.New("command missing config")

func run(cmd *cobra.Command, _ []string) error {
	conf, ok := config.FromContext(cmd.Context())
	if !ok {
		return ErrMissingConfig
	}

	if err := conf.Load(cmd); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cache, err := redis.Connect(conf)
	if err != nil {
		return err
	}
	defer cache.Close()

	return server.ListenAndServe(ctx, conf, cache)
}
