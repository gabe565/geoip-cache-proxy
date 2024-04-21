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

		upstreamReq, err := http.NewRequestWithContext(ctx, r.Method, u.String(), r.Body)
		if err != nil {
			log.Err(err).Msg("failed to create request")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			upstreamReq.Header.Set("X-Forwarded-For", host)
		} else {
			log.Warn().Err(err).Msg("failed to split remote address")
		}
		for k := range r.Header {
			upstreamReq.Header.Set(k, r.Header.Get(k))
		}

		var cacheStatus CacheStatus
		var upstreamResp *http.Response
		if upstreamResp, err = cache.GetCache(ctx, u, upstreamReq); err == nil {
			log.Trace().Msg("using cached response")
			cacheStatus = CacheHit
			defer func(Body io.ReadCloser) {
				_, _ = io.Copy(io.Discard, upstreamResp.Body)
				_ = Body.Close()
			}(upstreamResp.Body)
		} else {
			log.Trace().Err(err).Msg("failed to get cached response")
			upstreamResp = nil
		}

		if upstreamResp == nil {
			log.Trace().Msg("forwarding request to upstream")
			cacheStatus = CacheMiss
			upstreamResp, err = http.DefaultClient.Do(upstreamReq)
			if err != nil {
				log.Err(err).Msg("failed to forward to upstream")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func(Body io.ReadCloser) {
				_, _ = io.Copy(io.Discard, upstreamResp.Body)
				_ = Body.Close()
			}(upstreamResp.Body)

			if upstreamResp.StatusCode < 400 {
				upstreamResp, err = cache.SetCache(ctx, u, upstreamReq, upstreamResp)
				defer func(Body io.ReadCloser) {
					_, _ = io.Copy(io.Discard, upstreamResp.Body)
					_ = Body.Close()
				}(upstreamResp.Body)
				if err != nil {
					log.Err(err).Msg("failed to set cache response")
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
