package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/server/consts"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni/v3"
)

type ctxKey uint8

const logCtxKey ctxKey = iota

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger := log.With().
			Str("method", r.Method).
			Str("requestUrl", r.URL.String()).
			Str("remoteAddr", r.RemoteAddr).
			Str("userAgent", r.UserAgent()).
			Str("latency", time.Since(start).String()).
			Str("protocol", r.Proto).
			Logger()

		resp := negroni.NewResponseWriter(w)
		ctx := LogContext(r.Context(), logger)
		next.ServeHTTP(resp, r.WithContext(ctx))

		level := zerolog.DebugLevel
		if resp.Status() >= 400 {
			level = zerolog.InfoLevel
		}

		logger.WithLevel(level).
			Str("status", strconv.Itoa(resp.Status())).
			Str("responseSize", strconv.Itoa(resp.Size())).
			Str("upstreamUrl", resp.Header().Get(consts.UpstreamURLHeader)).
			Str("cacheStatus", resp.Header().Get(consts.CacheStatusHeader)).
			Msg("served request")
	})
}

func LogContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey, logger)
}

func LogFromContext(ctx context.Context) zerolog.Logger {
	return ctx.Value(logCtxKey).(zerolog.Logger)
}
