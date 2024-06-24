package redis

import (
	"bytes"
	"context"
	"io"
	"strconv"
)

type CacheReader struct {
	ctx    context.Context
	cache  *Client
	key    string
	buf    bytes.Buffer
	chunks int
	chunk  int
}

func (c *CacheReader) Read(p []byte) (int, error) {
	if c.buf.Len() != 0 {
		return c.buf.Read(p)
	}

	if c.chunk >= c.chunks {
		return 0, io.EOF
	}

	key := c.key + "_body_" + strconv.Itoa(c.chunk)
	b, err := c.cache.Do(c.ctx, c.cache.B().Get().Key(key).Build()).AsBytes()
	if err != nil {
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
