package ca

type CSRStatus string
type CSRType string

const (
	CSRStatusPending  CSRStatus = "pending"
	CSRStatusApproved CSRStatus = "approved"
	CSRStatusDenied   CSRStatus = "denied"

	CSRStatusDescriptionPending  string = "CSR submitted successfully. Await processing."
	CSRStatusDescriptionApproved string = "CSR approved. The certificate is ready for download."
	CSRStatusDescriptionDenied   string = "CSR denied. Please contact the administrator."

	CSRTypeWebhook       CSRType = "webhook"
	CSRTypePrivateServer CSRType = "private_server"
)

// Description returns a human-readable description for the given CSR status.
func (c CSRStatus) Description() string {
	switch c {
	case CSRStatusPending:
		return CSRStatusDescriptionPending
	case CSRStatusApproved:
		return CSRStatusDescriptionApproved
	case CSRStatusDenied:
		return CSRStatusDescriptionDenied
	default:
		return string(c)
	}
}

var allCSRTypes = map[CSRType]struct{}{
	CSRTypeWebhook:       {},
	CSRTypePrivateServer: {},
}

// IsValidCSRType checks if the given CSR type is valid.
func IsValidCSRType(t CSRType) bool {
	_, ok := allCSRTypes[t]
	return ok
}
