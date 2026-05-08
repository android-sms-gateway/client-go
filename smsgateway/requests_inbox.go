package smsgateway

import "time"

// InboxRefreshRequest represents a request to export messages.
type InboxRefreshRequest struct {
	DeviceID        string                `json:"deviceId"                  validate:"required,max=21"                    example:"PyDmBQZZXYmyxMwED8Fzy"`                    // ID of the device to refresh messages for".
	Since           time.Time             `json:"since"                     validate:"required,ltefield=Until"            example:"2024-01-01T00:00:00Z"  format:"date-time"` // Start of the time range to refresh.
	Until           time.Time             `json:"until"                     validate:"required,gtefield=Since"            example:"2024-01-01T23:59:59Z"  format:"date-time"` // End of the time range to refresh.
	MessageTypes    []IncomingMessageType `json:"messageTypes,omitempty"    validate:"omitempty,min=1,dive,oneof=SMS MMS"`                                                    // List of message types to refresh. By default, SMS messages are refreshed.
	TriggerWebhooks bool                  `json:"triggerWebhooks,omitempty" validate:"omitempty"                          example:"true"`                                     // Indicates whether to trigger webhooks for the refreshed messages.
}

// MessagesExportRequest represents a request to export messages.
//
// Deprecated: use InboxRefreshRequest instead.
type MessagesExportRequest = InboxRefreshRequest
