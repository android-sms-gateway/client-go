package smsgateway

import "time"

// MobileRegisterRequest represents a request to register a mobile device.
//
// The Name field contains the name of the device, and the PushToken field
// contains the FCM token of the device.
type MobileRegisterRequest struct {
	// Name of the device (optional)
	// +optional
	Name *string `json:"name,omitempty" validate:"omitempty,max=128" example:"Android Phone"`
	// FCM token of the device (optional)
	// +optional
	PushToken *string `json:"pushToken"      validate:"omitempty,max=256" example:"gHz-T6NezDlOfllr7F-Be"`
}

// MobileUpdateRequest represents a request to update a mobile device.
//
// The Id field contains the device ID.
//
// The PushToken field contains the FCM token of the device.
type MobileUpdateRequest struct {
	//nolint:revive,staticcheck // backward compatibility
	Id        string `json:"id"        example:"QslD_GefqiYV6RQXdkM6V"`                              // Device ID
	PushToken string `json:"pushToken" example:"gHz-T6NezDlOfllr7F-Be" validate:"omitempty,max=256"` // FCM token of the device (optional)
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
	ID string `json:"id"         validate:"required,max=36"     example:"PyDmBQZZXYmyxMwED8Fzy"`
	// State
	State ProcessingState `json:"state"      validate:"required"            example:"Pending"`
	// Recipients states
	Recipients []RecipientState `json:"recipients" validate:"required,min=1,dive"`
	// History of states
	States map[string]time.Time `json:"states"`
}

// MobilePatchMessageRequest represents a request to patch messages.
type MobilePatchMessageRequest []MobilePatchMessageItem
