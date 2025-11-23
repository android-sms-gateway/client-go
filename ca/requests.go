package ca

import "fmt"

// PostCSRRequest represents a request to post a Certificate Signing Request (CSR).
type PostCSRRequest struct {
	Type     CSRType           `json:"type,omitempty"     default:"webhook"`                                                                              // Type is the type of the CSR. By default, it is set to "webhook".
	Content  string            `json:"content"                              validate:"required,max=16384,startswith=-----BEGIN CERTIFICATE REQUEST-----"` // Content contains the CSR content and is required.
	Metadata map[string]string `json:"metadata,omitempty"                   validate:"dive,keys,max=64,endkeys,max=256"`                                  // Metadata includes additional metadata related to the CSR.
}

// Validate checks if the request is valid.
func (c PostCSRRequest) Validate() error {
	if c.Type != "" && !IsValidCSRType(c.Type) {
		return fmt.Errorf("%w: invalid csr type: %s", ErrValidationFailed, c.Type)
	}

	return nil
}
