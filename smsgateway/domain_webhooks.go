package smsgateway

type WebhookEvent = string

const (
	// Triggered when an SMS is received.
	WebhookEventSmsReceived WebhookEvent = "sms:received"
	// Triggered when an SMS is sent.
	WebhookEventSmsSent WebhookEvent = "sms:sent"
	// Triggered when an SMS is delivered.
	WebhookEventSmsDelivered WebhookEvent = "sms:delivered"
	// Triggered when an SMS processing fails.
	WebhookEventSmsFailed WebhookEvent = "sms:failed"
	// Triggered when the device pings the server.
	WebhookEventSystemPing WebhookEvent = "system:ping"
)

var allEventTypes = map[WebhookEvent]struct{}{
	WebhookEventSmsReceived:  {},
	WebhookEventSmsSent:      {},
	WebhookEventSmsDelivered: {},
	WebhookEventSmsFailed:    {},
	WebhookEventSystemPing:   {},
}

// WebhookEventTypes returns a slice of all supported webhook event types.
func WebhookEventTypes() []WebhookEvent {
	return []WebhookEvent{
		WebhookEventSmsReceived,
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

	// The URL the webhook will be sent to.
	URL string `json:"url" validate:"required,http_url" example:"https://example.com/webhook"`

	// The type of event the webhook is triggered for.
	Event WebhookEvent `json:"event" validate:"required" example:"sms:received"`
}
