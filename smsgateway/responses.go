package smsgateway

// Error response
type ErrorResponse struct {
	Message string `json:"message" example:"An error occurred"` // Error message
	Code    int32  `json:"code,omitempty"`                      // Error code
	Data    any    `json:"data,omitempty"`                      // Error context
}
