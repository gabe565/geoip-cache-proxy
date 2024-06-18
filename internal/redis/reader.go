package redis

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/redis/rueidis"
)

type CacheReader struct {
	ctx   context.Context
	cache *Client
	key   string
	buf   bytes.Buffer
	chunk int
}

func (c *CacheReader) Read(p []byte) (int, error) {
	if c.buf.Len() != 0 {
		return c.buf.Read(p)
	}

	key := c.key + "_body_" + strconv.Itoa(c.chunk)
	b, err := c.cache.Do(c.ctx, c.cache.B().Get().Key(key).Build()).AsBytes()
	if err != nil {
		var redisErr *rueidis.RedisError
		if errors.As(err, &redisErr) {
			if redisErr.IsNil() {
				return 0, io.EOF
			}
		}
		return 0, err
	}
	c.buf.Write(b)
	c.chunk++

	return c.buf.Read(p)
}

func (c *CacheReader) Close() error {
	locks.Unlock(c.key)
	return nil
}
