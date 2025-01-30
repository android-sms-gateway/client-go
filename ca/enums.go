package ca

type CSRStatus string

const (
	CSRStatusPending  CSRStatus = "pending"
	CSRStatusApproved CSRStatus = "approved"
	CSRStatusDenied   CSRStatus = "denied"

	CSRStatusDescriptionPending  string = "CSR submitted successfully. Await processing."
	CSRStatusDescriptionApproved string = "CSR approved. The certificate is ready for download."
	CSRStatusDescriptionDenied   string = "CSR denied. Please contact the administrator."
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
