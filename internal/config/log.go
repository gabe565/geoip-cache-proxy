package config

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
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
			log.Warn().Str("value", level).Msg("invalid log level. defaulting to info.")
			parsedLevel = zerolog.InfoLevel
		}
	}
	return parsedLevel
}

func logFormat(out io.Writer, format string) io.Writer {
	switch format {
	case "json", "j":
		return out
	default:
		sprintf := fmt.Sprintf
		var useColor bool
		switch format {
		case "auto", "a":
			if w, ok := out.(*os.File); ok {
				useColor = isatty.IsTerminal(w.Fd())
				if useColor {
					sprintf = color.New(color.Bold).Sprintf
				}
			}
		case "color", "c":
			useColor = true
			sprintf = color.New(color.Bold).Sprintf
		case "plain", "p":
		default:
			log.Warn().Str("value", format).Msg("invalid log formatter. defaulting to auto.")
		}

		return zerolog.ConsoleWriter{
			Out:     out,
			NoColor: !useColor,
			FormatMessage: func(i interface{}) string {
				return sprintf("%-45s", i)
			},
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
