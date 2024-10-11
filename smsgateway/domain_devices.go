package smsgateway

import "time"

// Device
type Device struct {
	ID        string     `json:"id" example:"PyDmBQZZXYmyxMwED8Fzy"`                 // ID
	Name      string     `json:"name" example:"My Device"`                           // Name
	CreatedAt time.Time  `json:"createdAt" example:"2020-01-01T00:00:00Z"`           // Created at (read only)
	UpdatedAt time.Time  `json:"updatedAt" example:"2020-01-01T00:00:00Z"`           // Updated at (read only)
	DeletedAt *time.Time `json:"deletedAt,omitempty" example:"2020-01-01T00:00:00Z"` // Deleted at (read only)

	LastSeen time.Time `json:"lastSeen" example:"2020-01-01T00:00:00Z"` // Last seen at (read only)
}
