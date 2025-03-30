package smsgateway_test

import (
	"errors"
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

// TestWebhookEventTypes tests that the event types returned by
// WebhookEventTypes are all valid.
func TestWebhookEventTypes(t *testing.T) {
	for _, v := range smsgateway.WebhookEventTypes() {
		if !smsgateway.IsValidWebhookEvent(v) {
			t.Errorf("event type %s is not valid", v)
		}
	}
}

// TestWebhook_Validate tests the Validate method of the Webhook struct.
func TestWebhook_Validate(t *testing.T) {
	tests := []struct {
		name    string
		webhook smsgateway.Webhook
		wantErr bool
		err     error
	}{
		{
			name: "Valid webhook with HTTPS URL",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: false,
		},
		{
			name: "Invalid event type",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: "invalid:event",
			},
			wantErr: true,
			err:     smsgateway.ErrValidationFailed,
		},
		{
			name: "Non-HTTPS URL",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "http://example.com/webhook",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: true,
			err:     smsgateway.ErrValidationFailed,
		},
		{
			name: "Empty URL",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: true,
			err:     smsgateway.ErrValidationFailed,
		},
		{
			name: "Valid webhook with sms:sent event",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: smsgateway.WebhookEventSmsSent,
			},
			wantErr: false,
		},
		{
			name: "Valid webhook with sms:delivered event",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: smsgateway.WebhookEventSmsDelivered,
			},
			wantErr: false,
		},
		{
			name: "Valid webhook with sms:failed event",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: smsgateway.WebhookEventSmsFailed,
			},
			wantErr: false,
		},
		{
			name: "Valid webhook with system:ping event",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https://example.com/webhook",
				Event: smsgateway.WebhookEventSystemPing,
			},
			wantErr: false,
		},
		{
			name: "URL with uppercase HTTPS",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "HTTPS://example.com/webhook",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: false,
		},
		{
			name: "FTP URL",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "ftp://example.com/webhook",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: true,
			err:     smsgateway.ErrValidationFailed,
		},
		{
			name: "Malformed URL",
			webhook: smsgateway.Webhook{
				ID:    "test-id",
				URL:   "https:/example.com",
				Event: smsgateway.WebhookEventSmsReceived,
			},
			wantErr: true,
			err:     smsgateway.ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.webhook.Validate()

			// Check if an error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, check the error message
			if tt.wantErr && !errors.Is(err, tt.err) {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.err)
			}
		})
	}
}
