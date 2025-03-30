package smsgateway

// Device registration request
type MobileRegisterRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,max=128" example:"Android Phone"`    // Device name
	PushToken *string `json:"pushToken" validate:"omitempty,max=256" example:"gHz-T6NezDlOfllr7F-Be"` // FCM token
}

// Device update request
type MobileUpdateRequest struct {
	//nolint:revive // backward compatibility
	// ID
	Id        string `json:"id" example:"QslD_GefqiYV6RQXdkM6V"`
	PushToken string `json:"pushToken" validate:"omitempty,max=256" example:"gHz-T6NezDlOfllr7F-Be"` // FCM token
}

// Device change password request
type MobileChangePasswordRequest struct {
	// Current password
	CurrentPassword string `json:"currentPassword" validate:"required" example:"cp2pydvxd2zwpx"`
	// New password, at least 14 characters
	NewPassword string `json:"newPassword" validate:"required,min=14" example:"cp2pydvxd2zwpx"`
}
