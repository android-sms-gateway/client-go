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

// Message
type Message struct {
	// ID (if not set - will be generated)
	ID string `json:"id,omitempty" validate:"omitempty,max=36" example:"PyDmBQZZXYmyxMwED8Fzy"`
	// Content
	Message string `json:"message" validate:"required,max=65535" example:"Hello World!"`
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

func (m Message) Validate() error {
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
