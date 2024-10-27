package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"gabe565.com/geoip-cache-proxy/internal/redis"
)

func Live() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusNoContent)
	}
}

func Ready(cache *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		if err := cache.Ping(r.Context()); err != nil {
			err = fmt.Errorf("failed to connect to redis: %w", err)
			slog.Error("Readiness check failed", "error", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
