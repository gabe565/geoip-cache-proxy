package middleware

import (
	"net/http"
	"time"
)

func Timeout(timeout time.Duration, h http.Handler) http.Handler {
	return http.TimeoutHandler(h, timeout, "timout exceeded")
}
