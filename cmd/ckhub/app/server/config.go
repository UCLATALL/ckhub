package server

import (
	"fmt"
	"os"

	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"github.com/uclatall/ckhub/sandbox"
	"github.com/uclatall/ckhub/sandbox/server"
)

// Config represents a server configuration.
type Config struct {
	Server  server.Config  `json:"server" yaml:"server"`
	Sandbox sandbox.Config `json:"sandbox" yaml:"sandbox"`
}

// NewConfig creates a new server configuration with the given options.
func NewConfig(options ...Option) (*Config, error) {
	config := &Config{
		Server: server.Config{
			Address: ":8080",
		},
		Sandbox: sandbox.Config{},
	}

	errs := make([]error, len(options))
	for i, option := range options {
		errs[i] = option.Apply(config)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Option is a generic interface of the server configuration option.
type Option interface {
	// Apply applies the option to the given server configuration.
	Apply(cfg *Config) error
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(cfg *Config) error

// Apply applies the option to the given server configuration.
func (o OptionFunc) Apply(cfg *Config) error {
	return o(cfg)
}

// Path creates a new option that reads configuration from the given path.
func Path(path string) OptionFunc {
	return func(cfg *Config) error {
		//nolint:gosec // Reads configuration from the given path.
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer func() { _ = f.Close() }()

		decoder := yaml.NewDecoder(f)
		decoder.KnownFields(true)
		err = decoder.Decode(cfg)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		return nil
	}
}
