package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (c *Config) Load(cmd *cobra.Command) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		optName := strings.ToUpper(f.Name)
		optName = strings.ReplaceAll(optName, "-", "_")
		varName := EnvPrefix + optName
		if val, ok := os.LookupEnv(varName); !f.Changed && ok {
			if err := f.Value.Set(val); err != nil {
				errs = append(errs, err)
			}
		}
	})
	initLog(cmd)
	return errors.Join(errs...)
}
