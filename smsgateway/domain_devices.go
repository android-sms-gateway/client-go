package smsgateway

import "time"

// Device represents a device registered on the server.
type Device struct {
	ID        string     `json:"id"                  example:"PyDmBQZZXYmyxMwED8Fzy"` // Device ID, read only.
	Name      string     `json:"name"                example:"My Device"`             // Device name.
	CreatedAt time.Time  `json:"createdAt"           example:"2020-01-01T00:00:00Z"`  // Time at which the device was created, read only.
	UpdatedAt time.Time  `json:"updatedAt"           example:"2020-01-01T00:00:00Z"`  // Time at which the device was last updated, read only.
	DeletedAt *time.Time `json:"deletedAt,omitempty" example:"2020-01-01T00:00:00Z"`  // Time at which the device was deleted, read only.

	LastSeen time.Time `json:"lastSeen" example:"2020-01-01T00:00:00Z"` // Time at which the device was last seen, read only.

	SimCards []SimCard `json:"simCards,omitempty"` // List of SIM cards in the device.
}

// SimCard represents a SIM card in an Android device.
type SimCard struct {
	SlotIndex   int     `json:"slotIndex"`             // 0-based physical slot index.
	SimNumber   int     `json:"simNumber"`             // 1-based slot number (1, 2, or 3).
	PhoneNumber *string `json:"phoneNumber,omitempty"` // Phone number associated with the SIM.
	CarrierName *string `json:"carrierName,omitempty"` // Carrier/network operator name (may be null).
	ICCID       *string `json:"iccid,omitempty"`       // Integrated Circuit Card Identifier (may be null).
}
