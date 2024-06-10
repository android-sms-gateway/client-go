package smsgateway

// The type of event.
type PushEventType string

const (
	// A message is enqueued.
	PushMessageEnqueued PushEventType = "MessageEnqueued"
	// Webhooks are updated.
	PushWebhooksUpdated PushEventType = "WebhooksUpdated"
)

// A push notification.
type PushNotification struct {
	// The token of the device that receives the notification.
	Token string `json:"token" validate:"required" example:"PyDmBQZZXYmyxMwED8Fzy"`
	// The type of event.
	Event PushEventType `json:"event" validate:"omitempty,oneof=MessageEnqueued WebhooksUpdated" default:"MessageEnqueued" example:"MessageEnqueued"`
	// The additional data associated with the event.
	Data map[string]string `json:"data"`
}
