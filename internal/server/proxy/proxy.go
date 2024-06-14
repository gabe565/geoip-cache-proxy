package proxy

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gabe565/geoip-cache-proxy/internal/cache"
	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/server/consts"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func Proxy(conf *config.Config, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := upstreamURL(host, r)
		log := log.Ctx(r.Context()).With().Str("upstreamUrl", u.String()).Logger()

		upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, u.String(), r.Body)
		if err != nil {
			log.Err(err).Msg("failed to create request")
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

		if upstreamResp, err = cache.Get(conf.CacheDir, upstreamReq); err == nil {
			log.Trace().Msg("using cached response")
			cacheStatus = CacheHit
		} else if errors.Is(err, os.ErrNotExist) {
			log.Trace().Msg("forwarding request to upstream")
			upstreamResp, err = http.DefaultClient.Do(upstreamReq)
			if err != nil {
				log.Err(err).Msg("failed to forward to upstream")
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}

			if upstreamResp.StatusCode < 300 {
				cacheWriter, err := cache.NewWriter(conf.CacheDir, upstreamReq, upstreamResp)
				if err != nil {
					log.Err(err).Msg("failed to set cache response")
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				defer func() {
					if err := cacheWriter.Close(); err != nil {
						log.Err(err).Msg("failed to close cache")
					}
				}()

				wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				wrapped.Tee(cacheWriter)
				w = wrapped
			} else {
				cacheStatus = CacheBypass
			}
		} else {
			log.Trace().Err(err).Msg("failed to get cached response")
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

func upstreamURL(host string, r *http.Request) url.URL {
	u := *r.URL
	u.Scheme = "https"
	u.Host = host
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
