package proxy

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"gabe565.com/geoip-cache-proxy/internal/config"
	"gabe565.com/geoip-cache-proxy/internal/redis"
	"gabe565.com/geoip-cache-proxy/internal/server/consts"
	geoipmiddleware "gabe565.com/geoip-cache-proxy/internal/server/middleware"
	"gabe565.com/utils/slogx"
	"github.com/go-chi/chi/v5/middleware"
)

func Proxy(conf *config.Config, cache *redis.Client, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := upstreamURL(host, r, conf.TranslateIngressNginxPaths)
		logger, ok := geoipmiddleware.LogFromContext(r.Context())
		if !ok {
			logger = slog.Default()
		}
		logger = logger.With("upstreamURL", u.String())

		upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, u.String(), r.Body)
		if err != nil {
			logger.Error("Failed to create request", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upstreamReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		for k := range r.Header {
			upstreamReq.Header.Set(k, r.Header.Get(k))
		}
		setAuth(conf, upstreamReq)

		var cacheStatus CacheStatus
		var upstreamResp *http.Response
		defer func() {
			if upstreamResp != nil {
				_, _ = io.Copy(io.Discard, upstreamResp.Body)
				_ = upstreamResp.Body.Close()
			}
		}()

		if upstreamResp, err = cache.Get(r.Context(), upstreamReq, conf.HTTPTimeout); err == nil {
			slogx.LoggerTrace(logger, "Using cached response")
			cacheStatus = CacheHit
		} else if errors.Is(err, redis.ErrNotExist) {
			slogx.LoggerTrace(logger, "Forwarding request to upstream")
			upstreamResp, err = http.DefaultClient.Do(upstreamReq)
			if err != nil {
				logger.Error("Failed to forward to upstream", "error", err)
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}

			if upstreamResp.StatusCode < 300 {
				if cacheWriter, err := cache.NewWriter(r.Context(), upstreamReq, upstreamResp, conf.CacheDuration); err == nil {
					defer func() {
						if err := cacheWriter.Close(); err != nil {
							logger.Error("Failed to close cache", "error", err)
						}
					}()

					wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
					wrapped.Tee(cacheWriter)
					w = wrapped
				} else {
					logger.Error("Failed to cache response", "error", err)
				}
			} else {
				cacheStatus = CacheBypass
			}
		} else {
			logger.Warn("Failed to get cached response", "error", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		for k := range upstreamResp.Header {
			w.Header().Set(k, upstreamResp.Header.Get(k))
		}
		w.Header().Set(consts.UpstreamURLHeader, u.String())
		w.Header().Set(consts.CacheStatusHeader, cacheStatus.String())
		w.WriteHeader(upstreamResp.StatusCode)
		_, _ = io.Copy(w, upstreamResp.Body)
	}
}

func upstreamURL(host string, r *http.Request, translatePaths bool) url.URL {
	u := *r.URL
	u.Scheme = "https"
	u.Host = host

	// If configured to do so, we want to translate a path like:
	//   https://download.maxmind.com/geoip/databases/GeoLite2-Country.tar.gz
	// to:
	//   https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz
	if translatePaths {
		if p, found := strings.CutSuffix(u.Path, ".tar.gz"); found {
			newPath := path.Join(p, "download")
			slog.Debug("Translate path", "from", u.Path, "to", newPath)
			u.Path = newPath
			q := u.Query()
			q.Set("suffix", "tar.gz")
			u.RawQuery = q.Encode()
		}
	}

	return u
}

func setAuth(conf *config.Config, r *http.Request) {
	if conf.AccountID != 0 && conf.LicenseKey != "" {
		q := r.URL.Query()
		if q.Has("license_key") {
			q.Set("license_key", conf.LicenseKey)
			r.URL.RawQuery = q.Encode()
		} else {
			r.SetBasicAuth(strconv.Itoa(conf.AccountID), conf.LicenseKey)
		}
	}
}
