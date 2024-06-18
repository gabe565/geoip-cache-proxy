package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/rueidis"
)

type CacheWriter struct {
	ctx        context.Context
	cache      *Client
	key        string
	expiration time.Time
	chunk      int
}

func (c *CacheWriter) Write(p []byte) (int, error) {
	key := c.key + "_body_" + strconv.Itoa(c.chunk)
	err := c.cache.Do(c.ctx,
		c.cache.B().Set().Key(key).Value(rueidis.BinaryString(p)).Exat(c.expiration).Build(),
	).Error()
	if err == nil {
		c.chunk++
	}
	return len(p), err
}

func (c *CacheWriter) Close() error {
	locks.Unlock(c.key)
	return nil
}
