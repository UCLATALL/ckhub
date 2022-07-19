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
	ID      uuid.UUID `json:"-"`
	Status  string    `json:"status,omitempty"`
	Errors  []Error   `json:"errors,omitempty"`
	Outputs []Output  `json:"outputs,omitempty"`
}

// Error represents a snippet execution error.
type Error struct {
	Data any            `json:"data"`
	Meta map[string]any `json:"metadata"`
}

// Output represents a snippet execution output.
type Output struct {
	Kind OutputKind     `json:"type"`
	Data any            `json:"data"`
	Meta map[string]any `json:"metadata"`
}

// OutputKind represetns an output kind.
type OutputKind uint

// Well-known output kinds.
const (
	OutputKindNone OutputKind = iota
	OutputKindDisplayData
	OutputKindStream
	outputKindCount
)

var outputKindOutput = []string{
	"none",
	"display_data",
	"stream",
	"invalid",
}

// String returns a string form of the output kind.
func (kind OutputKind) String() string {
	if kind >= outputKindCount {
		return fmt.Sprintf("%s (%d)", outputKindOutput[outputKindCount], kind)
	}
	return outputKindOutput[kind]
}

// ErrEventKindInvalid is returned when the output kind is invalid.
var ErrEventKindInvalid = errors.New("invalid output")

// MarshalText marshals output kind into text form.
func (kind OutputKind) MarshalText() ([]byte, error) {
	if kind >= outputKindCount {
		return nil, ErrEventKindInvalid
	}
	return []byte(outputKindOutput[kind]), nil
}
