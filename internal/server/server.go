package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/redis"
	"github.com/gabe565/geoip-cache-proxy/internal/server/api"
	"github.com/gabe565/geoip-cache-proxy/internal/server/handlers"
	"github.com/gabe565/geoip-cache-proxy/internal/server/handlers/proxy"
	geoipmiddleware "github.com/gabe565/geoip-cache-proxy/internal/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(ctx context.Context, conf *config.Config, cache *redis.Client) error {
	group, ctx := errgroup.WithContext(ctx)

	download := NewDownload(conf, cache)
	group.Go(func() error {
		slog.Info("Starting download server", "address", conf.DownloadAddr)
		return download.ListenAndServe()
	})

	updates := NewUpdates(conf, cache)
	group.Go(func() error {
		slog.Info("Starting updates server", "address", conf.UpdatesAddr)
		return updates.ListenAndServe()
	})

	debug := NewDebug(conf, cache)
	if debug != nil {
		group.Go(func() error {
			slog.Info("Starting debug pprof server", "address", conf.DebugAddr)
			return debug.ListenAndServe()
		})
	}

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), conf.HTTPTimeout)
	defer shutdownCancel()

	group.Go(func() error {
		slog.Info("Stopping download server")
		return download.Shutdown(shutdownCtx)
	})
	group.Go(func() error {
		slog.Info("Stopping updates server")
		return updates.Shutdown(shutdownCtx)
	})
	group.Go(func() error {
		slog.Info("Stopping debug pprof server")
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
		r := newMux(conf)
		r.Get("/livez", api.Live())
		r.Get("/readyz", api.Ready(cache))
		r.Mount("/debug", middleware.Profiler())
		return &http.Server{
			Addr:           conf.DebugAddr,
			Handler:        r,
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
	r.Get("/robots.txt", handlers.RobotsTxt)
	return r
}
