package ca

import "net/http"

type Option func(*Config)

type Config struct {
	client  *http.Client // Optional HTTP Client, defaults to `http.DefaultClient`
	baseURL string       // Optional base URL, defaults to `https://ca.sms-gate.app/api/v1`
}

func (c Config) Client() *http.Client {
	if c.client == nil {
		return http.DefaultClient
	}
	return c.client
}

func (c Config) BaseURL() string {
	if c.baseURL == "" {
		return BASE_URL
	}
	return c.baseURL
}

func WithClient(client *http.Client) Option {
	return func(c *Config) {
		c.client = client
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.baseURL = baseURL
	}
}
