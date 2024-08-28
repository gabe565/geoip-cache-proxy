package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/geoip-cache-proxy/cmd"
	"github.com/gabe565/geoip-cache-proxy/internal/config"
)

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)
	if err := cmd.New().Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
