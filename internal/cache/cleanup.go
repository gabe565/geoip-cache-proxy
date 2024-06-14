package cache

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/gabe565/geoip-cache-proxy/internal/config"
	"github.com/rs/zerolog/log"
)

func BeginCleanup(ctx context.Context, conf *config.Config) {
	go func() {
		ticker := time.NewTicker(conf.CleanupEvery)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := filepath.Walk(conf.CacheDir, func(path string, info fs.FileInfo, err error) error {
					if err != nil || info.IsDir() {
						return err
					}

					if time.Since(info.ModTime()) > conf.CacheDuration {
						log.Trace().Str("path", path).Msg("cleaning up")

						mu.Lock()
						defer mu.Unlock()

						if err := os.RemoveAll(path); err != nil {
							log.Err(err).Str("path", path).Msg("failed to remove file")
						}
					}

					return nil
				}); err != nil {
					log.Err(err).Msg("failed to walk data dir")
				}
			}
		}
	}()
}
