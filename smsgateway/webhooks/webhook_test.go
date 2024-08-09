package webhooks_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway/webhooks"
)

// TestIsValid tests the IsValid function.
func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		e    webhooks.EventType
		want bool
	}{
		{
			name: "Valid event type",
			e:    webhooks.EventTypeSmsReceived,
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
			if got := webhooks.IsValid(tt.e); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
