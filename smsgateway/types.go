package smsgateway

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrConflictFields   = errors.New("conflict fields")
)
