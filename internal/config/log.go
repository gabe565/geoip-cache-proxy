package config

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const LevelTrace = slog.Level(-5)

//go:generate go run github.com/dmarkham/enumer -type LogFormat -trimprefix Format -transform lower -text

type LogFormat uint8

const (
	FormatAuto LogFormat = iota
	FormatColor
	FormatPlain
	FormatJSON
)

func (c *Config) LogLevel() (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(c.logLevel))
	if err != nil {
		switch c.logLevel {
		case "trace":
			level = LevelTrace
			err = nil
		default:
			level = slog.LevelInfo
		}
	}
	return level, err
}

func (c *Config) LogFormat() (LogFormat, error) {
	var format LogFormat
	err := format.UnmarshalText([]byte(c.logFormat))
	if err != nil {
		format = FormatAuto
	}
	return format, err
}

func (c *Config) InitLog(w io.Writer) {
	level, err := c.LogLevel()
	if err != nil {
		defer func() {
			slog.Warn("Invalid log level. Defaulting to info.", "value", c.logLevel)
		}()
		c.logLevel = strings.ToLower(level.String())
	}

	format, err := c.LogFormat()
	if err != nil {
		defer func() {
			slog.Warn("Invalid log format. Defaulting to auto.", "value", c.logFormat)
		}()
		c.logFormat = format.String()
	}

	InitLog(w, level, format)
}

func InitLog(w io.Writer, level slog.Level, format LogFormat) {
	switch format {
	case FormatJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		})))
	default:
		var color bool
		switch format {
		case FormatAuto:
			if f, ok := w.(*os.File); ok {
				color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
			}
		case FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
