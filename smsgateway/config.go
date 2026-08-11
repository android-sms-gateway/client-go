package smsgateway

import (
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Client         *http.Client  // Optional HTTP Client, defaults to `http.DefaultClient`
	BaseURL        string        // Optional base URL, defaults to `https://api.sms-gate.app/3rdparty/v1`
	User           string        // Basic Auth username
	Password       string        // Basic Auth password
	Token          string        // Bearer token, has priority over Basic Auth
	DeviceCacheTTL time.Duration // Optional TTL for cached device listing entries, defaults to 60s
	Encryptor      Encryptor     // Optional E2E Encryptor; nil sends plaintext with no device lookup
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
// These are used when no Bearer token is set (see [WithJWTAuth]).
func (c Config) WithBasicAuth(user, password string) Config {
	c.User = user
	c.Password = password
	return c
}

// WithDeviceCacheTTL sets the TTL for cached device listing entries.
// A non-positive value falls back to the default (60s, matching the TS SDK).
func (c Config) WithDeviceCacheTTL(ttl time.Duration) Config {
	c.DeviceCacheTTL = ttl
	return c
}

// WithEncryptor sets the E2E Encryptor for the API client.
// Without an encryptor, messages are sent plaintext and no device listing
// lookup happens. With one installed, E2E is applied by default to keyed
// devices; see [WithE2EEncryption] for per-send control.
func (c Config) WithEncryptor(encryptor Encryptor) Config {
	c.Encryptor = encryptor
	return c
}

func (c Config) Validate() error {
	if c.User == "" && c.Password == "" && c.Token == "" {
		return fmt.Errorf("%w: missing auth credentials", ErrInvalidConfig)
	}
	return nil
}
