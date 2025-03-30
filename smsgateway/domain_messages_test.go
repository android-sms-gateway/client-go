package smsgateway_test

import (
	"errors"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message smsgateway.Message
		err     error
	}{
		{
			name:    "Valid - neither TTL nor ValidUntil set",
			message: smsgateway.Message{},
			err:     nil,
		},
		{
			name: "Valid - only TTL set",
			message: smsgateway.Message{
				TTL: func() *uint64 { val := uint64(3600); return &val }(),
			},
			err: nil,
		},
		{
			name: "Valid - only ValidUntil set",
			message: smsgateway.Message{
				ValidUntil: func() *time.Time { val := time.Now().Add(time.Hour); return &val }(),
			},
			err: nil,
		},
		{
			name: "Invalid - both TTL and ValidUntil set",
			message: smsgateway.Message{
				TTL:        func() *uint64 { val := uint64(3600); return &val }(),
				ValidUntil: func() *time.Time { val := time.Now().Add(time.Hour); return &val }(),
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
