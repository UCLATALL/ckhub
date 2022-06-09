package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger implements fast, leveled, structured logger.
type Logger struct {
	core  zapcore.Core
	clock zapcore.Clock

	level zapcore.Level
	name  string
	hooks []Hook
}

var encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
	MessageKey:     FieldMessage,
	LevelKey:       FieldLevel,
	TimeKey:        FieldTime,
	NameKey:        FieldName,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
})

// NewLogger creates a new logger with the given options.
func NewLogger(options ...Option) Logger {
	var log Logger

	log.clock = zapcore.DefaultClock
	log.core = zapcore.NewCore(encoder, os.Stderr, &log.level)
	log.level = zapcore.InfoLevel

	for _, option := range options {
		option.Apply(&log)
	}

	return log
}

// NopLogger creates a new logger that does nothing.
func NopLogger() Logger {
	return Logger{
		clock: zapcore.DefaultClock,
		core:  zapcore.NewNopCore(),
	}
}

// Output logs message at the given level. The message includes given fields,
// as well as any fields accumulated on the logger.
func (log Logger) Output(lvl Level, msg string, fields ...zap.Field) {
	entry := log.core.Check(zapcore.Entry{
		Level:      zapcore.Level(lvl),
		LoggerName: log.name,
		Message:    msg,
		Time:       log.clock.Now(),
	}, nil)
	if entry == nil {
		return
	}

	for _, hook := range log.hooks {
		hook.Hook(entry)
	}
	entry.Write(fields...)
}

// Debug logs message at the debug level. The message includes given fields,
// as well as any fields accumulated on the logger.
func (log Logger) Debug(msg string, fields ...Field) {
	log.Output(LevelDebug, msg, fields...)
}

// Info logs message at the info level. The message includes given fields,
// as well as any fields accumulated on the logger.
func (log Logger) Info(msg string, fields ...Field) {
	log.Output(LevelInfo, msg, fields...)
}

// Warn logs message at the warn level. The message includes given fields,
// as well as any fields accumulated on the logger.
func (log Logger) Warn(msg string, fields ...Field) {
	log.Output(LevelWarn, msg, fields...)
}

// Error logs message at the error level. The message includes given fields,
// as well as any fields accumulated on the logger.
func (log Logger) Error(msg string, fields ...Field) {
	log.Output(LevelError, msg, fields...)
}

// Fields creates a child logger and appends given fields to it.
func (log Logger) Fields(fields ...Field) Logger {
	log.core = log.core.With(fields)
	return log
}

// Hooks creates a child logger and appends given hooks to it.
func (log Logger) Hooks(hooks ...Hook) Logger {
	log.hooks = append(log.hooks, hooks...)
	return log
}

// Name creates a child logger set the given name to it.
func (log Logger) Name(name string) Logger {
	log.name = name
	return log
}
