package smsgateway

import "time"

// Push request
type UpstreamPushRequest = []PushNotification

// Messages export request
type MessagesExportRequest struct {
	// DeviceID is the ID of the device to export messages for.
	DeviceID string `json:"deviceId" example:"PyDmBQZZXYmyxMwED8Fzy" validate:"required,max=21"`
	// Since is the start of the time range to export.
	Since time.Time `json:"since" example:"2024-01-01T00:00:00Z" validate:"required,ltefield=Until"`
	// Until is the end of the time range to export.
	Until time.Time `json:"until" example:"2024-01-01T23:59:59Z" validate:"required,gtefield=Since"`
}
