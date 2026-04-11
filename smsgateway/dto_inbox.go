package smsgateway

import (
	"time"
)

// IncomingMessageType represents the type of incoming message.
type IncomingMessageType string

const (
	IncomingMessageTypeSMS           IncomingMessageType = "SMS"            // SMS message
	IncomingMessageTypeDataSMS       IncomingMessageType = "DATA_SMS"       // Data SMS message
	IncomingMessageTypeMMS           IncomingMessageType = "MMS"            // MMS message
	IncomingMessageTypeMmsDownloaded IncomingMessageType = "MMS_DOWNLOADED" // Downloaded MMS message
)

// IncomingMessage represents an incoming message received by the device.
//
// ID is the incoming message ID.
// Type is the type of the incoming message (SMS, DATA_SMS, MMS, MMS_DOWNLOADED).
// Sender is the incoming sender phone number.
// Recipient is the recipient phone number on the device.
// SimNumber is the SIM slot number.
// ContentPreview is the message body preview or metadata.
// CreatedAt is the message received timestamp.
type IncomingMessage struct {
	ID             string              `json:"id"                  example:"PyDmBQZZXYmyxMwED8Fzy" validate:"required"`                    // Incoming message ID
	Type           IncomingMessageType `json:"type"                example:"SMS"                   validate:"required"`                    // Message type
	Sender         string              `json:"sender"              example:"+79990001234"          validate:"required"`                    // Incoming sender phone number
	Recipient      *string             `json:"recipient,omitempty" example:"+79990001234"          validate:"optional"`                    // Recipient phone number on the device
	SimNumber      *uint8              `json:"simNumber,omitempty" example:"1"                     validate:"optional"`                    // SIM slot number
	ContentPreview string              `json:"contentPreview"      example:"Hello World!"          validate:"required"`                    // Message body preview or metadata
	CreatedAt      time.Time           `json:"createdAt"           example:"2020-01-01T00:00:00Z"  validate:"required" format:"date-time"` // Message received timestamp
}
