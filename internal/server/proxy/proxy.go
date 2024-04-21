package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/cache"
	"github.com/gabe565/geoip-cache-proxy/internal/server/consts"
	"github.com/gabe565/geoip-cache-proxy/internal/server/middleware"
)

func Proxy(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		u := buildURL(host, r)
		log := middleware.LogFromContext(r.Context()).With().Str("upstreamUrl", u.String()).Logger()

		req, err := http.NewRequestWithContext(ctx, r.Method, u.String(), r.Body)
		if err != nil {
			log.Err(err).Msg("failed to create request")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			req.Header.Set("X-Forwarded-For", host)
		} else {
			log.Warn().Err(err).Msg("failed to split remote address")
		}
		for k := range r.Header {
			req.Header.Set(k, r.Header.Get(k))
		}

		var cacheStatus CacheStatus
		var resp *http.Response
		if resp, err = cache.GetCache(ctx, u, req); err == nil {
			log.Trace().Msg("using cached response")
			cacheStatus = CacheHit
		} else {
			log.Trace().Msg("failed to get cached response")
			resp = nil
		}

		if resp == nil {
			log.Trace().Msg("forwarding request to upstream")
			cacheStatus = CacheMiss
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Err(err).Msg("failed to forward to upstream")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}()

			if resp.StatusCode < 300 {
				resp, err = cache.SetCache(ctx, u, req, resp)
				if err != nil {
					log.Err(err).Msg("failed to set cache response")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				cacheStatus = CacheBypass
			}
		}

		for k := range resp.Header {
			w.Header().Set(k, resp.Header.Get(k))
		}
		w.Header().Set(consts.UpstreamURLHeader, u.String())
		w.Header().Set(consts.CacheStatusHeader, cacheStatus.String())
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
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
