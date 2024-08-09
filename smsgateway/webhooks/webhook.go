package webhooks

type EventType string

const (
	// Triggered when an SMS is received.
	EventTypeSmsReceived EventType = "sms:received"
	// Triggered when an SMS is sent.
	EventTypeSmsSent EventType = "sms:sent"
	// Triggered when the device pings the server.
	EventTypeSystemPing EventType = "system:ping"
)

var allEventTypes = map[EventType]struct{}{
	EventTypeSmsReceived: {},
	EventTypeSmsSent:     {},
	EventTypeSystemPing:  {},
}

// IsValid checks if the given event type is valid.
//
// e is the event type to be checked.
// Returns true if the event type is valid, false otherwise.
func IsValid(e EventType) bool {
	_, ok := allEventTypes[e]
	return ok
}

// A webhook configuration.
type Webhook struct {
	// The unique identifier of the webhook.
	ID string `json:"id" validate:"max=36" example:"123e4567-e89b-12d3-a456-426614174000"`

	// The URL the webhook will be sent to.
	URL string `json:"url" validate:"required,http_url" example:"https://example.com/webhook"`

	// The type of event the webhook is triggered for.
	Event EventType `json:"event" validate:"required" example:"sms:received"`
}
