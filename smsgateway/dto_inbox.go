package smsgateway

import (
	"time"
)

// IncomingMessageType represents the type of inbox message.
type IncomingMessageType string

const (
	IncomingMessageTypeSMS           IncomingMessageType = "SMS"            // SMS message
	IncomingMessageTypeDataSMS       IncomingMessageType = "DATA_SMS"       // Data SMS message
	IncomingMessageTypeMMS           IncomingMessageType = "MMS"            // MMS message
	IncomingMessageTypeMmsDownloaded IncomingMessageType = "MMS_DOWNLOADED" // Downloaded MMS message
)

// IncomingMessage represents an inbox message received by the device.
//
// ID is the inbox message ID.
// Type is the type of the inbox message (SMS, DATA_SMS, MMS, MMS_DOWNLOADED).
// Sender is the inbox sender phone number.
// Recipient is the recipient phone number on the device.
// SimNumber is the SIM slot number.
// ContentPreview is the message body preview or metadata.
// IsEncrypted indicates whether the message is encrypted.
// Attachments is the list of MMS attachments.
// CreatedAt is the message received timestamp.
type IncomingMessage struct {
	ID             string              `json:"id"                    example:"PyDmBQZZXYmyxMwED8Fzy" validate:"required"`                    // Inbox message ID
	Type           IncomingMessageType `json:"type"                  example:"SMS"                   validate:"required"`                    // Message type
	Sender         string              `json:"sender"                example:"+79990001234"          validate:"required"`                    // Inbox sender phone number
	Recipient      *string             `json:"recipient,omitempty"   example:"+79990001234"          validate:"optional"`                    // Recipient phone number on the device
	SimNumber      *uint8              `json:"simNumber,omitempty"   example:"1"                     validate:"optional"`                    // SIM slot number
	ContentPreview string              `json:"contentPreview"        example:"Hello World!"          validate:"required"`                    // Message body preview or metadata
	IsEncrypted    bool                `json:"isEncrypted"           example:"true"`                                                         // Whether the message is encrypted
	Attachments    []InboxAttachment   `json:"attachments,omitempty"`                                                                        // MMS attachments
	CreatedAt      time.Time           `json:"createdAt"             example:"2020-01-01T00:00:00Z"  validate:"required" format:"date-time"` // Message received timestamp
}

// InboxAttachment represents an MMS attachment for an inbox message.
type InboxAttachment struct {
	PartID      int64  `json:"partId"      example:"1"`          // Attachment part ID
	Name        string `json:"name"        example:"photo.jpg"`  // Attachment file name
	Size        int64  `json:"size"        example:"102400"`     // Attachment file size
	ContentType string `json:"contentType" example:"image/jpeg"` // Attachment MIME type
}
