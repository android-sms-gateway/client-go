package smsgateway

import "time"

// MobileRegisterRequest represents a request to register a mobile device.
type MobileRegisterRequest struct {
	Name      *string   `json:"name,omitempty"     validate:"omitempty,max=128" example:"Android Phone"`         // Name of the device (optional)
	PushToken *string   `json:"pushToken"          validate:"omitempty,max=256" example:"gHz-T6NezDlOfllr7F-Be"` // FCM token of the device (optional)
	SimCards  []SimCard `json:"simCards,omitempty"`                                                              // SIM cards (optional)
}

// MobileUpdateRequest represents a request to update a mobile device.
type MobileUpdateRequest struct {
	//nolint:revive,staticcheck // backward compatibility
	Id        string    `json:"id"                 example:"QslD_GefqiYV6RQXdkM6V"`                              // Device ID
	PushToken string    `json:"pushToken"          example:"gHz-T6NezDlOfllr7F-Be" validate:"omitempty,max=256"` // FCM token of the device (optional)
	SimCards  []SimCard `json:"simCards,omitempty"`                                                              // SIM cards (optional)
}

// MobileChangePasswordRequest represents a request to change the password of a mobile device.
//
// The CurrentPassword field contains the current password of the device.
//
// The NewPassword field contains the new password of the device. It must be at least 14 characters long.
type MobileChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"        example:"cp2pydvxd2zwpx"` // Current password
	NewPassword     string `json:"newPassword"     validate:"required,min=14" example:"cp2pydvxd2zwpx"` // New password, at least 14 characters
}

// MobilePatchMessageItem represents a single message patch request.
type MobilePatchMessageItem struct {
	// Message ID
	ID string `json:"id" validate:"required,max=36" example:"PyDmBQZZXYmyxMwED8Fzy"`
	// State
	State ProcessingState `json:"state" validate:"required" example:"Pending"`
	// Recipients states
	Recipients []RecipientState `json:"recipients" validate:"required,min=1,dive"`
	// History of states
	States map[string]time.Time `json:"states"`
}

// MobilePatchMessageRequest represents a request to patch messages.
type MobilePatchMessageRequest []MobilePatchMessageItem

// MobilePostInboxRequestItemAttachment represents an attachment in an inbox message upload.
type MobilePostInboxRequestItemAttachment struct {
	PartID      int64  `json:"partId"         validate:"required"         example:"1"`
	ContentType string `json:"contentType"    validate:"required,max=128" example:"image/jpeg"`
	Name        string `json:"name"           validate:"required,max=512" example:"part001"`
	Size        *int64 `json:"size,omitempty"                             example:"12345"`
	Data        []byte `json:"data"           validate:"required"`
}

// MobilePostInboxRequestItem represents a single inbox message upload.
type MobilePostInboxRequestItem struct {
	ID          string                                 `json:"id"                    validate:"required,max=36"   example:"PyDmBQZZXYmyxMwED8Fzy"`
	Type        IncomingMessageType                    `json:"type"                  validate:"required"          example:"SMS"`
	Sender      string                                 `json:"sender"                validate:"required,max=512"  example:"+79990001234"`
	Recipient   *string                                `json:"recipient,omitempty"   validate:"omitempty,max=512" example:"+79990001234"`
	SimNumber   *uint8                                 `json:"simNumber,omitempty"                                example:"1"`
	Content     string                                 `json:"content"               validate:"required"`
	IsEncrypted bool                                   `json:"isEncrypted"           validate:"required"          example:"true"`
	CreatedAt   time.Time                              `json:"createdAt"             validate:"required"                                          format:"date-time"`
	Attachments []MobilePostInboxRequestItemAttachment `json:"attachments,omitempty" validate:"omitempty,dive"`
}

// MobilePostInboxRequest represents a batch of inbox messages from a device.
type MobilePostInboxRequest []MobilePostInboxRequestItem
