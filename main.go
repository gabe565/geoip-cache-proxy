package main

import (
	"log/slog"
	"net/http"
	"os"

	"gabe565.com/geoip-cache-proxy/cmd"
	"gabe565.com/geoip-cache-proxy/internal/config"
	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/httpx"
	"gabe565.com/utils/slogx"
)

var version = "beta"

func main() {
	config.InitLog(os.Stderr, slogx.LevelInfo, slogx.FormatAuto)
	root := cmd.New(cobrax.WithVersion(version))
	http.DefaultTransport = httpx.NewUserAgentTransport(nil, cobrax.BuildUserAgent(root))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
