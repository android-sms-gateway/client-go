package smsgateway

import "time"

// Device self-information response
type MobileDeviceResponse struct {
	Device     *Device `json:"device,omitempty"`     // Device information, empty if device is not registered on the server
	ExternalIP string  `json:"externalIp,omitempty"` // External IP
}

// Device registration response
type MobileRegisterResponse struct {
	Id       string `json:"id" example:"QslD_GefqiYV6RQXdkM6V"`          // New device ID
	Token    string `json:"token" example:"bP0ZdK6rC6hCYZSjzmqhQ"`       // Device access token
	Login    string `json:"login" example:"VQ4GII"`                      // User login
	Password string `json:"password,omitempty" example:"cp2pydvxd2zwpx"` // User password, empty for existing user
}

// User one-time code response
type MobileUserCodeResponse struct {
	Code       string    `json:"code" example:"123456"`                     // One-time code
	ValidUntil time.Time `json:"validUntil" example:"2020-01-01T00:00:00Z"` // One-time code expiration time
}
