package sandbox

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"go.uber.org/multierr"

	"github.com/uclatall/ckhub/pkg/jupyter"
)

// Kernel is a thin wrapper around a jupyter kernel that provides access to
// the kernel metdata.
type Kernel struct {
	client *jupyter.Client

	name  string
	init  string
	limit int64

	mu        sync.RWMutex
	close     bool
	instances []*jupyter.Kernel
	total     int64
}

// ErrKernelClosed is returned when the kernel is closed.
var ErrKernelClosed = errors.New("kernel closed")

// SpawnInstance creates a new instance of the kernel.
func (k *Kernel) SpawnInstance(ctx context.Context) error {
	if atomic.LoadInt64(&k.total) >= k.limit {
		return nil
	}

	k.mu.RLock()
	if k.close {
		k.mu.RUnlock()
		return ErrKernelClosed
	}
	k.mu.RUnlock()

	atomic.AddInt64(&k.total, 1)

	kernel, err := k.client.CreateKernel(ctx, k.name)
	if err != nil {
		atomic.AddInt64(&k.total, -1)
		return fmt.Errorf("failed to create kernel: %w", err)
	}

	if k.init != "" {
		_, err := executeCode(kernel, uuid.New(), k.init)
		if err != nil {
			_ = k.client.RemoveKernel(ctx, kernel)
			atomic.AddInt64(&k.total, -1)
			return fmt.Errorf("failed to init kernel: %w", err)
		}
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.close {
		_ = k.client.RemoveKernel(ctx, kernel)
		atomic.AddInt64(&k.total, -1)
		return ErrKernelClosed
	}

	k.instances = append(k.instances, kernel)

	return nil
}

// ErrTooManyRequests is returned when the kernel is at its limit.
var ErrTooManyRequests = errors.New("too many requests")

// ExecuteSnippet executes the given snippet.
func (k *Kernel) ExecuteSnippet(
	ctx context.Context,
	snippet *Snippet,
) (*Result, error) {

	k.mu.Lock()

	if k.close {
		return nil, ErrKernelClosed
	}

	if len(k.instances) == 0 {
		k.mu.Unlock()
		return nil, ErrTooManyRequests
	}

	kernel := k.instances[0]
	k.instances = k.instances[1:]

	go func() {
		_ = k.SpawnInstance(context.Background())
	}()

	k.mu.Unlock()

	result, err := executeCode(kernel, snippet.ID, snippet.Source)

	go func() {
		_ = k.client.RemoveKernel(context.Background(), kernel)
		atomic.AddInt64(&k.total, -1)
	}()

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Destroy destroys all kernel instances.
func (k *Kernel) Destroy() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.close = true

	errs := make([]error, len(k.instances))
	for i, kernel := range k.instances {
		errs[i] = k.client.RemoveKernel(context.Background(), kernel)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return err
	}
	return nil
}

// Instances returns the number of kernel instances.
func (k *Kernel) Instances() int {
	return int(atomic.LoadInt64(&k.total))
}

func executeCode(kernel *jupyter.Kernel, id uuid.UUID, code string) (*Result, error) {
	err := kernel.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kernel: %w", err)
	}
	defer func() { _ = kernel.Close() }()

	err = kernel.Execute(id, code)
	if err != nil {
		return nil, fmt.Errorf("failed to execute code: %w", err)
	}

	result := &Result{}

loop:
	for {
		msg, err := kernel.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		if !msg.IsChildByParentMsgID(id.String()) {
			continue
		}

		switch msg := msg.(type) {
		case *jupyter.MessageDisplayData:
			for mime, content := range msg.Content.Data {
				if strings.Split(mime, "/")[0] == "text" {
					result.Events = append(result.Events, Event{
						Kind:    EventKindOutput,
						Message: string(content),
					})
					continue
				}
				result.Events = append(result.Events, Event{
					Kind:    EventKindPayload,
					Mime:    mime,
					Message: base64.StdEncoding.EncodeToString(content),
				})
			}
		case *jupyter.MessageError:
			result.Errors = append(result.Errors, Error{
				Message: msg.Content.EValue,
			})
		case *jupyter.MessageExecuteReply:
			result.Status = msg.Content.Status
		case *jupyter.MessageStream:
			switch msg.Content.Name {
			case "stderr":
				result.Events = append(result.Events, Event{
					Kind:    EventKindError,
					Message: msg.Content.Text,
				})
			case "stdout":
				result.Events = append(result.Events, Event{
					Kind:    EventKindOutput,
					Message: msg.Content.Text,
				})
			default:
				result.Events = append(result.Events, Event{
					Kind:    EventKindNone,
					Message: msg.Content.Text,
				})
			}
		case *jupyter.MessageStatus:
			if msg.Content.ExecutionState == jupyter.StateIdle {
				break loop
			}
		}
	}

	return result, nil
}
