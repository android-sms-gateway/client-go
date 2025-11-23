package smsgateway

import "time"

// TokenRequest represents a request to obtain an access token.
//
// The TTL field defines the requested lifetime of the access token in seconds.
// A value of 0 will result in a token with the maximum allowed lifetime.
//
// The Scopes field defines the scopes for which the access token is valid.
// At least one scope must be provided.
type TokenRequest struct {
	TTL    uint64   `json:"ttl,omitempty"`                                         // lifetime of the access token in seconds
	Scopes []string `json:"scopes"        validate:"required,min=1,dive,required"` // scopes for which the access token is valid
}

// TokenResponse represents a response to a TokenRequest.
//
// The ID field contains a unique identifier for the access token.
//
// The TokenType field contains the type of the access token, which is "Bearer".
//
// The AccessToken field contains the actual access token.
//
// The ExpiresAt field contains the time at which the access token is no longer valid.
type TokenResponse struct {
	ID          string    `json:"id"`           // unique identifier for the access token
	TokenType   string    `json:"token_type"`   // type of the access token
	AccessToken string    `json:"access_token"` // actual access token
	ExpiresAt   time.Time `json:"expires_at"`   // time at which the access token is no longer valid
}
