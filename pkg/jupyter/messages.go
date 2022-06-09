package jupyter

import (
	"bytes"
	"errors"
	"fmt"
)

// Message contains details about a jupyter message.
type Message interface {
	GetMsgType() MsgType
	IsChildByParentMsgID(u string) bool
}

// MessageBase contains base structure of the message.
type MessageBase struct {
	Header       Header       `json:"header"`
	ParentHeader ParentHeader `json:"parent_header"`
}

// GetMsgType returns header of the message
func (m MessageBase) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageBase) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// Header contains the structure of the message header.
type Header struct {
	MsgID string `json:"msg_id"`
	// Session  string     `json:"session"`
	// Username string `json:"username"`
	// Date     *time.Time `json:"date"`
	MsgType MsgType `json:"msg_type"`
	Version string  `json:"version"`
}

// ParentHeader contains the structure of the message header.
type ParentHeader struct {
	MsgID   string  `json:"msg_id"`
	MsgType MsgType `json:"msg_type"`
}

// MetaData contains the structure of the message metadata.
type MetaData struct{}

// MessageDisplayData contains details about a jupyter message
// for MsgTypeDisplayData.
type MessageDisplayData struct {
	Header       Header                    `json:"header"`
	ParentHeader ParentHeader              `json:"parent_header"`
	MetaData     MetaData                  `json:"metadata"`
	Content      MessageDisplayDataContent `json:"content"`
	Channel      Channel                   `json:"channel"`
}

// MessageDisplayDataContent contains the structure
// of the MsgTypeDisplayData message content.
type MessageDisplayDataContent struct {
	Meta any               `json:"metadata,omitempty"`
	Data map[string][]byte `json:"data"`
}

// GetMsgType returns header of the message
func (m MessageDisplayData) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageDisplayData) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MessageExecuteRequest contains details about a jupyter message
// for MsgTypeExecuteRequest.
type MessageExecuteRequest struct {
	Header       Header                       `json:"header"`
	ParentHeader ParentHeader                 `json:"parent_header"`
	MetaData     MetaData                     `json:"metadata"`
	Content      MessageExecuteRequestContent `json:"content"`
}

// MessageExecuteRequestContent contains the structure
// of the MsgTypeExecuteRequest message content.
type MessageExecuteRequestContent struct {
	// The code to execute
	Code string `json:"code"`
	// Whether to execute the code as quietly as possible
	// The default is `false`
	Silent bool `json:"silent"`
	// Whether to store history of the execution
	// The default `true` if silent is False
	// It is forced to  `false ` if silent is `true`
	StoreHistory bool `json:"store_history"`
	// Whether to allow stdin requests.
	// The default is `true`.
	AllowStdin bool `json:"allow_stdin"`
	// Whether to the abort execution queue on an error
	// The default is `false`
	StopOnError bool `json:"stop_on_error"`

	// A mapping of names to expressions to be evaluated.
	// UserExpressions json `json:"user_expressions"`
}

// GetMsgType returns header of the message
func (m MessageExecuteRequest) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageExecuteRequest) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MessageStream contains details about a jupyter message
// for MsgTypeStream.
type MessageStream struct {
	Header       Header               `json:"header"`
	ParentHeader ParentHeader         `json:"parent_header"`
	MetaData     MetaData             `json:"metadata"`
	Content      MessageStreamContent `json:"content"`
	Channel      Channel              `json:"channel"`
}

// MessageStreamContent contains the structure
// of the MsgTypeStream message content.
type MessageStreamContent struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// GetMsgType returns header of the message
func (m MessageStream) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageStream) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MessageError contains details about a jupyter message
// for MsgTypeError.
type MessageError struct {
	Header       Header              `json:"header"`
	ParentHeader ParentHeader        `json:"parent_header"`
	MetaData     MetaData            `json:"metadata"`
	Content      MessageErrorContent `json:"content"`
	Channel      Channel             `json:"channel"`
}

// MessageErrorContent contains the structure
// of the MsgTypeError message content.
type MessageErrorContent struct {
	EName     string   `json:"ename"`
	EValue    string   `json:"evalue"`
	Traceback []string `json:"traceback"`
}

// GetMsgType returns header of the message
func (m MessageError) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageError) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MessageExecuteReply contains details about a jupyter message
// for MsgTypeExecuteReply.
type MessageExecuteReply struct {
	Header       Header                     `json:"header"`
	ParentHeader ParentHeader               `json:"parent_header"`
	MetaData     MetaData                   `json:"metadata"`
	Content      MessageExecuteReplyContent `json:"content"`
	Channel      Channel                    `json:"channel"`
}

// MessageExecuteReplyContent contains the structure
// of the MsgTypeExecuteReply message content.
type MessageExecuteReplyContent struct {
	Status         string `json:"status"`
	ExecutionCount int    `json:"execution_count"`
}

// GetMsgType returns header of the message
func (m MessageExecuteReply) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageExecuteReply) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MessageStatus contains details about a jupyter message
// for MsgTypeStatus.
type MessageStatus struct {
	Header       Header               `json:"header"`
	ParentHeader ParentHeader         `json:"parent_header"`
	MetaData     MetaData             `json:"metadata"`
	Content      MessageStatusContent `json:"content"`
	Channel      Channel              `json:"channel"`
}

// MessageStatusContent contains the structure
// of the MsgTypeStatus message content.
type MessageStatusContent struct {
	ExecutionState State `json:"execution_state"`
}

// GetMsgType returns header of the message
func (m MessageStatus) GetMsgType() MsgType {
	return m.Header.MsgType
}

// IsChildByParentMsgID determines whether the message
// is child of given message ID or not
func (m MessageStatus) IsChildByParentMsgID(u string) bool {
	parentMsgID := u
	return parentMsgID == m.ParentHeader.MsgID
}

// MsgType represents a msgType of a message.
type MsgType uint

// Well-known msgTypes.
const (
	MsgTypeNone           MsgType = iota
	MsgTypeExecuteRequest         // +
	MsgTypeExecuteReply           // +
	MsgTypeInspectRequest
	MsgTypeInspectReply
	MsgTypeCompleteRequest
	MsgTypeCompleteReply
	MsgTypeHistoryRequest
	MsgTypeHistoryReply
	MsgTypeIsCompleteRequest
	MsgTypeIsCompleteReply
	MsgTypeConnectRequest
	MsgTypeConnectReply
	MsgTypeCommInfoRequest
	MsgTypeCommInfoReply
	MsgTypeKernelInfoRequest
	MsgTypeKernelInfoReply
	MsgTypeShutdownRequest
	MsgTypeShutdownReply
	MsgTypeInterruptRequest
	MsgTypeInterruptReply
	MsgTypeDebugRequest
	MsgTypeDebugReply
	MsgTypeStream
	MsgTypeDisplayData
	MsgTypeUpdateDisplayData
	MsgTypeExecuteInput
	MsgTypeExecuteResult
	MsgTypeError
	MsgTypeStatus
	MsgTypeClearOutput
	MsgTypeDebugEvent
	MsgTypeInputRequest
	MsgTypeInputReply
	MsgTypeCommMsg
	MsgTypeCommClose
	msgTypeCount
)

var msgTypeOutput = [...]string{
	"",
	"execute_request",
	"execute_reply",
	"inspect_request",
	"inspect_reply",
	"complete_request",
	"complete_reply",
	"history_request",
	"history_reply",
	"is_complete_request",
	"is_complete_reply",
	"connect_request",
	"connect_reply",
	"comm_info_request",
	"comm_info_reply",
	"kernel_info_request",
	"kernel_info_reply",
	"shutdown_request",
	"shutdown_reply",
	"interrupt_request",
	"interrupt_reply",
	"debug_request",
	"debug_reply",
	"stream",
	"display_data",
	"update_display_data",
	"execute_input",
	"execute_result",
	"error",
	"status",
	"clear_output",
	"debug_event",
	"input_request",
	"input_reply",
	"comm_msg",
	"comm_close",
	"invalid",
}

// String returns a string form of the msgType.
func (s MsgType) String() string {
	if s >= msgTypeCount {
		return fmt.Sprintf("%s (%d)", msgTypeOutput[msgTypeCount], s)
	}
	return msgTypeOutput[s]
}

// ErrMsgTypeInvalid is returns when msgType is invalid.
var ErrMsgTypeInvalid = errors.New("invalid msgType")

// MarshalText encodes the msgType to the text form.
func (s MsgType) MarshalText() ([]byte, error) {
	if s >= msgTypeCount {
		return nil, ErrMsgTypeInvalid
	}
	return []byte(msgTypeOutput[s]), nil
}

var msgTypeInput = map[string]MsgType{
	"execute_request":     MsgTypeExecuteRequest,
	"execute_reply":       MsgTypeExecuteReply,
	"inspect_request":     MsgTypeInspectRequest,
	"inspect_reply":       MsgTypeInspectReply,
	"complete_request":    MsgTypeCompleteRequest,
	"complete_reply":      MsgTypeCompleteReply,
	"history_request":     MsgTypeHistoryRequest,
	"history_reply":       MsgTypeHistoryReply,
	"is_complete_request": MsgTypeIsCompleteRequest,
	"is_complete_reply":   MsgTypeIsCompleteReply,
	"connect_request":     MsgTypeConnectRequest,
	"connect_reply":       MsgTypeConnectReply,
	"comm_info_request":   MsgTypeCommInfoRequest,
	"comm_info_reply":     MsgTypeCommInfoReply,
	"kernel_info_request": MsgTypeKernelInfoRequest,
	"kernel_info_reply":   MsgTypeKernelInfoReply,
	"shutdown_request":    MsgTypeShutdownRequest,
	"shutdown_reply":      MsgTypeShutdownReply,
	"interrupt_request":   MsgTypeInterruptRequest,
	"interrupt_reply":     MsgTypeInterruptReply,
	"debug_request":       MsgTypeDebugRequest,
	"debug_reply":         MsgTypeDebugReply,
	"stream":              MsgTypeStream,
	"display_data":        MsgTypeDisplayData,
	"update_display_data": MsgTypeUpdateDisplayData,
	"execute_input":       MsgTypeExecuteInput,
	"execute_result":      MsgTypeExecuteResult,
	"error":               MsgTypeError,
	"status":              MsgTypeStatus,
	"clear_output":        MsgTypeClearOutput,
	"debug_event":         MsgTypeDebugEvent,
	"input_request":       MsgTypeInputRequest,
	"input_reply":         MsgTypeInputReply,
	"comm_msg":            MsgTypeCommMsg,
	"comm_close":          MsgTypeCommClose,
}

// UnmarshalText decodes the text form of the msgType.
func (s *MsgType) UnmarshalText(text []byte) error {
	msgType, ok := msgTypeInput[string(bytes.ToLower(text))]
	if !ok {
		return ErrMsgTypeInvalid
	}
	*s = msgType
	return nil
}

// Channel represents a channel of a message.
type Channel uint

// Well-known channels.
const (
	ChannelNone Channel = iota
	ChannelShell
	ChannelIOPub
	ChannelStdin
	ChannelControl
	ChannelHeartbeat
	channelCount
)

var channelOutput = [...]string{
	"none",
	"shell",
	"iopub",
	"stdin",
	"control",
	"heartbeat",
	"invalid",
}

// String returns a string form of the channel.
func (s Channel) String() string {
	if s >= channelCount {
		return fmt.Sprintf("%s (%d)", channelOutput[channelCount], s)
	}
	return channelOutput[s]
}

// ErrChannelInvalid is returns when channel is invalid.
var ErrChannelInvalid = errors.New("invalid channel")

// MarshalText encodes the channel to the text form.
func (s Channel) MarshalText() ([]byte, error) {
	if s >= channelCount {
		return nil, ErrChannelInvalid
	}
	return []byte(channelOutput[s]), nil
}

var channelInput = map[string]Channel{
	"shell":     ChannelShell,
	"iopub":     ChannelIOPub,
	"stdin":     ChannelStdin,
	"control":   ChannelControl,
	"heartbeat": ChannelHeartbeat,
}

// UnmarshalText decodes the text form of the channel.
func (s *Channel) UnmarshalText(text []byte) error {
	channel, ok := channelInput[string(bytes.ToLower(text))]
	if !ok {
		return ErrChannelInvalid
	}
	*s = channel
	return nil
}

// State represents an execution state of a kernel.
type State uint

// Well-known execution states.
const (
	StateNone State = iota
	StateStarting
	StateIdle
	StateBusy
	StateTerminating
	StateRestarting
	StateAutoRestarting
	StateDead
	StateConnected
	StateConnecting
	StateDisconnected
	StateInitializing
	stateCount
)

var stateOutput = [...]string{
	"none",
	"starting",
	"idle",
	"busy",
	"terminating",
	"restarting",
	"autorestarting",
	"dead",
	"connected",
	"connecting",
	"disconnected",
	"initializing",
	"invalid",
}

// String returns a string form of the execution state.
func (s State) String() string {
	if s >= stateCount {
		return fmt.Sprintf("%s (%d)", stateOutput[stateCount], s)
	}
	return stateOutput[s]
}

// ErrStateInvalid is returns when execution state is invalid.
var ErrStateInvalid = errors.New("invalid state")

// MarshalText encodes the execution state to the text form.
func (s State) MarshalText() ([]byte, error) {
	if s >= stateCount {
		return nil, ErrStateInvalid
	}
	return []byte(stateOutput[s]), nil
}

var stateInput = map[string]State{
	"starting":       StateStarting,
	"idle":           StateIdle,
	"busy":           StateBusy,
	"terminating":    StateTerminating,
	"restarting":     StateRestarting,
	"autorestarting": StateAutoRestarting,
	"dead":           StateDead,
	"connected":      StateConnected,
	"connecting":     StateConnecting,
	"disconnected":   StateDisconnected,
	"initializing":   StateInitializing,
}

// UnmarshalText decodes the text form of the execution state.
func (s *State) UnmarshalText(text []byte) error {
	state, ok := stateInput[string(bytes.ToLower(text))]
	if !ok {
		return ErrStateInvalid
	}
	*s = state
	return nil
}
