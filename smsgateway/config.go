package smsgateway

import (
	"fmt"
	"net/http"
)

type Config struct {
	Client   *http.Client // Optional HTTP Client, defaults to `http.DefaultClient`
	BaseURL  string       // Optional base URL, defaults to `https://api.sms-gate.app/3rdparty/v1`
	User     string       // Basic Auth username
	Password string       // Basic Auth password
	Token    string       // Bearer token, has priority over Basic Auth
}

// WithClient sets the HTTP client for the API client.
// If the client is nil, it defaults to `http.DefaultClient`.
// This is useful for testing or custom HTTP clients.
func (c Config) WithClient(client *http.Client) Config {
	if client == nil {
		client = http.DefaultClient
	}
	c.Client = client
	return c
}

// WithBaseURL sets the base URL for the API client.
// If the base URL is empty, it defaults to the constant `BaseURL`.
// This is useful for setting a custom base URL for the API client.
func (c Config) WithBaseURL(baseURL string) Config {
	if baseURL == "" {
		baseURL = BaseURL
	}
	c.BaseURL = baseURL
	return c
}

// WithJWTAuth sets the Bearer token for the API client.
// This is useful for setting a custom Bearer token for the API client.
// If the token is empty, it defaults to an empty string.
func (c Config) WithJWTAuth(token string) Config {
	c.Token = token
	return c
}

// WithBasicAuth sets the Basic Auth credentials for the API client.
// This is useful for setting custom Basic Auth credentials for the API client.
// If the user or password is empty, it defaults to an empty string.
func (c Config) WithBasicAuth(user, password string) Config {
	c.User = user
	c.Password = password
	return c
}

func (c Config) Validate() error {
	if c.User == "" && c.Password == "" && c.Token == "" {
		return fmt.Errorf("%w: missing auth credentials", ErrInvalidConfig)
	}
	return nil
}
