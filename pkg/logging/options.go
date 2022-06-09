package logging

// Option is a generic interface of the logger configuration option.
type Option interface {
	// Apply applies the option to the given logger.
	Apply(log *Logger)
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(log *Logger)

// Apply applies the option to the given logger.
func (o OptionFunc) Apply(log *Logger) {
	o(log)
}

// Hooks creates a new options that appends given hooks to the logger.
func Hooks(hooks ...Hook) OptionFunc {
	return func(log *Logger) {
		log.hooks = append(log.hooks, hooks...)
	}
}

// Fields creates a new option that appends given fields to the logger.
func Fields(fields ...Field) OptionFunc {
	return func(log *Logger) {
		log.core = log.core.With(fields)
	}
}

// Name creates a new options that sets the given name for the logger.
func Name(name string) OptionFunc {
	return func(log *Logger) {
		log.name = name
	}
}
