package server

import (
	"github.com/spf13/pflag"

	"github.com/uclatall/ckhub/pkg/logging"
)

// Flags describes command-line flags for server management command.
type Flags struct {
	config string
	debug  bool
}

// NewFlags creates command-line flags for the server management command.
func NewFlags() Flags {
	return Flags{
		config: "",
		debug:  false,
	}
}

// Register registers flags to the given flagset.
func (f *Flags) Register(flags *pflag.FlagSet) {
	flags.StringVarP(&f.config, "config", "c", f.config, "path to the server config")
	flags.BoolVar(&f.debug, "debug", f.debug, "enable verbose logging")
}

// Logger creates an option that configures a logger with provided flags.
func (f *Flags) Logger() logging.OptionFunc {
	return func(log *logging.Logger) {
		if f.debug {
			logging.LevelDebug.Apply(log)
		}
	}
}

// Config creates an option that configures a server with provided flags.
func (f *Flags) Config() OptionFunc {
	return func(cfg *Config) error {
		if f.config != "" {
			return Path(f.config).Apply(cfg)
		}
		return nil
	}
}
