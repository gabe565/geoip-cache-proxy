package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/geoip-cache-proxy/cmd"
	"github.com/gabe565/geoip-cache-proxy/internal/config"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)
	root := cmd.New(cmd.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
