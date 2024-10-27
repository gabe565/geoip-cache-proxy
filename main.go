package main

import (
	"log/slog"
	"os"

	"gabe565.com/geoip-cache-proxy/cmd"
	"gabe565.com/geoip-cache-proxy/internal/config"
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
