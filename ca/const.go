package ca

const BaseURL = "https://ca.sms-gate.app/api/v1"

//nolint:revive,staticcheck // backward compatibility
const BASE_URL = BaseURL

var (
	//nolint:gochecknoglobals // constant
	emptyHeaders = map[string]string{}
)
