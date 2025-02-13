package smsgateway

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

// Error response
type ErrorResponse struct {
	Message string `json:"message" example:"An error occurred"` // Error message
	Code    int32  `json:"code,omitempty"`                      // Error code
	Data    any    `json:"data,omitempty"`                      // Error context
}
