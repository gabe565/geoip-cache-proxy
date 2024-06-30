package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/rueidis"
)

type CacheWriter struct {
	ctx           context.Context
	cache         *Client
	key           string
	expiration    time.Time
	chunk         int
	contentLength int64
	written       int64
}

func (c *CacheWriter) Write(p []byte) (int, error) {
	key := c.key + "_body_" + strconv.Itoa(c.chunk)
	err := c.cache.Do(c.ctx,
		c.cache.B().Set().Key(key).Value(rueidis.BinaryString(p)).Exat(c.expiration).Build(),
	).Error()
	if err == nil {
		c.chunk++
		c.written += int64(len(p))
	}
	return len(p), err
}

func (c *CacheWriter) Close() error {
	defer locks.Unlock(c.key)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if c.contentLength != -1 && c.written != c.contentLength {
		return nil
	}
	key := c.key + "_chunks"
	return c.cache.Do(ctx,
		c.cache.B().Set().Key(key).Value(strconv.Itoa(c.chunk)).Exat(c.expiration).Build(),
	).Error()
}
