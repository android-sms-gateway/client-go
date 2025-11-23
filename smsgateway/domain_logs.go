package smsgateway

import "time"

type LogEntryPriority string

const (
	LogEntryPriorityDebug LogEntryPriority = "DEBUG"
	LogEntryPriorityInfo  LogEntryPriority = "INFO"
	LogEntryPriorityWarn  LogEntryPriority = "WARN"
	LogEntryPriorityError LogEntryPriority = "ERROR"
)

// LogEntry represents a log entry.
type LogEntry struct {
	// A unique identifier for the log entry.
	ID uint64 `json:"id"`
	// The priority level of the log entry.
	Priority LogEntryPriority `json:"priority"`
	// The module or component of the system that generated the log entry.
	Module string `json:"module"`
	// A message describing the log event.
	Message string `json:"message"`
	// Additional context information related to the log entry, typically including data relevant to the log event.
	Context map[string]string `json:"context"`
	// The timestamp when this log entry was created.
	CreatedAt time.Time `json:"createdAt"`
}
