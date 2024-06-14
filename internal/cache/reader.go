package cache

import (
	"bufio"
	"net/http"
	"os"
)

func Get(path string, req *http.Request) (*http.Response, error) {
	headerPath, bodyPath := dataPaths(path, req)

	mu.RLock()

	header, err := os.Open(headerPath)
	if err != nil {
		mu.RUnlock()
		return nil, err
	}
	defer func() {
		_ = header.Close()
	}()

	resp, err := http.ReadResponse(bufio.NewReader(header), req)
	if err != nil {
		mu.RUnlock()
		return nil, err
	}

	body, err := os.Open(bodyPath)
	if err != nil {
		mu.RUnlock()
		return nil, err
	}

	resp.Body = ROFile{body}
	return resp, nil
}

type ROFile struct {
	*os.File
}

func (c ROFile) Close() error {
	err := c.File.Close()
	mu.RUnlock()
	return err
}
