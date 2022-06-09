package logging

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

// Hook is a generic interface of the logging hook.
type Hook interface {
	// Hook intercepts the given log entry.
	Hook(entry *zapcore.CheckedEntry)
}

// HookFunc is an adapter to allow the use of ordinary functions as hooks.
type HookFunc func(entry *zapcore.CheckedEntry)

// Hook intercepts the given log entry.
func (h HookFunc) Hook(entry *zapcore.CheckedEntry) {
	h(entry)
}

// Span creates a new hook that appends an operation duration to the log entry.
func Span() HookFunc {
	start := time.Now()
	return func(entry *zapcore.CheckedEntry) {
		entry.Write(Duration(FieldSpan, time.Since(start)))
	}
}

// Trace creates a new hook that appends trace to the log entry.
func Trace(id fmt.Stringer) HookFunc {
	start := time.Now()
	return func(entry *zapcore.CheckedEntry) {
		entry.Write(
			Stringer(FieldTrace, id),
			Duration(FieldSpan, time.Since(start)),
		)
	}
}
