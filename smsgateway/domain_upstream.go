//nolint:lll // validator tags
package smsgateway

// PushEventType is the type of a push notification.
type PushEventType string

const (
	PushMessageEnqueued         PushEventType = "MessageEnqueued"         // Message is enqueued.
	PushMessageCancelled        PushEventType = "MessageCancelled"        // Message cancellation is requested.
	PushWebhooksUpdated         PushEventType = "WebhooksUpdated"         // Webhooks are updated.
	PushMessagesExportRequested PushEventType = "MessagesExportRequested" // Messages export is requested.
	PushSettingsUpdated         PushEventType = "SettingsUpdated"         // Settings are updated.
)

// PushNotification represents a push notification.
//
// The token of the device that receives the notification.
type PushNotification struct {
	Token     string            `json:"token"                   validate:"required"                                                                                       example:"PyDmBQZZXYmyxMwED8Fzy"` // The token of the device that receives the notification.
	Event     PushEventType     `json:"event"                   validate:"oneof=MessageEnqueued MessageCancelled WebhooksUpdated MessagesExportRequested SettingsUpdated" example:"MessageEnqueued"`       // The type of event.
	Data      map[string]string `json:"data"`                                                                                                                                                              // The additional data associated with the event.
	TokenType *PushTokenType    `json:"pushTokenType,omitempty" validate:"omitempty,oneof=token fid"                                                                      example:"fid"`                   // Which FCM target field to use for the token.
}
