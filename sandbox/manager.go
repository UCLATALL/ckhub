package sandbox

import (
	"context"
	"errors"
	"time"

	"go.uber.org/multierr"

	"github.com/uclatall/ckhub/pkg/logging"
)

// Manager implements a sandbox management service.
type Manager struct {
	log logging.Logger

	kernels map[string]*Kernel
	spawn   time.Duration
}

// NewManager creates a new sandbox manager with the given options.
func NewManager(options ...Option) (*Manager, error) {
	manager := &Manager{
		log:     logging.NopLogger(),
		kernels: make(map[string]*Kernel),
	}

	errs := make([]error, len(options))
	for i, option := range options {
		errs[i] = option.Apply(manager)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// Run kicks off the manager. It blocks the execution and interrupts when
// the given context is canceled.
func (m *Manager) Run(ctx context.Context) error {
	log := m.log

	ticker := time.NewTicker(500 * time.Millisecond)

loop:
	for {
		select {
		case <-ticker.C:
			for name, kernel := range m.kernels {
				err := kernel.SpawnInstance(ctx)
				if err != nil {
					log.Error(
						"failed to spawn new kernel",
						logging.String("name", name),
						logging.Error(err),
					)
				}
			}
		case <-ctx.Done():
			log.Debug("manager shutdown", logging.Error(ctx.Err()))
			break loop
		}
	}

	log = log.Hooks(logging.Span())

	for _, kernel := range m.kernels {
		err := kernel.Destroy()
		if err != nil {
			log.Error("failed to destroy kernel", logging.Error(err))
		}
	}

	log.Debug("manager stopped")

	return nil
}

// ErrKernelNotFound is returned when a kernel is not found.
var ErrKernelNotFound = errors.New("kernel not found")

// ExecuteSnippet executes the given snippet.
func (m *Manager) ExecuteSnippet(ctx context.Context, snippet *Snippet) (*Result, error) {
	kernel, ok := m.kernels[snippet.Kernel]
	if !ok {
		return nil, ErrKernelNotFound
	}

	return kernel.ExecuteSnippet(ctx, snippet)
}

// Kernels returns the number of available kernels.
func (m *Manager) Kernels() int {
	total := 0
	for _, kernel := range m.kernels {
		total += kernel.Instances()
	}
	return total
}
