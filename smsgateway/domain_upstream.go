//nolint:lll // validator tags
package smsgateway

// PushEventType is the type of a push notification.
type PushEventType string

const (
	PushMessageEnqueued         PushEventType = "MessageEnqueued"         // Message is enqueued.
	PushWebhooksUpdated         PushEventType = "WebhooksUpdated"         // Webhooks are updated.
	PushMessagesExportRequested PushEventType = "MessagesExportRequested" // Messages export is requested.
	PushSettingsUpdated         PushEventType = "SettingsUpdated"         // Settings are updated.
)

// PushNotification represents a push notification.
//
// The token of the device that receives the notification.
type PushNotification struct {
	Token string            `json:"token" validate:"required"                                                                      example:"PyDmBQZZXYmyxMwED8Fzy"` // The token of the device that receives the notification.
	Event PushEventType     `json:"event" validate:"oneof=MessageEnqueued WebhooksUpdated MessagesExportRequested SettingsUpdated" example:"MessageEnqueued"`       // The type of event.
	Data  map[string]string `json:"data"`                                                                                                                           // The additional data associated with the event.
}
