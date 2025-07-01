//nolint:lll // validator tags
package smsgateway

import (
	"fmt"
	"time"
)

type (
	// Processing state
	ProcessingState string

	// Message priority
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

// Text SMS message
type TextMessage struct {
	// Message text
	Text string `json:"text" validate:"required,min=1,max=65535" example:"Hello World!"`
}

// Data SMS message
type DataMessage struct {
	// Base64-encoded payload
	Data string `json:"data" validate:"required,base64,min=4,max=65535" example:"SGVsbG8gV29ybGQh" format:"byte"`
	// Destination port
	Port uint16 `json:"port" validate:"required,min=1,max=65535" example:"53739"`
}

// Message
type Message struct {
	// ID (if not set - will be generated)
	ID string `json:"id,omitempty" validate:"omitempty,max=36" example:"PyDmBQZZXYmyxMwED8Fzy"`
	// Message content
	// Deprecated: use TextMessage instead
	Message string `json:"message,omitempty" validate:"omitempty,max=65535" example:"Hello World!"`

	// Text message
	TextMessage *TextMessage `json:"textMessage,omitempty" validate:"omitempty"`
	// Data message
	DataMessage *DataMessage `json:"dataMessage,omitempty" validate:"omitempty"`

	// Recipients (phone numbers)
	PhoneNumbers []string `json:"phoneNumbers" validate:"required,min=1,max=100,dive,required,min=1,max=128" example:"79990001234"`
	// Is encrypted
	IsEncrypted bool `json:"isEncrypted,omitempty" example:"true"`

	// SIM card number (1-3), if not set - default SIM will be used
	SimNumber *uint8 `json:"simNumber,omitempty" validate:"omitempty,max=3" example:"1"`
	// With delivery report
	WithDeliveryReport *bool `json:"withDeliveryReport,omitempty" example:"true"`
	// Priority, messages with values greater than `99` will bypass limits and delays
	Priority MessagePriority `json:"priority,omitempty" validate:"omitempty,min=-128,max=127" example:"0" default:"0"`

	// Time to live in seconds (conflicts with `validUntil`)
	TTL *uint64 `json:"ttl,omitempty" validate:"omitempty,min=5" example:"86400"`
	// Valid until (conflicts with `ttl`)
	ValidUntil *time.Time `json:"validUntil,omitempty" example:"2020-01-01T00:00:00Z"`
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

// Message state
type MessageState struct {
	// Message ID
	ID string `json:"id,omitempty" validate:"omitempty,max=36" example:"PyDmBQZZXYmyxMwED8Fzy"`
	// State
	State ProcessingState `json:"state" validate:"required" example:"Pending"`
	// Hashed
	IsHashed bool `json:"isHashed" example:"false"`
	// Encrypted
	IsEncrypted bool `json:"isEncrypted" example:"false"`
	// Recipients states
	Recipients []RecipientState `json:"recipients" validate:"required,min=1,dive"`
	// History of states
	States map[string]time.Time `json:"states"`
}

func (m MessageState) Validate() error {
	for k := range m.States {
		if _, ok := allProcessStates[ProcessingState(k)]; !ok {
			return fmt.Errorf("%w: invalid state value: %s", ErrValidationFailed, k)
		}
	}

	return nil
}

// Recipient state
type RecipientState struct {
	// Phone number or first 16 symbols of SHA256 hash
	PhoneNumber string `json:"phoneNumber" validate:"required,min=1,max=128" example:"79990001234"`
	// State
	State ProcessingState `json:"state" validate:"required" example:"Pending"`
	// Error (for `Failed` state)
	Error *string `json:"error,omitempty" example:"timeout"`
}
