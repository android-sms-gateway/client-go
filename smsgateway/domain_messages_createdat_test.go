package smsgateway_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestMessageState_JSON_Marshal_CreatedAt(t *testing.T) {
	tests := []struct {
		name       string
		createdAt  time.Time
		wantString string
	}{
		{
			name:       "UTC time serialized as RFC3339",
			createdAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			wantString: `"createdAt":"2020-01-01T00:00:00Z"`,
		},
		{
			name:       "zero value is still serialized (no omitempty)",
			createdAt:  time.Time{},
			wantString: `"createdAt":"0001-01-01T00:00:00Z"`,
		},
		{
			name:       "offset time preserved as RFC3339 with offset",
			createdAt:  time.Date(2025, 2, 14, 7, 0, 52, 245000000, time.FixedZone("", 3*60*60)),
			wantString: `"createdAt":"2025-02-14T07:00:52.245+03:00"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := smsgateway.MessageState{
				ID:        "msg-1",
				DeviceID:  "dev-1",
				State:     smsgateway.ProcessingStatePending,
				CreatedAt: tt.createdAt,
			}

			data, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if !strings.Contains(string(data), tt.wantString) {
				t.Errorf("json.Marshal() = %s, want substring %s", data, tt.wantString)
			}
		})
	}
}

func TestMessageState_JSON_Unmarshal_CreatedAt(t *testing.T) {
	tests := []struct {
		name     string
		jsonBody string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "RFC3339 UTC",
			jsonBody: `{"id":"msg-1","createdAt":"2020-01-01T00:00:00Z"}`,
			wantTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "RFC3339 with offset normalized to UTC",
			jsonBody: `{"id":"msg-1","createdAt":"2025-02-14T07:00:52.245+03:00"}`,
			wantTime: time.Date(2025, 2, 14, 4, 0, 52, 245000000, time.UTC),
		},
		{
			name:     "field absent does not break parsing",
			jsonBody: `{"id":"msg-1","state":"Pending"}`,
			wantTime: time.Time{},
		},
		{
			name:     "explicit null yields zero value",
			jsonBody: `{"id":"msg-1","createdAt":null}`,
			wantTime: time.Time{},
		},
		{
			name:     "invalid RFC3339 value returns error",
			jsonBody: `{"id":"msg-1","createdAt":"not-a-date"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m smsgateway.MessageState

			err := json.Unmarshal([]byte(tt.jsonBody), &m)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("json.Unmarshal() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if !m.CreatedAt.Equal(tt.wantTime) {
				t.Errorf("CreatedAt = %v, want %v", m.CreatedAt, tt.wantTime)
			}
		})
	}
}

func TestMessageState_JSON_RoundTrip_CreatedAt(t *testing.T) {
	original := smsgateway.MessageState{
		ID:       "msg-1",
		DeviceID: "dev-1",
		State:    smsgateway.ProcessingStateDelivered,
		CreatedAt: time.Date(
			2023, 6, 15, 12, 30, 45, 0, time.UTC,
		),
		Recipients: []smsgateway.RecipientState{
			{PhoneNumber: "79990001234", State: smsgateway.ProcessingStateDelivered},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded smsgateway.MessageState
	if jsonErr := json.Unmarshal(data, &decoded); jsonErr != nil {
		t.Fatalf("json.Unmarshal() error = %v", jsonErr)
	}

	if !decoded.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("round-trip CreatedAt = %v, want %v", decoded.CreatedAt, original.CreatedAt)
	}
}
