package ca

import "fmt"

// PostCSRRequest represents a request to post a Certificate Signing Request (CSR).
type PostCSRRequest struct {
	// Type is the type of the CSR. By default, it is set to "webhook".
	Type CSRType `json:"type,omitempty" default:"webhook"`
	// Content contains the CSR content and is required.
	Content string `json:"content" validate:"required,max=16384,startswith=-----BEGIN CERTIFICATE REQUEST-----"`
	// Metadata includes additional metadata related to the CSR.
	Metadata map[string]string `json:"metadata,omitempty" validate:"dive,keys,max=64,endkeys,max=256"`
}

// Validate checks if the request is valid.
func (c PostCSRRequest) Validate() error {
	if c.Type != "" && !IsValidCSRType(c.Type) {
		return fmt.Errorf("%w: invalid csr type: %s", ErrValidationFailed, c.Type)
	}

	return nil
}
