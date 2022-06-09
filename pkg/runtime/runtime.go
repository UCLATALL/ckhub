package runtime

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/multierr"

	"github.com/uclatall/ckhub/pkg/logging"
)

// Runtime implements an application runtime.
type Runtime struct {
	log logging.Logger

	services []Service

	signals []os.Signal
	timeout time.Duration
}

// NewRuntime creates a new application runtime with the given options.
func NewRuntime(options ...Option) (*Runtime, error) {
	run := &Runtime{
		log:     logging.NopLogger(),
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		timeout: 30 * time.Second,
	}

	errs := make([]error, len(options))
	for i, option := range options {
		errs[i] = option.Apply(run)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return nil, err
	}
	return run, nil
}

// ErrShutdownInterrupt is returned when the runtime shutdown is interrupted.
var ErrShutdownInterrupt = errors.New("shutdown interrupted")

// Run kicks off configured services, blocks on the signal channel, and then
// gracefully shuts them down.
func (r *Runtime) Run(ctx context.Context) error {
	log := r.log

	rctx, cancel := NewContext(ctx)
	defer cancel()

	total := int64(len(r.services))
	tasks := total

	done := make(chan struct{})
	if total == 0 {
		close(done)
	}

	log.Debug("starting runtime")
	errs := make([]error, total)
	for i := int64(0); i < total; i++ {
		go func(n int64) {
			errs[n] = r.services[n].Run(rctx)
			if atomic.AddInt64(&tasks, -1) == 0 {
				close(done)
			}
		}(i)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, r.signals...)
	defer signal.Stop(sigs)

	select {
	case <-done:
		log.Debug("runtime stopped")
		return multierr.Combine(errs...)
	case sig := <-sigs:
		log.Debug("runtime shutdown", logging.Stringer("signal", sig))
		cancel()
	case <-ctx.Done():
		log.Debug("runtime shutdown", logging.Error(ctx.Err()))
		cancel()
	}

	log = log.Hooks(logging.Span())

	sctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	select {
	case <-done:
		log.Debug("runtime stopped")
		return multierr.Combine(errs...)
	case sig := <-sigs:
		log.Debug("runtime interrupted", logging.Stringer("signal", sig))
		return ErrShutdownInterrupt
	case <-sctx.Done():
		log.Debug("runtime interrupted", logging.Error(sctx.Err()))
		return ErrShutdownInterrupt
	}
}
