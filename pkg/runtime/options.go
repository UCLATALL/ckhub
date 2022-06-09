package runtime

import (
	"os"
	"time"

	"github.com/uclatall/ckhub/pkg/logging"
)

// Option is a generic interface of the runtime configuration option.
type Option interface {
	// Apply applies the option to the given runtime.
	Apply(r *Runtime) error
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(r *Runtime) error

// Apply applies the option to the given runtime.
func (o OptionFunc) Apply(r *Runtime) error {
	return o(r)
}

// Logger creates a new option that sets the logger for the runtime.
func Logger(log logging.Logger) OptionFunc {
	return func(r *Runtime) error {
		r.log = log
		return nil
	}
}

// Services creates a new option that appends given services to the runtime.
func Services(services ...Service) OptionFunc {
	return func(r *Runtime) error {
		r.services = append(r.services, services...)
		return nil
	}
}

// Signals creates a new option that sets signals for the runtime shutdown.
func Signals(sigs ...os.Signal) OptionFunc {
	return func(r *Runtime) error {
		r.signals = sigs
		return nil
	}
}

// Timeout creates a new option that sets timeout for the runtime shutdown.
func Timeout(timeout time.Duration) OptionFunc {
	return func(r *Runtime) error {
		r.timeout = timeout
		return nil
	}
}
