package proxy

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/redis"
	"github.com/gabe565/geoip-cache-proxy/internal/server/consts"
	"github.com/rs/zerolog/log"
)

func Proxy(conf *config.Config, cache *redis.Client, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := buildURL(host, r)
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
		if upstreamResp, err = cache.GetCache(r.Context(), u, upstreamReq); err == nil {
			log.Trace().Msg("using cached response")
			cacheStatus = CacheHit
			defer func(Body io.ReadCloser) {
				_, _ = io.Copy(io.Discard, upstreamResp.Body)
				_ = Body.Close()
			}(upstreamResp.Body)
		} else if !errors.Is(err, redis.ErrNotExist) {
			log.Trace().Err(err).Msg("failed to get cached response")
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		if cacheStatus == CacheMiss {
			log.Trace().Msg("forwarding request to upstream")
			upstreamResp, err = http.DefaultClient.Do(upstreamReq)
			if err != nil {
				log.Err(err).Msg("failed to forward to upstream")
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
			defer func(Body io.ReadCloser) {
				_, _ = io.Copy(io.Discard, upstreamResp.Body)
				_ = Body.Close()
			}(upstreamResp.Body)

			if upstreamResp.StatusCode < 400 {
				if err = cache.SetCache(r.Context(), u, upstreamReq, upstreamResp, conf.CacheDuration); err != nil {
					log.Err(err).Msg("failed to set cache response")
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
			} else {
				cacheStatus = CacheBypass
			}
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

func buildURL(host string, r *http.Request) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
		Fragment: r.URL.Fragment,
	}
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
