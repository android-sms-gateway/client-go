package smsgateway_test

import (
	"errors"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestMessage_GetTextMessage(t *testing.T) {
	tests := []struct {
		name         string
		message      smsgateway.Message
		expectedText *smsgateway.TextMessage
	}{
		{
			name: "TextMessage field set",
			message: smsgateway.Message{
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				PhoneNumbers: []string{"1234567890"},
			},
			expectedText: &smsgateway.TextMessage{
				Text: "Hello World!",
			},
		},
		{
			name: "Message field set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				PhoneNumbers: []string{"1234567890"},
			},
			expectedText: &smsgateway.TextMessage{
				Text: "Hello World!",
			},
		},
		{
			name: "Neither TextMessage nor Message set",
			message: smsgateway.Message{
				PhoneNumbers: []string{"1234567890"},
			},
			expectedText: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.GetTextMessage()

			if result == nil && tt.expectedText != nil {
				t.Errorf("GetTextMessage() = nil, expected %v", tt.expectedText)
			}

			if result != nil && tt.expectedText == nil {
				t.Errorf("GetTextMessage() = %v, expected nil", result)
				return
			}

			if result != nil && tt.expectedText != nil && result.Text != tt.expectedText.Text {
				t.Errorf("GetTextMessage() = %v, expected %v", result, tt.expectedText)
			}
		})
	}
}

func TestMessage_GetDataMessage(t *testing.T) {
	tests := []struct {
		name         string
		message      smsgateway.Message
		expectedData *smsgateway.DataMessage
	}{
		{
			name: "DataMessage field set",
			message: smsgateway.Message{
				DataMessage: &smsgateway.DataMessage{
					Data: "SGVsbG8gV29ybGQh",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			expectedData: &smsgateway.DataMessage{
				Data: "SGVsbG8gV29ybGQh",
				Port: 1,
			},
		},
		{
			name: "DataMessage field not set",
			message: smsgateway.Message{
				PhoneNumbers: []string{"1234567890"},
			},
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.GetDataMessage()

			if result == nil && tt.expectedData != nil {
				t.Errorf("GetDataMessage() = nil, expected %v", tt.expectedData)
			}

			if result != nil && tt.expectedData == nil {
				t.Errorf("GetDataMessage() = %v, expected nil", result)
				return
			}

			if result != nil && tt.expectedData != nil {
				if result.Data != tt.expectedData.Data || result.Port != tt.expectedData.Port {
					t.Errorf("GetDataMessage() = %v, expected %v", result, tt.expectedData)
				}
			}
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message smsgateway.Message
		err     error
	}{
		{
			name: "Valid - only Message field set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Valid - only TextMessage field set",
			message: smsgateway.Message{
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Valid - only DataMessage field set",
			message: smsgateway.Message{
				DataMessage: &smsgateway.DataMessage{
					Data: "SGVsbG8gV29ybGQh",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Invalid - no message fields set",
			message: smsgateway.Message{
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrValidationFailed,
		},
		{
			name: "Invalid - multiple message fields set (Message + TextMessage)",
			message: smsgateway.Message{
				Message: "Hello World!",
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
		{
			name: "Invalid - multiple message fields set (Message + DataMessage)",
			message: smsgateway.Message{
				Message: "Hello World!",
				DataMessage: &smsgateway.DataMessage{
					Data: "SGVsbG8gV29ybGQh",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
		{
			name: "Invalid - multiple message fields set (TextMessage + DataMessage)",
			message: smsgateway.Message{
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				DataMessage: &smsgateway.DataMessage{
					Data: "SGVsbG8gV29ybGQh",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
		{
			name: "Invalid - all message fields set",
			message: smsgateway.Message{
				Message: "Hello World!",
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				DataMessage: &smsgateway.DataMessage{
					Data: "SGVsbG8gV29ybGQh",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
		{
			name: "Edge case - empty Message field",
			message: smsgateway.Message{
				Message:      "",
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrValidationFailed, // Empty string is treated as field not set
		},
		{
			name: "Edge case - empty TextMessage field",
			message: smsgateway.Message{
				TextMessage: &smsgateway.TextMessage{
					Text: "",
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil, // Empty text is valid, validation is only for field presence
		},
		{
			name: "Edge case - empty DataMessage field",
			message: smsgateway.Message{
				DataMessage: &smsgateway.DataMessage{
					Data: "",
					Port: 1,
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil, // Empty data is valid, validation is only for field presence
		},
		{
			name: "Valid - neither TTL nor ValidUntil set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Valid - only TTL set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				TTL:          func() *uint64 { val := uint64(3600); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Valid - only ValidUntil set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				ValidUntil:   func() *time.Time { val := time.Now().Add(time.Hour); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Invalid - both TTL and ValidUntil set",
			message: smsgateway.Message{
				Message:      "Hello World!",
				TTL:          func() *uint64 { val := uint64(3600); return &val }(),
				ValidUntil:   func() *time.Time { val := time.Now().Add(time.Hour); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()

			if tt.err == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, expected no error", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, expected error")
					return
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("Validate() error = %v, want %v", err, tt.err)
				}
			}
		})
	}
}

func TestMessageState_Validate(t *testing.T) {
	tests := []struct {
		name    string
		states  map[string]time.Time
		wantErr bool
	}{
		{
			name:    "Empty states",
			states:  map[string]time.Time{},
			wantErr: false,
		},
		{
			name: "Valid states",
			states: map[string]time.Time{
				string(smsgateway.ProcessingStatePending):   time.Now(),
				string(smsgateway.ProcessingStateProcessed): time.Now(),
				string(smsgateway.ProcessingStateSent):      time.Now(),
				string(smsgateway.ProcessingStateDelivered): time.Now(),
				string(smsgateway.ProcessingStateFailed):    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Invalid state",
			states: map[string]time.Time{
				string(smsgateway.ProcessingStatePending): time.Now(),
				"InvalidState": time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := smsgateway.MessageState{
				States: tt.states,
			}

			err := m.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}

				if !errors.Is(err, smsgateway.ErrValidationFailed) {
					t.Errorf("Validate() error = %v, want error type %v", err, smsgateway.ErrValidationFailed)
				}
			} else if err != nil {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
