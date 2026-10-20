package webhooks_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway/webhooks"
)

// TestIsValidEventType tests the IsValidEventType function.
func TestIsValidEventType(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := webhooks.IsValidEventType(tt.e); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
