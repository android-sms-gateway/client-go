//nolint:lll // validator tags
package smsgateway

import (
	"fmt"
	"time"
)

type (
	// ProcessingState represents the state of a message.
	ProcessingState string
	// MessagePriority represents the priority of a message.
	MessagePriority int8
)

const (
	ProcessingStatePending   ProcessingState = "Pending"   // Pending
	ProcessingStateProcessed ProcessingState = "Processed" // Processed (received by device)
	ProcessingStateSent      ProcessingState = "Sent"      // Sent
	ProcessingStateDelivered ProcessingState = "Delivered" // Delivered
	ProcessingStateFailed    ProcessingState = "Failed"    // Failed

	PriorityMinimum         MessagePriority = -128
	PriorityDefault         MessagePriority = 0
	PriorityBypassThreshold MessagePriority = 100 // Threshold at which messages bypass limits and delays
	PriorityMaximum         MessagePriority = 127
)

//nolint:gochecknoglobals // lookup table
var allProcessStates = map[ProcessingState]struct{}{
	ProcessingStatePending:   {},
	ProcessingStateProcessed: {},
	ProcessingStateSent:      {},
	ProcessingStateDelivered: {},
	ProcessingStateFailed:    {},
}

// TextMessage represents an SMS message with a text body.
//
// Text is the message text.
type TextMessage struct {
	// Text is the message text.
	Text string `json:"text" validate:"required,min=1,max=65535" example:"Hello World!"`
}

// DataMessage represents an SMS message with a binary payload.
//
// Data is the base64-encoded payload.
//
// Port is the destination port.
type DataMessage struct {
	// Data is the base64-encoded payload.
	Data string `json:"data" validate:"required,base64,min=4,max=65535" example:"SGVsbG8gV29ybGQh" format:"byte"`
	// Port is the destination port.
	Port uint16 `json:"port" validate:"required,min=1,max=65535"        example:"53739"`
}

// Message represents an SMS message.
//
// ID is the message ID (if not set - will be generated).
// DeviceID is the optional device ID for explicit selection.
// Message is the message content (deprecated, use TextMessage instead).
// TextMessage is the text message.
// DataMessage is the data message.
// PhoneNumbers is the list of phone numbers.
// IsEncrypted is true if the message is encrypted.
// SimNumber is the SIM card number (1-3), if not set - default SIM will be used.
// WithDeliveryReport is true if the message should request a delivery report.
// Priority is the priority of the message, messages with values greater than `99` will bypass limits and delays.
// TTL is the time to live in seconds (conflicts with `ValidUntil`).
// ValidUntil is the time until the message is valid (conflicts with `TTL`).
type Message struct {
	ID       string `json:"id,omitempty"       validate:"omitempty,max=36" example:"PyDmBQZZXYmyxMwED8Fzy"` // ID (if not set - will be generated)
	DeviceID string `json:"deviceId,omitempty" validate:"omitempty,max=21" example:"PyDmBQZZXYmyxMwED8Fzy"` // Optional device ID for explicit selection

	Message string `json:"message,omitempty" validate:"omitempty,max=65535" example:"Hello World!"` // Message content (deprecated, use TextMessage instead)

	TextMessage *TextMessage `json:"textMessage,omitempty" validate:"omitempty"` // Text message
	DataMessage *DataMessage `json:"dataMessage,omitempty" validate:"omitempty"` // Data message

	PhoneNumbers []string `json:"phoneNumbers"          validate:"required,min=1,max=100,dive,required,min=1,max=128" example:"79990001234"` // Recipients (phone numbers)
	IsEncrypted  bool     `json:"isEncrypted,omitempty"                                                               example:"true"`        // Is encrypted

	SimNumber          *uint8          `json:"simNumber,omitempty"          validate:"omitempty,max=3"            example:"1"`                // SIM card number (1-3), if not set - default SIM will be used
	WithDeliveryReport *bool           `json:"withDeliveryReport,omitempty"                                       example:"true"`             // With delivery report
	Priority           MessagePriority `json:"priority,omitempty"           validate:"omitempty,min=-128,max=127" example:"0"    default:"0"` // Priority, messages with values greater than `99` will bypass limits and delays

	TTL        *uint64    `json:"ttl,omitempty"        validate:"omitempty,min=5" example:"86400"`                // Time to live in seconds (conflicts with `ValidUntil`)
	ValidUntil *time.Time `json:"validUntil,omitempty"                            example:"2020-01-01T00:00:00Z"` // Valid until (conflicts with `TTL`)
}

// GetTextMessage returns the TextMessage, if it was set explicitly, or
// constructs it from the deprecated Message field and returns it.
// If neither TextMessage nor Message are set, returns nil.
func (m *Message) GetTextMessage() *TextMessage {
	if m.TextMessage != nil {
		return m.TextMessage
	}

	if m.Message != "" {
		return &TextMessage{
			Text: m.Message,
		}
	}

	return nil
}

// GetDataMessage returns the DataMessage, if it was set explicitly, or nil otherwise.
func (m *Message) GetDataMessage() *DataMessage {
	return m.DataMessage
}

// Validate validates the Message structure.
func (m *Message) Validate() error {
	fields := []bool{
		m.Message != "",
		m.TextMessage != nil,
		m.DataMessage != nil,
	}

	filled := 0
	for _, f := range fields {
		if f {
			filled++
		}
	}

	if filled == 0 {
		return fmt.Errorf("%w: must specify exactly one of: textMessage or dataMessage", ErrValidationFailed)
	}
	if filled > 1 {
		return fmt.Errorf("%w: must specify exactly one of: textMessage or dataMessage", ErrConflictFields)
	}

	if m.TTL != nil && m.ValidUntil != nil {
		return fmt.Errorf("%w: ttl and validUntil", ErrConflictFields)
	}

	return nil
}

// MessageState represents the state of a message.
//
// MessageState is a struct used to communicate the state of a message
// between the client and the server. It contains the message ID, device ID,
// state, and hashed and encrypted flags. Additionally, it contains a slice
// of RecipientState, representing the state of each recipient, and a map
// of states, representing the history of states for the message.
type MessageState struct {
	ID          string               `json:"id"          validate:"required,max=36"     example:"PyDmBQZZXYmyxMwED8Fzy"` // Message ID
	DeviceID    string               `json:"deviceId"    validate:"required,max=21"     example:"PyDmBQZZXYmyxMwED8Fzy"` // Device ID
	State       ProcessingState      `json:"state"       validate:"required"            example:"Pending"`               // State
	IsHashed    bool                 `json:"isHashed"                                   example:"false"`                 // Hashed
	IsEncrypted bool                 `json:"isEncrypted"                                example:"false"`                 // Encrypted
	Recipients  []RecipientState     `json:"recipients"  validate:"required,min=1,dive"`                                 // Recipients states
	States      map[string]time.Time `json:"states"`                                                                     // History of states
}

func (m MessageState) Validate() error {
	for k := range m.States {
		if _, ok := allProcessStates[ProcessingState(k)]; !ok {
			return fmt.Errorf("%w: invalid state value: %s", ErrValidationFailed, k)
		}
	}

	return nil
}

// RecipientState represents the state of a recipient.
//
// RecipientState is a struct used to communicate the state of a recipient
// between the client and the server. It contains the phone number or first 16
// symbols of the SHA256 hash, state, and error information.
type RecipientState struct {
	PhoneNumber string          `json:"phoneNumber"     validate:"required,min=1,max=128" example:"79990001234"` // Phone number or first 16 symbols of SHA256 hash
	State       ProcessingState `json:"state"           validate:"required"               example:"Pending"`     // State
	Error       *string         `json:"error,omitempty"                                   example:"timeout"`     // Error (for `Failed` state)
}
