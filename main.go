package main

import (
	"os"

	"github.com/gabe565/geoip-cache-proxy/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		log.Err(err).Msg("Exiting due to an error")
		os.Exit(1)
	}
}
