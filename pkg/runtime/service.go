package runtime

import (
	"context"
	"sync/atomic"

	"go.uber.org/multierr"
)

// Service is a generic interface of the long-running application routine.
type Service interface {
	// Run kicks off the service. It blocks the execution and interrupts when
	// the context is canceled.
	Run(ctx context.Context) error
}

// ServiceFunc is an adapter to allow the use of ordinary functions as services.
type ServiceFunc func(ctx context.Context) error

// Run kicks off the service. It blocks the execution and interrupts when
// the context is canceled.
func (svc ServiceFunc) Run(ctx context.Context) error {
	return svc(ctx)
}

// Group represents a group of services.
type Group []Service

// Apply appends the group to the application runtime.
func (g Group) Apply(r *Runtime) error {
	r.services = append(r.services, g)
	return nil
}

// Run kicks off underlying services. It blocks the execution and interrupts
// when the context is canceled, or any of the services terminate.
func (g Group) Run(ctx context.Context) error {
	rctx, cancel := NewContext(ctx)
	defer cancel()

	total := int64(len(g))
	tasks := total

	done := make(chan struct{})
	if total == 0 {
		close(done)
	}

	errs := make([]error, total)
	for i := int64(0); i < total; i++ {
		go func(n int64) {
			defer cancel()

			errs[n] = g[n].Run(rctx)
			if atomic.AddInt64(&tasks, -1) == 0 {
				close(done)
			}
		}(i)
	}

	select {
	case <-done:
	case <-ctx.Done():
		cancel()
		<-done
	}

	return multierr.Combine(errs...)
}
