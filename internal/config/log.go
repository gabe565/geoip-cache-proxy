package config

import (
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func logLevel(level string) zerolog.Level {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil || parsedLevel == zerolog.NoLevel {
		if level == "warning" {
			parsedLevel = zerolog.WarnLevel
		} else {
			log.Warn().Str("value", level).Msg("Invalid log level. Defaulting to info.")
			parsedLevel = zerolog.InfoLevel
		}
	}
	return parsedLevel
}

const (
	FormatJSON  = "json"
	FormatAuto  = "auto"
	FormatColor = "color"
	FormatPlain = "plain"
)

func logFormat(out io.Writer, format string) io.Writer {
	switch format {
	case FormatJSON:
		return out
	default:
		var useColor bool
		switch format {
		case FormatAuto:
			if w, ok := out.(*os.File); ok {
				useColor = isatty.IsTerminal(w.Fd())
			}
		case FormatColor:
			useColor = true
		case FormatPlain:
		default:
			log.Warn().Str("value", format).Msg("Invalid log formatter. Defaulting to auto.")
		}

		return zerolog.ConsoleWriter{
			Out:        out,
			NoColor:    !useColor,
			TimeFormat: time.DateTime,
		}
	}
}

func initLog(cmd *cobra.Command) {
	level, err := cmd.Flags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(logLevel(level))

	format, err := cmd.Flags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(logFormat(cmd.ErrOrStderr(), format))
}
