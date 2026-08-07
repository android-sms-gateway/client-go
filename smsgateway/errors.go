package smsgateway

import "errors"

var (
	ErrConflictFields   = errors.New("conflict fields")
	ErrInvalidConfig    = errors.New("invalid config")
	ErrValidationFailed = errors.New("validation failed")

	// ErrDeviceIDRequired is returned when an E2E-encrypted message is sent
	// without a DeviceID. The device listing key lookup requires a target
	// device; with an empty DeviceID the server would route the message to a
	// random device and the E2E payload would be undecryptable.
	ErrDeviceIDRequired = errors.New("device id is required for E2E encrypted messages")

	// ErrE2ENotConfigured is returned when the target device has no publicKey
	// or no keyVersion. There is intentionally NO plaintext fallback: the
	// message would be undecryptable by the device.
	ErrE2ENotConfigured = errors.New("E2E encryption is not configured for the target device")

	// ErrDeviceNotFound is returned when a message targets a device that does
	// not exist in the device listing.
	ErrDeviceNotFound = errors.New("device not found")
)
