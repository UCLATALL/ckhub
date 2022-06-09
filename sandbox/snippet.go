package sandbox

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Snippet represents a snippet to be executed in the sandbox.
type Snippet struct {
	ID     uuid.UUID
	Kernel string
	Source string
}

// Result represents a snippet execution result.
type Result struct {
	ID     uuid.UUID `json:"-"`
	Status string    `json:"status,omitempty"`
	Errors []Error   `json:"errors,omitempty"`
	Events []Event   `json:"events,omitempty"`
}

// Error represents a snippet execution error.
type Error struct {
	Message string `json:"message"`
}

// Event represents a snippet execution event.
type Event struct {
	Kind    EventKind `json:"kind"`
	Mime    string    `json:"mime,omitempty"`
	Message string    `json:"message,omitempty"`
}

// EventKind represetns an event kind.
type EventKind uint

// Well-known event kinds.
const (
	EventKindNone EventKind = iota
	EventKindOutput
	EventKindError
	EventKindPayload
	eventKindCount
)

var eventKindOutput = []string{
	"none",
	"output",
	"error",
	"payload",
	"invalid",
}

// String returns a string form of the event kind.
func (kind EventKind) String() string {
	if kind >= eventKindCount {
		return fmt.Sprintf("%s (%d)", eventKindOutput[eventKindCount], kind)
	}
	return eventKindOutput[kind]
}

// ErrEventKindInvalid is returned when the event kind is invalid.
var ErrEventKindInvalid = errors.New("invalid event kind")

// MarshalText marshals event kind into text form.
func (kind EventKind) MarshalText() ([]byte, error) {
	if kind >= eventKindCount {
		return nil, ErrEventKindInvalid
	}
	return []byte(eventKindOutput[kind]), nil
}
