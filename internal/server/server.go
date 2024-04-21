package server

import (
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/server/middleware"
	"github.com/gabe565/geoip-cache-proxy/internal/server/proxy"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func ListenAndServe(conf *config.Config) error {
	var group errgroup.Group

	group.Go(func() error {
		server := NewDownload(conf)
		log.Info().Str("address", conf.DownloadAddr).Msg("starting download server")
		return server.ListenAndServe()
	})

	group.Go(func() error {
		server := NewUpdates(conf)
		log.Info().Str("address", conf.UpdatesAddr).Msg("starting updates server")
		return server.ListenAndServe()
	})

	group.Go(func() error {
		if server := NewDebug(conf); server != nil {
			log.Info().Str("address", conf.DebugAddr).Msg("starting debug pprof server")
			return server.ListenAndServe()
		}
		return nil
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
		return &http.Server{
			Addr:              conf.DebugAddr,
			ReadHeaderTimeout: 3 * time.Second,
		}
	}
	return nil
}
