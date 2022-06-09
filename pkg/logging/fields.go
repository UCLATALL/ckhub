package logging

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

// Field represents a field of the logging structured context.
type Field = zapcore.Field

// Well-known field names.
const (
	FieldCaller  = "caller"
	FieldError   = "error"
	FieldLevel   = "level"
	FieldMessage = "message"
	FieldName    = "logger"
	FieldSpan    = "duration"
	FieldTime    = "time"
	FieldTrace   = "trace"
)

// Bool creates a new field with the given key and boolean value.
func Bool(key string, b bool) Field {
	if b {
		return Field{Key: key, Type: zapcore.BoolType, Integer: 1}
	}
	return Field{Key: key, Type: zapcore.BoolType, Integer: 0}
}

// Bytes creates a new field with the given key and bytes value.
func Bytes(key string, p []byte) Field {
	buf := make([]byte, 0, len(p)*2)
	for i := 0; i < len(p); i++ {
		strconv.AppendInt(p, int64(p[i]), 16)
	}
	return Field{Key: key, Type: zapcore.StringType, String: string(buf)}
}

// Duration creates a new field with the given key and duration.
func Duration(key string, d time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(d)}
}

// Error creates a new field with the given key and error.
func Error(err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: FieldError, Type: zapcore.ErrorType, Interface: err}
}

// Float32 creates a new field with the given key and float32 value.
func Float32(key string, f float32) Field {
	return Field{
		Key:     key,
		Type:    zapcore.Float32Type,
		Integer: int64(math.Float32bits(f)),
	}
}

// Float64 creates a new field with the given key and float64 value.
func Float64(key string, f float64) Field {
	return Field{
		Key:     key,
		Type:    zapcore.Float64Type,
		Integer: int64(math.Float64bits(f)),
	}
}

// Int creates a new field with the given key and integer number.
func Int(key string, i int) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: int64(i)}
}

// Int32 creates a new field with the given key and int32 value.
func Int32(key string, i int32) Field {
	return Field{Key: key, Type: zapcore.Int32Type, Integer: int64(i)}
}

// Int64 creates a new field with the given key and int64 value.
func Int64(key string, i int64) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: i}
}

// Skip creates a new no-op field, which is often useful when handling invalid
// inputs in other field constructors.
func Skip() Field {
	return Field{Type: zapcore.SkipType}
}

// String creates a new field with the given key and string.
func String(key, s string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: s}
}

// Stringer creates a new field with the given key and output of the value's
// string method. The string method is called lazily.
func Stringer(key string, v fmt.Stringer) Field {
	return Field{Key: key, Type: zapcore.StringerType, Interface: v}
}

// Stringf creates a new field with the given key and format specifier.
func Stringf(key, format string, a ...any) Field {
	return Field{
		Key:    key,
		Type:   zapcore.StringType,
		String: fmt.Sprintf(format, a...),
	}
}

var (
	minTime = time.Unix(0, math.MinInt64)
	maxTime = time.Unix(0, math.MaxInt64)
)

// Time creates a new field with the given key and time.
func Time(key string, t time.Time) Field {
	if t.Before(minTime) || t.After(maxTime) {
		return Field{Key: key, Type: zapcore.TimeFullType, Interface: t}
	}
	return Field{
		Key:       key,
		Type:      zapcore.TimeType,
		Integer:   t.UnixNano(),
		Interface: t.Location(),
	}
}

// Uint creates a new field with the given key and unsigned integer number.
func Uint(key string, i uint) Field {
	return Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(i)}
}

// Uint32 creates a new field with the given key and uint32 value.
func Uint32(key string, i uint32) Field {
	return Field{Key: key, Type: zapcore.Uint32Type, Integer: int64(i)}
}

// Uint64 creates a new field with the given key and uint64 value.
func Uint64(key string, i uint64) Field {
	return Field{Key: key, Type: zapcore.Uint64Type, Integer: int64(i)}
}
