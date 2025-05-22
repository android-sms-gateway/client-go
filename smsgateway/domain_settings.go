//nolint:lll // validator tags
package smsgateway

import "fmt"

// LimitPeriod defines the period for message sending limits.
type LimitPeriod string

const (
	// Disabled indicates no limit period.
	Disabled LimitPeriod = "Disabled"
	// PerMinute sets the limit period to per minute.
	PerMinute LimitPeriod = "PerMinute"
	// PerHour sets the limit period to per hour.
	PerHour LimitPeriod = "PerHour"
	// PerDay sets the limit period to per day.
	PerDay LimitPeriod = "PerDay"
)

// SimSelectionMode defines how SIM cards are selected for sending messages.
type SimSelectionMode string

const (
	// OSDefault uses the OS default SIM selection.
	OSDefault SimSelectionMode = "OSDefault"
	// RoundRobin cycles through SIM cards.
	RoundRobin SimSelectionMode = "RoundRobin"
	// Random selects SIM cards randomly.
	Random SimSelectionMode = "Random"
)

// DeviceSettings represents the overall configuration settings for a device.
type DeviceSettings struct {
	// Encryption contains settings related to message encryption.
	Encryption *SettingsEncryption `json:"encryption,omitempty"`

	// Messages contains settings related to message handling.
	Messages *SettingsMessages `json:"messages,omitempty"`

	// Ping contains settings related to ping functionality.
	Ping *SettingsPing `json:"ping,omitempty"`

	// Logs contains settings related to logging.
	Logs *SettingsLogs `json:"logs,omitempty"`

	// Webhooks contains settings related to webhook functionality.
	Webhooks *SettingsWebhooks `json:"webhooks,omitempty"`
}

func (s DeviceSettings) Validate() error {
	if s.Messages != nil {
		return s.Messages.Validate()
	}
	return nil
}

// SettingsEncryption contains settings related to message encryption.
type SettingsEncryption struct {
	// Passphrase is the encryption passphrase. If nil or empty, encryption is disabled.
	Passphrase *string `json:"passphrase,omitempty"`
}

// SettingsMessages contains settings related to message handling.
type SettingsMessages struct {
	// SendIntervalMin is the minimum interval between message sends (in seconds).
	// Must be at least 1 when provided.
	SendIntervalMin *int `json:"send_interval_min,omitempty" validate:"omitempty,min=1"`

	// SendIntervalMax is the maximum interval between message sends (in seconds).
	// Must be at least 1 when provided and greater than or equal to SendIntervalMin.
	SendIntervalMax *int `json:"send_interval_max,omitempty" validate:"omitempty,min=1"`

	// LimitPeriod defines the period for message sending limits.
	// Valid values are "Disabled", "PerMinute", "PerHour", or "PerDay".
	LimitPeriod *LimitPeriod `json:"limit_period,omitempty" validate:"omitempty,oneof=Disabled PerMinute PerHour PerDay"`

	// LimitValue is the maximum number of messages allowed per limit period.
	// Must be at least 1 when provided.
	LimitValue *int `json:"limit_value,omitempty" validate:"omitempty,min=1"`

	// SimSelectionMode defines how SIM cards are selected for sending messages.
	// Valid values are "OSDefault", "RoundRobin", or "Random".
	SimSelectionMode *SimSelectionMode `json:"sim_selection_mode,omitempty" validate:"omitempty,oneof=OSDefault RoundRobin Random"`

	// LogLifetimeDays is the number of days to retain message logs.
	// Must be at least 1 when provided.
	LogLifetimeDays *int `json:"log_lifetime_days,omitempty" validate:"omitempty,min=1"`
}

func (s SettingsMessages) Validate() error {
	if s.SendIntervalMax != nil && s.SendIntervalMin != nil && *s.SendIntervalMax < *s.SendIntervalMin {
		return fmt.Errorf("%w: sendIntervalMax must be greater than or equal to sendIntervalMin", ErrValidationFailed)
	}

	return nil
}

// SettingsPing contains settings related to ping functionality.
type SettingsPing struct {
	// IntervalSeconds is the interval between ping requests (in seconds).
	// Must be at least 1 when provided.
	IntervalSeconds *int `json:"interval_seconds,omitempty" validate:"omitempty,min=1"`
}

// SettingsLogs contains settings related to logging.
type SettingsLogs struct {
	// LifetimeDays is the number of days to retain logs.
	// Must be at least 1 when provided.
	LifetimeDays *int `json:"lifetime_days,omitempty" validate:"omitempty,min=1"`
}

// SettingsWebhooks contains settings related to webhook functionality.
type SettingsWebhooks struct {
	// InternetRequired indicates whether internet access is required for webhooks.
	InternetRequired *bool `json:"internet_required,omitempty"`

	// RetryCount is the number of times to retry failed webhook deliveries.
	// Must be at least 1 when provided.
	RetryCount *int `json:"retry_count,omitempty" validate:"omitempty,min=1"`

	// SigningKey is the secret key used for signing webhook payloads.
	SigningKey *string `json:"signing_key,omitempty"`
}
