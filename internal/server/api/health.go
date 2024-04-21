package api

import (
	"fmt"
	"net/http"

	"github.com/gabe565/geoip-cache-proxy/internal/cache"
	"github.com/rs/zerolog/log"
)

func Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusNoContent)
}

func Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if err := cache.Client.Ping(r.Context()).Err(); err != nil {
		err = fmt.Errorf("failed to connect to redis: %w", err)
		log.Err(err).Msg("readiness check failed")
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
