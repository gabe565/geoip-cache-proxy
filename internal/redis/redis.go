package redis

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"gabe565.com/geoip-cache-proxy/internal/config"
	"gabe565.com/geoip-cache-proxy/internal/util"
	"github.com/redis/rueidis"
)

//nolint:gochecknoglobals
var locks = util.NewKeyedLock()

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

	slog.Info("Connected to redis", "address", addr, "db", conf.RedisDB)
	return &Client{client}, nil
}

func (c *Client) Close() {
	slog.Info("Disconnecting from redis")
	c.Client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Do(ctx, c.B().Ping().Build()).Error()
}

func FormatCacheKey(req *http.Request) string {
	u := *req.URL
	q := u.Query()
	q.Del("db_md5")
	u.RawQuery = q.Encode()
	sum := sha256.Sum256([]byte(u.String()))
	return hex.EncodeToString(sum[:])
}

var ErrNotExist = errors.New("key not found")

func (c *Client) Get(ctx context.Context, req *http.Request, httpTimeout time.Duration) (*http.Response, error) {
	key := FormatCacheKey(req)

	locks.Lock(key)
	var resp *http.Response
	defer func() {
		if resp == nil {
			locks.Unlock(key)
		}
	}()

	ttl, err := c.Do(ctx, c.B().Ttl().Key(key+"_headers").Build()).AsInt64()
	if err != nil {
		return nil, err
	}
	if ttl == -2 || ttl < int64(httpTimeout/time.Second) {
		return nil, ErrNotExist
	}

	r, err := c.Do(ctx, c.B().Get().Key(key+"_headers").Build()).AsReader()
	if err != nil {
		var redisErr *rueidis.RedisError
		if errors.As(err, &redisErr) && redisErr.IsNil() {
			return nil, ErrNotExist
		}
		return nil, err
	}

	chunks, err := c.Do(ctx, c.B().Get().Key(key+"_chunks").Build()).AsInt64()
	if err != nil {
		var redisErr *rueidis.RedisError
		if errors.As(err, &redisErr) && redisErr.IsNil() {
			return nil, ErrNotExist
		}
		return nil, err
	}

	if resp, err = http.ReadResponse(bufio.NewReader(r), req); err != nil {
		return nil, err
	}

	resp.Body = &CacheReader{
		ctx:    ctx,
		cache:  c,
		key:    key,
		chunks: int(chunks),
	}
	return resp, nil
}

func (c *Client) NewWriter(
	ctx context.Context,
	req *http.Request,
	resp *http.Response,
	expiration time.Duration,
) (io.WriteCloser, error) {
	key := FormatCacheKey(req)

	locks.Lock(key)
	var w io.WriteCloser
	defer func() {
		if w == nil {
			locks.Unlock(key)
		}
	}()

	b, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, err
	}

	exp := time.Now().Add(expiration)

	if err := c.Do(ctx,
		c.B().Set().Key(key+"_headers").Value(rueidis.BinaryString(b)).Exat(exp).Build(),
	).Error(); err != nil {
		return nil, err
	}

	w = &CacheWriter{
		ctx:           ctx,
		cache:         c,
		key:           key,
		expiration:    exp,
		contentLength: resp.ContentLength,
	}
	return w, nil
}
