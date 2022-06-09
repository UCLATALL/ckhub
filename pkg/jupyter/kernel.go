package jupyter

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// Kernel contains details of the jupyter kernel.
type Kernel struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	ChanURL string    `json:"-"`

	mu   sync.RWMutex
	conn *websocket.Conn
}

// Connect establishes a connection to the jupyter kernel.
func (k *Kernel) Connect() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	conn, err := websocket.Dial(k.ChanURL, "", "http://localhost")
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	k.conn = conn

	return nil
}

// ErrNotConnected is returned when the kernel connection is not established.
var ErrNotConnected = errors.New("not connected")

// Execute executes the given code in the jupyter kernel.
func (k *Kernel) Execute(id uuid.UUID, code string) error {
	return k.WriteMessage(&MessageExecuteRequest{
		Header: Header{
			MsgID:   id.String(),
			MsgType: MsgTypeExecuteRequest,
		},
		Content: MessageExecuteRequestContent{
			Code: code,
		},
	})
}

var bufferPool = &sync.Pool{
	New: func() any {
		return make([]byte, 0, 4096)
	},
}

// ReadMessage reads a next message from the kernel.
func (k *Kernel) ReadMessage() (Message, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.conn == nil {
		return nil, ErrNotConnected
	}

	buf, _ := bufferPool.Get().([]byte)
	buf = buf[:0]
	//nolint:staticcheck // Slices is a pointer-like variables.
	defer bufferPool.Put(buf)

	err := websocket.Message.Receive(k.conn, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var base MessageBase
	err = json.Unmarshal(buf, &base)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	switch base.GetMsgType() {
	case MsgTypeStream:
		msg := new(MessageStream)
		err = json.Unmarshal(buf, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		return msg, nil
	case MsgTypeError:
		msg := new(MessageError)
		err = json.Unmarshal(buf, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		return msg, nil
	case MsgTypeExecuteReply:
		msg := new(MessageExecuteReply)
		err = json.Unmarshal(buf, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		return msg, nil
	case MsgTypeStatus:
		msg := new(MessageStatus)
		err = json.Unmarshal(buf, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		return msg, nil
	case MsgTypeDisplayData:
		msg := new(MessageDisplayData)
		err = json.Unmarshal(buf, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		return msg, nil
	default:
		return base, nil
	}
}

// WriteMessage writes given message to the kernel.
func (k *Kernel) WriteMessage(msg Message) error {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if k.conn == nil {
		return ErrNotConnected
	}

	err := json.NewEncoder(k.conn).Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// Close closes the connection to the jupyter kernel.
func (k *Kernel) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.conn == nil {
		return nil
	}

	err := k.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
