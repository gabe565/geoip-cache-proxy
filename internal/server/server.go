package server

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/redis"
	"github.com/gabe565/geoip-cache-proxy/internal/server/api"
	geoipmiddleware "github.com/gabe565/geoip-cache-proxy/internal/server/middleware"
	"github.com/gabe565/geoip-cache-proxy/internal/server/proxy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, conf *config.Config, cache *redis.Client) error {
	group, ctx := errgroup.WithContext(ctx)

	download := NewDownload(conf, cache)
	group.Go(func() error {
		log.Info().Str("address", conf.DownloadAddr).Msg("starting download server")
		return download.ListenAndServe()
	})

	updates := NewUpdates(conf, cache)
	group.Go(func() error {
		log.Info().Str("address", conf.UpdatesAddr).Msg("starting updates server")
		return updates.ListenAndServe()
	})

	debug := NewDebug(conf, cache)
	if debug != nil {
		group.Go(func() error {
			log.Info().Str("address", conf.DebugAddr).Msg("starting debug pprof server")
			return debug.ListenAndServe()
		})
	}

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), conf.HTTPTimeout)
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

	if err := group.Wait(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func NewDownload(conf *config.Config, cache *redis.Client) *http.Server {
	r := newMux(conf)
	r.Get("/*", proxy.Proxy(conf, cache, conf.DownloadHost))

	return &http.Server{
		Addr:           conf.DownloadAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 1024, // 1MiB
	}
}

func NewUpdates(conf *config.Config, cache *redis.Client) *http.Server {
	r := newMux(conf)
	r.Get("/*", proxy.Proxy(conf, cache, conf.UpdatesHost))

	return &http.Server{
		Addr:           conf.UpdatesAddr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 1024, // 1MiB
	}
}

func NewDebug(conf *config.Config, cache *redis.Client) *http.Server {
	if conf.DebugAddr != "" {
		http.Handle("/livez", api.Live())
		http.Handle("/readyz", http.TimeoutHandler(api.Ready(cache), conf.HTTPTimeout, "timeout exceeded"))
		return &http.Server{
			Addr:           conf.DebugAddr,
			ReadTimeout:    10 * time.Second,
			MaxHeaderBytes: 1024 * 1024, // 1MiB
		}
	}
	return nil
}

func newMux(conf *config.Config) *chi.Mux {
	r := chi.NewMux()
	r.Use(middleware.RealIP)
	r.Use(geoipmiddleware.Log)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(conf.HTTPTimeout))
	return r
}
