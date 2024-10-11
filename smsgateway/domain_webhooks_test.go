package smsgateway_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

// TestIsValidEventType tests the IsValidEventType function.
func TestIsValidEventType(t *testing.T) {
	tests := []struct {
		name string
		e    smsgateway.WebhookEvent
		want bool
	}{
		{
			name: "Valid event type",
			e:    smsgateway.WebhookEventSmsDelivered,
			want: true,
		},
		{
			name: "Invalid event type",
			e:    "invalid:event",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := smsgateway.IsValidWebhookEvent(tt.e); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
