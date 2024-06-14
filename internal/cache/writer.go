package cache

import (
	"net/http"
	"net/http/httputil"
	"os"
)

func NewWriter(path string, req *http.Request, resp *http.Response) (*Writer, error) {
	headerPath, bodyPath := dataPaths(path, req)

	b, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, err
	}

	mu.Lock()

	if err := os.WriteFile(headerPath, b, 0o400); err != nil {
		mu.Unlock()
		return nil, err
	}

	f, err := os.OpenFile(bodyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o400)
	if err != nil {
		mu.Unlock()
		return nil, err
	}

	return &Writer{f: f}, nil
}

type Writer struct {
	f *os.File
}

func (c Writer) Write(p []byte) (int, error) {
	return c.f.Write(p)
}

func (c Writer) Close() error {
	err := c.f.Close()
	mu.Unlock()
	return err
}
