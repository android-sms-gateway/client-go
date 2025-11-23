package smsgateway

type HealthStatus string

const (
	HealthStatusPass HealthStatus = "pass"
	HealthStatusWarn HealthStatus = "warn"
	HealthStatusFail HealthStatus = "fail"
)

// HealthCheck represents the result of a health check.
//
// Description is a human-readable description of the check.
//
// ObservedUnit is the unit of measurement for the observed value.
//
// ObservedValue is the observed value of the check.
//
// Status is the status of the check.
// It can be one of the following values: "pass", "warn", or "fail".
type HealthCheck struct {
	// A human-readable description of the check.
	Description string `json:"description,omitempty"`
	// Unit of measurement for the observed value.
	ObservedUnit string `json:"observedUnit,omitempty"`
	// Observed value of the check.
	ObservedValue int `json:"observedValue"`
	// Status of the check.
	// It can be one of the following values: "pass", "warn", or "fail".
	Status HealthStatus `json:"status"`
}

// HealthChecks is a map of check names to their respective details.
type HealthChecks map[string]HealthCheck

// HealthResponse represents the result of a health check.
//
// Status is the overall status of the application.
// It can be one of the following values: "pass", "warn", or "fail".
//
// Version is the version of the application.
//
// ReleaseID is the release ID of the application.
// It is used to identify the version of the application.
//
// Checks is a map of check names to their respective details.
type HealthResponse struct {
	// Overall status of the application.
	// It can be one of the following values: "pass", "warn", or "fail".
	Status HealthStatus `json:"status"`
	// Version of the application.
	Version string `json:"version,omitempty"`
	// Release ID of the application.
	// It is used to identify the version of the application.
	ReleaseID int `json:"releaseId,omitempty"`
	// A map of check names to their respective details.
	Checks HealthChecks `json:"checks,omitempty"`
}
