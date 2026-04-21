package smsgateway

import "time"

// InboxRefreshRequest represents a request to export messages.
type InboxRefreshRequest struct {
	DeviceID        string                `json:"deviceId"                  validate:"required,max=21"                    example:"PyDmBQZZXYmyxMwED8Fzy"` // DeviceID is the ID of the device to export messages for.
	Since           time.Time             `json:"since"                     validate:"required,ltefield=Until"            example:"2024-01-01T00:00:00Z"`  // Since is the start of the time range to export.
	Until           time.Time             `json:"until"                     validate:"required,gtefield=Since"            example:"2024-01-01T23:59:59Z"`  // Until is the end of the time range to export.
	MessageTypes    []IncomingMessageType `json:"messageTypes,omitempty"    validate:"omitempty,min=1,dive,oneof=SMS MMS"`                                 // MessageTypes is the list of message types to export. By default, SMS messages are exported.
	TriggerWebhooks bool                  `json:"triggerWebhooks,omitempty" validate:"omitempty"                          example:"true"`                  // TriggerWebhooks indicates whether to trigger webhooks for the exported messages.
}

// MessagesExportRequest represents a request to export messages.
//
// Deprecated: use InboxRefreshRequest instead.
type MessagesExportRequest = InboxRefreshRequest
