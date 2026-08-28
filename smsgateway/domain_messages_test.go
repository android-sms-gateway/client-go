package smsgateway_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
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

func TestMessage_GetMmsMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     smsgateway.Message
		expectedMms *smsgateway.MmsMessage
	}{
		{
			name: "MmsMessage field set",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Text: ptr("World"),
				},
				PhoneNumbers: []string{"1234567890"},
			},
			expectedMms: &smsgateway.MmsMessage{
				Text: ptr("World"),
			},
		},
		{
			name: "MmsMessage field not set",
			message: smsgateway.Message{
				PhoneNumbers: []string{"1234567890"},
			},
			expectedMms: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.GetMmsMessage()

			if result == nil && tt.expectedMms != nil {
				t.Errorf("GetMmsMessage() = nil, expected %v", tt.expectedMms)
			}

			if result != nil && tt.expectedMms == nil {
				t.Errorf("GetMmsMessage() = %v, expected nil", result)
				return
			}

			if result != nil && tt.expectedMms != nil && !reflect.DeepEqual(result, tt.expectedMms) {
				t.Errorf("GetMmsMessage() = %v, expected %v", result, tt.expectedMms)
			}
		})
	}
}

func TestMmsMessage_Validate_NilReceiver(t *testing.T) {
	var m *smsgateway.MmsMessage

	if err := m.Validate(); err == nil {
		t.Errorf("Validate() error = nil, expected error for nil receiver")
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
		{
			name: "Valid - ScheduleAt set to future time",
			message: smsgateway.Message{
				Message:      "Hello World!",
				ScheduleAt:   func() *time.Time { val := time.Now().Add(time.Hour); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Invalid - ScheduleAt set to past time",
			message: smsgateway.Message{
				Message:      "Hello World!",
				ScheduleAt:   func() *time.Time { val := time.Now().Add(-time.Hour); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrValidationFailed,
		},
		{
			name: "Invalid - ScheduleAt set to current time",
			message: smsgateway.Message{
				Message:      "Hello World!",
				ScheduleAt:   func() *time.Time { val := time.Now(); return &val }(),
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrValidationFailed,
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

func TestMmsMessage_MarshalFixture(t *testing.T) {
	mms := smsgateway.MmsMessage{
		Subject: ptr("Hello"),
		Text:    ptr("World"),
		Attachments: []smsgateway.MmsAttachment{
			{
				ContentType: "image/png",
				Name:        ptr("picture.png"),
				Data:        "BASE64DATA",
			},
		},
	}

	expected := `{"subject":"Hello","text":"World","attachments":[{"contentType":"image/png","name":"picture.png","data":"BASE64DATA"}]}`

	got, err := json.Marshal(mms)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if string(got) != expected {
		t.Errorf("Marshal() = %s, want %s", got, expected)
	}
}

func TestMmsMessage_MarshalOmitEmptyAttachments(t *testing.T) {
	mms := smsgateway.MmsMessage{
		Subject: ptr("Hello"),
		Text:    ptr("World"),
	}

	got, err := json.Marshal(mms)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	if strings.Contains(string(got), "attachments") {
		t.Errorf("Marshal() = %s, should not contain \"attachments\" key", got)
	}

	expected := `{"subject":"Hello","text":"World"}`
	if string(got) != expected {
		t.Errorf("Marshal() = %s, want %s", got, expected)
	}
}

func TestMmsMessage_RoundTrip(t *testing.T) {
	original := smsgateway.MmsMessage{
		Subject: ptr("Hello"),
		Text:    ptr("World"),
		Attachments: []smsgateway.MmsAttachment{
			{
				ContentType: "image/png",
				Name:        ptr("picture.png"),
				Data:        "BASE64DATA",
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded smsgateway.MmsMessage
	if uerr := json.Unmarshal(data, &decoded); uerr != nil {
		t.Fatalf("Unmarshal() error = %v", uerr)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("RoundTrip() = %+v, want %+v", decoded, original)
	}
}

func TestMmsMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message smsgateway.MmsMessage
		wantErr bool
	}{
		{
			name: "Valid - text set",
			message: smsgateway.MmsMessage{
				Text: ptr("World"),
			},
			wantErr: false,
		},
		{
			name: "Valid - one attachment",
			message: smsgateway.MmsMessage{
				Attachments: []smsgateway.MmsAttachment{
					{ContentType: "image/png", Data: "BASE64DATA"},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid - neither text nor attachments",
			message: smsgateway.MmsMessage{
				Subject: ptr("Hello"),
			},
			wantErr: true,
		},
		{
			name: "Valid - attachment missing contentType (tag-based, server-side)",
			message: smsgateway.MmsMessage{
				Attachments: []smsgateway.MmsAttachment{
					{Data: "BASE64DATA"},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid - attachment missing data (tag-based, server-side)",
			message: smsgateway.MmsMessage{
				Attachments: []smsgateway.MmsAttachment{
					{ContentType: "image/png"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_Validate_MmsMessage(t *testing.T) {
	tests := []struct {
		name    string
		message smsgateway.Message
		err     error
	}{
		{
			name: "Valid - only MmsMessage (text)",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Text: ptr("World"),
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Valid - only MmsMessage (attachment)",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Attachments: []smsgateway.MmsAttachment{
						{ContentType: "image/png", Data: "BASE64DATA"},
					},
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: nil,
		},
		{
			name: "Invalid - MmsMessage with neither text nor attachments",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Subject: ptr("Hello"),
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrValidationFailed,
		},
		{
			name: "Invalid - MmsMessage combined with TextMessage",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Text: ptr("World"),
				},
				TextMessage: &smsgateway.TextMessage{
					Text: "Hello World!",
				},
				PhoneNumbers: []string{"1234567890"},
			},
			err: smsgateway.ErrConflictFields,
		},
		{
			name: "Invalid - MmsMessage combined with Message",
			message: smsgateway.Message{
				MmsMessage: &smsgateway.MmsMessage{
					Text: ptr("World"),
				},
				Message:      "Hello World!",
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
