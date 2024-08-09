package proxy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/redis"
	"github.com/gabe565/geoip-cache-proxy/internal/server/consts"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func Proxy(conf *config.Config, cache *redis.Client, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := upstreamURL(host, r, conf.TranslateIngressNginxPaths)
		log := log.Ctx(r.Context()).With().Str("upstreamUrl", u.String()).Logger()

		upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, u.String(), r.Body)
		if err != nil {
			log.Err(err).Msg("Failed to create request")
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
			log.Trace().Msg("Using cached response")
			cacheStatus = CacheHit
		} else if errors.Is(err, redis.ErrNotExist) {
			log.Trace().Msg("Forwarding request to upstream")
			upstreamResp, err = http.DefaultClient.Do(upstreamReq)
			if err != nil {
				log.Err(err).Msg("Failed to forward to upstream")
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}

			if upstreamResp.StatusCode < 300 {
				if cacheWriter, err := cache.NewWriter(r.Context(), upstreamReq, upstreamResp, conf.CacheDuration); err == nil {
					defer func() {
						if err := cacheWriter.Close(); err != nil {
							log.Err(err).Msg("Failed to close cache")
						}
					}()

					wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
					wrapped.Tee(cacheWriter)
					w = wrapped
				} else {
					log.Err(err).Msg("Failed to cache response")
				}
			} else {
				cacheStatus = CacheBypass
			}
		} else {
			log.Trace().Err(err).Msg("Failed to get cached response")
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
		log.Debug().Msg("translating paths for ingress-nginx")
		pat := regexp.MustCompile(`(.+)\.tar\.gz`)

		matches := pat.FindStringSubmatch(u.Path)
		if len(matches) > 1 {
			newPath := fmt.Sprintf("%s/download?suffix=tar.gz", matches[1])
			log.Debug().Msg(fmt.Sprintf("translating %s into %s", u.Path, newPath))
			u.Path = newPath
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
