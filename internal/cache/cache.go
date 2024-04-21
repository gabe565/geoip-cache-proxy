package cache

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/redis/go-redis/v9"
)

//nolint:gochecknoglobals
var Client *redis.Client

func Connect(conf *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	return nil
}

func FormatCacheKey(u url.URL, req *http.Request) string {
	return req.Method + "_" + u.String()
}

func GetCache(ctx context.Context, u url.URL, req *http.Request) (*http.Response, error) {
	b, err := Client.Get(ctx, FormatCacheKey(u, req)).Bytes()
	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SetCache(ctx context.Context, u url.URL, req *http.Request, resp *http.Response, expiration time.Duration) (*http.Response, error) {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if err := Client.Set(ctx, FormatCacheKey(u, req), b, expiration).Err(); err != nil {
		return nil, err
	}

	resp, err = http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
