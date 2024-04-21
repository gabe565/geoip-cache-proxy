package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "GEOIP_"

func (c *Config) Load(cmd *cobra.Command) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if val, ok := os.LookupEnv(EnvName(f.Name)); !f.Changed && ok {
			if err := f.Value.Set(val); err != nil {
				errs = append(errs, err)
			}
		}
	})
	initLog(cmd)
	return errors.Join(errs...)
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
