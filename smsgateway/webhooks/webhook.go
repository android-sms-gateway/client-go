package webhooks

import "github.com/android-sms-gateway/client-go/smsgateway"

// EventType is deprecated: use smsgateway package instead.
type EventType = smsgateway.WebhookEvent

const (
	// Deprecated: use smsgateway package instead.
	EventTypeSmsReceived EventType = smsgateway.WebhookEventSmsReceived
	// Deprecated: use smsgateway package instead.
	EventTypeSmsSent EventType = smsgateway.WebhookEventSmsSent
	// Deprecated: use smsgateway package instead.
	EventTypeSmsDelivered EventType = smsgateway.WebhookEventSmsDelivered
	// Deprecated: use smsgateway package instead.
	EventTypeSmsFailed EventType = smsgateway.WebhookEventSmsFailed
	// Deprecated: use smsgateway package instead.
	EventTypeSystemPing EventType = smsgateway.WebhookEventSystemPing
)

// Deprecated: use smsgateway package instead.
func IsValidEventType(e EventType) bool {
	return smsgateway.IsValidWebhookEvent(e)
}

// Webhook is deprecated: use smsgateway package instead.
type Webhook = smsgateway.Webhook
