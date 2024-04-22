package redis

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/redis/rueidis"
	"github.com/rs/zerolog/log"
)

type Client struct {
	rueidis.Client
}

func Connect(conf *config.Config) (*Client, error) {
	addr := net.JoinHostPort(conf.RedisHost, strconv.Itoa(int(conf.RedisPort)))
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{addr},
		Password:     conf.RedisPassword,
		SelectDB:     conf.RedisDB,
		DisableCache: true,
	})
	if err != nil {
		return nil, err
	}

	log.Info().Str("addr", addr).Int("db", conf.RedisDB).Msg("connected to redis")
	return &Client{client}, nil
}

func (c *Client) Close() {
	log.Info().Msg("disconnecting from redis")
	c.Client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Do(ctx, c.B().Ping().Build()).Error()
}

func FormatCacheKey(u url.URL, req *http.Request) string {
	q := u.Query()
	q.Del("db_md5")
	u.RawQuery = q.Encode()
	key := req.Method + " " + u.String() + " " + req.Header.Get("Authorization")
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

var ErrNotExist = errors.New("key not found")

func (c *Client) GetCache(ctx context.Context, u url.URL, req *http.Request) (*http.Response, error) {
	r, err := c.Do(ctx, c.B().Get().Key(FormatCacheKey(u, req)).Build()).AsReader()
	if err != nil {
		var redisErr *rueidis.RedisError
		if errors.As(err, &redisErr) {
			if redisErr.IsNil() {
				return nil, ErrNotExist
			}
		}
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(r), req)
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

	if err := c.Do(ctx,
		c.B().Set().Key(FormatCacheKey(u, req)).Value(rueidis.BinaryString(b)).Ex(expiration).Build(),
	).Error(); err != nil {
		return err
	}

	return nil
}
