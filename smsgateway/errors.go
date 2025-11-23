package smsgateway

import "errors"

var (
	ErrConflictFields   = errors.New("conflict fields")
	ErrInvalidConfig    = errors.New("invalid config")
	ErrValidationFailed = errors.New("validation failed")
)
