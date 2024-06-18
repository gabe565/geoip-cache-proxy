package redis

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/gabe565/geoip-cache-proxy/internal/util"
	"github.com/redis/rueidis"
	"github.com/rs/zerolog/log"
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

func FormatCacheKey(req *http.Request) string {
	u := req.URL
	q := u.Query()
	q.Del("db_md5")
	u.RawQuery = q.Encode()
	sum := sha256.Sum256([]byte(u.String()))
	return hex.EncodeToString(sum[:])
}

var ErrNotExist = errors.New("key not found")

func (c *Client) Get(ctx context.Context, req *http.Request) (*http.Response, error) {
	key := FormatCacheKey(req)

	locks.Lock(key)

	r, err := c.Do(ctx, c.B().Get().Key(key+"_headers").Build()).AsReader()
	if err != nil {
		locks.Unlock(key)
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
		locks.Unlock(key)
		return nil, err
	}

	resp.Body = &CacheReader{
		ctx:   ctx,
		cache: c,
		key:   key,
	}
	return resp, nil
}

func (c *Client) NewWriter(ctx context.Context, req *http.Request, resp *http.Response, expiration time.Duration) (io.WriteCloser, error) {
	key := FormatCacheKey(req)
	locks.Lock(key)

	b, err := httputil.DumpResponse(resp, false)
	if err != nil {
		locks.Unlock(key)
		return nil, err
	}

	exp := time.Now().Add(expiration)

	if err := c.Do(ctx,
		c.B().Set().Key(key+"_headers").Value(rueidis.BinaryString(b)).Exat(exp).Build(),
	).Error(); err != nil {
		locks.Unlock(key)
		return nil, err
	}

	w := &CacheWriter{
		ctx:        ctx,
		cache:      c,
		key:        key,
		expiration: exp,
	}
	return w, nil
}
