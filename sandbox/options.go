package sandbox

import "github.com/uclatall/ckhub/pkg/logging"

// Option is a generic interface of the sandbox configuration option.
type Option interface {
	// Apply applies the option to the given sandbox.
	Apply(srv *Manager) error
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(srv *Manager) error

// Apply applies the option to the given sandbox.
func (o OptionFunc) Apply(srv *Manager) error {
	return o(srv)
}

// Logger creates a new option that sets the logger for the sandbox.
func Logger(log logging.Logger) OptionFunc {
	return func(srv *Manager) error {
		srv.log = log
		return nil
	}
}
