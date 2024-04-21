package server

import (
	"net/http"
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

	return group.Wait()
}

func NewDownload(conf *config.Config) *http.Server {
	return &http.Server{
		Addr:              conf.DownloadAddr,
		Handler:           middleware.Log(proxy.Proxy(conf.DownloadHost)),
		ReadHeaderTimeout: 3 * time.Second,
	}
}

func NewUpdates(conf *config.Config) *http.Server {
	return &http.Server{
		Addr:              conf.UpdatesAddr,
		Handler:           middleware.Log(proxy.Proxy(conf.UpdatesHost)),
		ReadHeaderTimeout: 3 * time.Second,
	}
}
