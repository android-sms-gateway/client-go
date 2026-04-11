package smsgateway

import "time"

// SmsEventPayload represents the base payload for message-related events.
//
// MessageID is the unique identifier of the message.
// PhoneNumber is the phone number of the sender (for incoming messages) or recipient (for outgoing messages).
// Sender is the phone number of the message sender. For incoming messages, this is the external sender; for outgoing messages, this is the device's phone number. May be empty for device's phone number.
// Recipient is the phone number of the message recipient. For incoming messages, this is the device's/SIM's phone number that received the message; for outgoing messages, this is the recipient's phone number. May be nil for device's phone number.
// SimNumber is the SIM card number that sent or received the SMS. May be nil if the SIM cannot be determined or the default was used.
type SmsEventPayload struct {
	MessageID   string  `json:"messageId"           example:"PyDmBQZZXYmyxMwED8Fzy"` // The unique identifier of the message.
	PhoneNumber string  `json:"phoneNumber"         example:"+79990001234"`          // The phone number of the sender (for incoming messages) or recipient (for outgoing messages).
	Sender      string  `json:"sender"              example:"+79990001234"`          // The phone number of the message sender.
	Recipient   *string `json:"recipient,omitempty" example:"+79990001234"`          // The phone number of the message recipient.
	SimNumber   *uint8  `json:"simNumber,omitempty" example:"1"`                     // The SIM card number that sent or received the SMS.
}

// SmsReceivedPayload represents the payload of `sms:received` event.
//
// SmsEventPayload is the base payload for message-related events.
// Message is the content of the SMS message received.
// ReceivedAt is the timestamp when the SMS message was received.
type SmsReceivedPayload struct {
	SmsEventPayload

	Message    string    `json:"message"    example:"Hello World!"`         // The content of the SMS message received.
	ReceivedAt time.Time `json:"receivedAt" example:"2020-01-01T00:00:00Z"` // The timestamp when the SMS message was received.
}

// SmsSentPayload represents the payload of `sms:sent` event.
//
// SmsEventPayload is the base payload for message-related events.
// SentAt is the timestamp when the SMS message was sent.
type SmsSentPayload struct {
	SmsEventPayload

	SentAt time.Time `json:"sentAt" example:"2020-01-01T00:00:00Z"` // The timestamp when the SMS message was sent.
}

// SmsDeliveredPayload represents the payload of `sms:delivered` event.
//
// SmsEventPayload is the base payload for message-related events.
// DeliveredAt is the timestamp when the SMS message was delivered.
type SmsDeliveredPayload struct {
	SmsEventPayload

	DeliveredAt time.Time `json:"deliveredAt" example:"2020-01-01T00:00:00Z"` // The timestamp when the SMS message was delivered.
}

// SmsFailedPayload represents the payload of `sms:failed` event.
//
// SmsEventPayload is the base payload for message-related events.
// FailedAt is the timestamp when the SMS message failed.
// Reason is the reason for the failure.
type SmsFailedPayload struct {
	SmsEventPayload

	FailedAt time.Time `json:"failedAt" example:"2020-01-01T00:00:00Z"` // The timestamp when the SMS message failed.
	Reason   string    `json:"reason"   example:"timeout"`              // The reason for the failure.
}

// SmsDataReceivedPayload represents the payload of `sms:data-received` event.
//
// SmsEventPayload is the base payload for message-related events.
// Data is the Base64-encoded content of the SMS message received.
// ReceivedAt is the timestamp when the SMS message was received.
type SmsDataReceivedPayload struct {
	SmsEventPayload

	Data       string    `json:"data"       example:"SGVsbG8gV29ybGQh"     format:"byte"` // Base64-encoded content of the SMS message received.
	ReceivedAt time.Time `json:"receivedAt" example:"2020-01-01T00:00:00Z"`               // The timestamp when the SMS message was received.
}

// MmsReceivedPayload represents the payload of `mms:received` event.
//
// SmsEventPayload is the base payload for message-related events.
// TransactionID is the unique MMS transaction identifier.
// Subject is the message subject line.
// ContentClass is the MMS content classification.
// Size is the attachment size in bytes.
// ReceivedAt is the timestamp when the MMS message was received.
type MmsReceivedPayload struct {
	SmsEventPayload

	TransactionID string    `json:"transactionId"     example:"abc123"`               // Unique MMS transaction identifier
	Subject       *string   `json:"subject,omitempty" example:"Hello"`                // Message subject line
	ContentClass  string    `json:"contentClass"      example:"MMS"`                  // MMS content classification
	Size          int       `json:"size"              example:"1024"`                 // Attachment size in bytes
	ReceivedAt    time.Time `json:"receivedAt"        example:"2020-01-01T00:00:00Z"` // The timestamp when the MMS message was received.
}

// MmsDownloadedAttachment represents metadata for a non-text MMS part.
//
// PartID is the _id from content://mms/part.
// ContentType is the MIME type of the attachment (e.g. image/jpeg, audio/amr).
// Name is the filename of the attachment, if present.
// Data is the Base64-encoded attachment data, nil if unavailable.
// Size is the size in bytes, 0 if unknown.
type MmsDownloadedAttachment struct {
	PartID      int     `json:"partId"         example:"1"`                // The _id from content://mms/part.
	ContentType string  `json:"contentType"    example:"image/jpeg"`       // MIME type of the attachment
	Name        *string `json:"name,omitempty" example:"photo.jpg"`        // Filename of the attachment, if present.
	Data        *string `json:"data,omitempty" example:"SGVsbG8gV29ybGQh"` // Base64-encoded attachment data, nil if unavailable.
	Size        *int    `json:"size,omitempty" example:"1024"`             // Size in bytes, nil if unknown.
}

// MmsDownloadedPayload represents the payload of `mms:downloaded` event.
//
// SmsEventPayload is the base payload for message-related events.
// Subject is the message subject line.
// Body is the aggregated text content of the MMS message.
// Attachments is the metadata for non-text MMS parts, including optional Base64 content.
// ReceivedAt is the timestamp when the MMS message was received.
type MmsDownloadedPayload struct {
	SmsEventPayload

	Subject     *string                   `json:"subject,omitempty" example:"Hello"`                // Message subject line.
	Body        *string                   `json:"body,omitempty"    example:"Hello World!"`         // Aggregated text content of the MMS message.
	Attachments []MmsDownloadedAttachment `json:"attachments"       example:"[{}]"`                 // Metadata for non-text MMS parts, including optional Base64 content.
	ReceivedAt  time.Time                 `json:"receivedAt"        example:"2020-01-01T00:00:00Z"` // The timestamp when the MMS message was received.
}

// SystemPingPayload represents the payload of `system:ping` event.
type SystemPingPayload struct{}
