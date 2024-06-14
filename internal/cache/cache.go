package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
)

//nolint:gochecknoglobals
var mu sync.RWMutex

func cacheKey(req *http.Request) string {
	u := *req.URL
	q := u.Query()
	q.Del("db_md5")
	u.RawQuery = q.Encode()
	sum := sha256.Sum256([]byte(u.String()))
	encoded := hex.EncodeToString(sum[:])
	return encoded
}

func dataPaths(path string, req *http.Request) (string, string) {
	key := cacheKey(req)
	return filepath.Join(path, key+"_headers.bin"), filepath.Join(path, key+"_body.bin")
}

func EnsureCache(conf *config.Config) error {
	return os.MkdirAll(conf.CacheDir, 0o700)
}
