package smsgateway

// ErrorResponse represents a response to a request in case of an error.
//
// Message is an error message.
//
// Code is an error code, which is omitted if not specified.
//
// Data is an error context, which is omitted if not specified.
type ErrorResponse struct {
	Message string `json:"message"        example:"An error occurred"` // Error message
	Code    int32  `json:"code,omitempty"`                             // Error code
	Data    any    `json:"data,omitempty"`                             // Error context
}
