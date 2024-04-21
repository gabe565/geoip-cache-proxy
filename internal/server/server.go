package server

import (
	"context"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/server/api"
	"github.com/gabe565/geoip-cache-proxy/internal/server/middleware"
	"github.com/gabe565/geoip-cache-proxy/internal/server/proxy"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, conf *config.Config) error {
	group, ctx := errgroup.WithContext(ctx)

	download := NewDownload(conf)
	group.Go(func() error {
		log.Info().Str("address", conf.DownloadAddr).Msg("starting download server")
		return download.ListenAndServe()
	})

	updates := NewUpdates(conf)
	group.Go(func() error {
		log.Info().Str("address", conf.UpdatesAddr).Msg("starting updates server")
		return updates.ListenAndServe()
	})

	debug := NewDebug(conf)
	if debug != nil {
		group.Go(func() error {
			log.Info().Str("address", conf.DebugAddr).Msg("starting debug pprof server")
			return debug.ListenAndServe()
		})
	}

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	group.Go(func() error {
		log.Info().Msg("stopping download server")
		return download.Shutdown(shutdownCtx)
	})
	group.Go(func() error {
		log.Info().Msg("stopping updates server")
		return updates.Shutdown(shutdownCtx)
	})
	group.Go(func() error {
		log.Info().Msg("stopping debug pprof server")
		return debug.Shutdown(shutdownCtx)
	})

	return group.Wait()
}

func NewDownload(conf *config.Config) *http.Server {
	return &http.Server{
		Addr:              conf.DownloadAddr,
		Handler:           middleware.Log(proxy.Proxy(conf, conf.DownloadHost)),
		ReadHeaderTimeout: 3 * time.Second,
	}
}

func NewUpdates(conf *config.Config) *http.Server {
	return &http.Server{
		Addr:              conf.UpdatesAddr,
		Handler:           middleware.Log(proxy.Proxy(conf, conf.UpdatesHost)),
		ReadHeaderTimeout: 3 * time.Second,
	}
}

func NewDebug(conf *config.Config) *http.Server {
	if conf.DebugAddr != "" {
		http.HandleFunc("/livez", api.Live)
		http.HandleFunc("/readyz", api.Ready)
		return &http.Server{
			Addr:              conf.DebugAddr,
			ReadHeaderTimeout: 3 * time.Second,
		}
	}
	return nil
}
