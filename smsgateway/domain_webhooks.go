package smsgateway

import (
	"fmt"
	"strings"
)

type WebhookEvent = string

const (
	WebhookEventMmsReceived     WebhookEvent = "mms:received"      // Triggered when an MMS is received.
	WebhookEventSmsDataReceived WebhookEvent = "sms:data-received" // Triggered when a data SMS is received.
	WebhookEventSmsDelivered    WebhookEvent = "sms:delivered"     // Triggered when an SMS is delivered.
	WebhookEventSmsFailed       WebhookEvent = "sms:failed"        // Triggered when an SMS processing fails.
	WebhookEventSmsReceived     WebhookEvent = "sms:received"      // Triggered when an SMS is received.
	WebhookEventSmsSent         WebhookEvent = "sms:sent"          // Triggered when an SMS is sent.
	WebhookEventSystemPing      WebhookEvent = "system:ping"       // Triggered when the device pings the server.
)

//nolint:gochecknoglobals // lookup table
var allEventTypes = map[WebhookEvent]struct{}{
	WebhookEventSmsReceived:     {},
	WebhookEventSmsDataReceived: {},
	WebhookEventSmsSent:         {},
	WebhookEventSmsDelivered:    {},
	WebhookEventSmsFailed:       {},
	WebhookEventSystemPing:      {},
	WebhookEventMmsReceived:     {},
}

// WebhookEventTypes returns a slice of all supported webhook event types.
func WebhookEventTypes() []WebhookEvent {
	return []WebhookEvent{
		WebhookEventSmsReceived,
		WebhookEventSmsDataReceived,
		WebhookEventSmsSent,
		WebhookEventSmsDelivered,
		WebhookEventSmsFailed,
		WebhookEventSystemPing,
		WebhookEventMmsReceived,
	}
}

// IsValidWebhookEvent checks if the webhook event is a valid type.
// It takes a webhook event type and returns true if the event is valid, false otherwise.
func IsValidWebhookEvent(e WebhookEvent) bool {
	_, ok := allEventTypes[e]
	return ok
}

// Webhook represents a webhook configuration.
//
// ID is the unique identifier of the webhook.
//
// DeviceID is the unique identifier of the device the webhook is associated with.
//
// URL is the URL the webhook will be sent to.
//
// Event is the type of event the webhook is triggered for.
type Webhook struct {
	ID       string       `json:"id,omitempty" validate:"max=36"            example:"123e4567-e89b-12d3-a456-426614174000"` // The unique identifier of the webhook.
	DeviceID *string      `json:"deviceId"     validate:"omitempty,max=21"  example:"PyDmBQZZXYmyxMwED8Fzy"`                // The unique identifier of the device the webhook is associated with.
	URL      string       `json:"url"          validate:"required,http_url" example:"https://example.com/webhook"`          // The URL the webhook will be sent to.
	Event    WebhookEvent `json:"event"        validate:"required"          example:"sms:received"`                         // The type of event the webhook is triggered for.
}

// Validate checks if the webhook is configured correctly.
// Returns nil if validation passes, or an appropriate error otherwise.
func (w Webhook) Validate() error {
	if !IsValidWebhookEvent(w.Event) {
		return fmt.Errorf("%w: invalid event type", ErrValidationFailed)
	}

	if !strings.HasPrefix(strings.ToLower(w.URL), "https://") {
		return fmt.Errorf("%w: url must start with https://", ErrValidationFailed)
	}

	return nil
}
