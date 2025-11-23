package smsgateway

import "time"

// MobileDeviceResponse contains device information and external IP address.
//
// Device is empty if the device is not registered on the server.
type MobileDeviceResponse struct {
	// Device information, empty if device is not registered on the server
	Device *Device `json:"device,omitempty"`
	// External IP address
	ExternalIP string `json:"externalIp,omitempty"`
}

// MobileRegisterResponse contains device registration response.
//
// Id is the new device ID.
// Token is the device access token.
// Login is the user login.
// Password is the user password, empty for existing user.
type MobileRegisterResponse struct {
	//nolint:revive,staticcheck // backward compatibility
	Id       string `json:"id"                 example:"QslD_GefqiYV6RQXdkM6V"` // New device ID
	Token    string `json:"token"              example:"bP0ZdK6rC6hCYZSjzmqhQ"` // Device access token
	Login    string `json:"login"              example:"VQ4GII"`                // User login
	Password string `json:"password,omitempty" example:"cp2pydvxd2zwpx"`        // User password, empty for existing user
}

// MobileUserCodeResponse represents a one-time code response for mobile clients.
//
// Code is the one-time code sent to the user.
// ValidUntil is the one-time code expiration time.
type MobileUserCodeResponse struct {
	Code       string    `json:"code"       example:"123456"`               // One-time code sent to the user
	ValidUntil time.Time `json:"validUntil" example:"2020-01-01T00:00:00Z"` // One-time code expiration time
}

// MobileMessage represents a message for mobile clients.
//
// It contains the message information and message creation time.
type MobileMessage struct {
	Message // Message information

	CreatedAt time.Time `json:"createdAt" example:"2020-01-01T00:00:00Z"` // Message creation time
}

// MobileGetMessagesResponse represents a collection of messages for mobile clients.
type MobileGetMessagesResponse []MobileMessage
