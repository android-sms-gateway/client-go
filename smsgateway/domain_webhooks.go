package smsgateway

import (
	"fmt"
	"strings"
)

type WebhookEvent = string

const (
	// Triggered when an SMS is received.
	WebhookEventSmsReceived WebhookEvent = "sms:received"
	// Triggered when a data SMS is received.
	WebhookEventSmsDataReceived WebhookEvent = "sms:data-received"
	// Triggered when an SMS is sent.
	WebhookEventSmsSent WebhookEvent = "sms:sent"
	// Triggered when an SMS is delivered.
	WebhookEventSmsDelivered WebhookEvent = "sms:delivered"
	// Triggered when an SMS processing fails.
	WebhookEventSmsFailed WebhookEvent = "sms:failed"
	// Triggered when the device pings the server.
	WebhookEventSystemPing WebhookEvent = "system:ping"
)

//nolint:gochecknoglobals // lookup table
var allEventTypes = map[WebhookEvent]struct{}{
	WebhookEventSmsReceived:     {},
	WebhookEventSmsDataReceived: {},
	WebhookEventSmsSent:         {},
	WebhookEventSmsDelivered:    {},
	WebhookEventSmsFailed:       {},
	WebhookEventSystemPing:      {},
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
	}
}

// IsValid checks if the given event type is valid.
//
// e is the event type to be checked.
// Returns true if the event type is valid, false otherwise.
func IsValidWebhookEvent(e WebhookEvent) bool {
	_, ok := allEventTypes[e]
	return ok
}

// A webhook configuration.
type Webhook struct {
	// The unique identifier of the webhook.
	ID string `json:"id,omitempty" validate:"max=36" example:"123e4567-e89b-12d3-a456-426614174000"`

	// The unique identifier of the device the webhook is associated with.
	DeviceID *string `json:"deviceId" validate:"omitempty,max=21" example:"PyDmBQZZXYmyxMwED8Fzy"`

	// The URL the webhook will be sent to.
	URL string `json:"url" validate:"required,http_url" example:"https://example.com/webhook"`

	// The type of event the webhook is triggered for.
	Event WebhookEvent `json:"event" validate:"required" example:"sms:received"`
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
