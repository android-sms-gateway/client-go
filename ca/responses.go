package ca

// PostCSRResponse is a response to a request to post a Certificate Signing Request (CSR).
type PostCSRResponse struct {
	// RequestID is the ID of the request. Can be used to request status.
	RequestID string `json:"request_id"`
	// Status is the status of the requested certificate.
	Status CSRStatus `json:"status"`
	// Message is a human-readable description of the status.
	Message string `json:"message"`
	// Certificate is the certificate issued by the CA. This field is only present
	// if the status is `approved`.
	Certificate string `json:"certificate,omitempty"`
}

type GetCSRStatusResponse = PostCSRResponse
