package smsgateway

import "github.com/android-sms-gateway/client-go/smsgateway/webhooks"

// Deprecated: use webhooks package instead.
type WebhookEvent = webhooks.EventType

const (
	// Deprecated: use webhooks package instead.
	WebhookEventSmsReceived WebhookEvent = webhooks.EventTypeSmsReceived
	// Deprecated: use webhooks package instead.
	WebhookEventSystemPing WebhookEvent = webhooks.EventTypeSystemPing
)

// Deprecated: use webhook package instead.
type Webhook = webhooks.Webhook
