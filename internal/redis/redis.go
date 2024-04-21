package redis

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Client struct {
	*redis.Client
}

func Connect(ctx context.Context, conf *config.Config) (*Client, error) {
	addr := net.JoinHostPort(conf.RedisHost, strconv.Itoa(int(conf.RedisPort)))
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Str("addr", addr).Int("db", conf.RedisDB).Msg("connected to redis")
	return &Client{client}, nil
}

func (c *Client) Close() error {
	log.Info().Msg("disconnecting from redis")
	return c.Client.Close()
}

func FormatCacheKey(u url.URL, req *http.Request) string {
	key := req.Method + " " + u.String() + " " + req.Header.Get("Authorization")
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

var ErrNotExist = errors.New("key not found")

func (c *Client) GetCache(ctx context.Context, u url.URL, req *http.Request) (*http.Response, error) {
	b, err := c.Get(ctx, FormatCacheKey(u, req)).Bytes()
	if err != nil {
		if redis.HasErrorPrefix(err, "redis: nil") {
			return nil, ErrNotExist
		}
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(b)), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) SetCache(ctx context.Context, u url.URL, req *http.Request, resp *http.Response, expiration time.Duration) error {
	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}

	if err := c.Set(ctx, FormatCacheKey(u, req), b, expiration).Err(); err != nil {
		return err
	}

	return nil
}
