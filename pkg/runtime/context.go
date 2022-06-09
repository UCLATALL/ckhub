package runtime

import "context"

type runtimeCtx struct {
	context.Context
	valueCtx context.Context
}

// NewContext returns a new context that blocks propagation of the context
// cancelation, but value calls are forwarded to the provided parent.
func NewContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(&runtimeCtx{
		Context:  context.Background(),
		valueCtx: parent,
	})
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *runtimeCtx) Value(key any) any {
	return c.valueCtx.Value(key)
}
