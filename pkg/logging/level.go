package logging

import "go.uber.org/zap/zapcore"

// Level represents a logging priority. Higher levels are more important.
type Level int8

// Well-known logging priority levels.
const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
)

// Apply applies the level to the logger.
func (lvl Level) Apply(log *Logger) {
	log.level = zapcore.Level(lvl)
}
