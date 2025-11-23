package smsgateway_test

import (
	"net/http"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestConfig_WithClient(t *testing.T) {
	client := &http.Client{}
	tests := []struct {
		name     string
		client   *http.Client
		expected *http.Client
	}{
		{
			name:     "with custom client",
			client:   client,
			expected: client,
		},
		{
			name:     "with nil client should use default",
			client:   nil,
			expected: http.DefaultClient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := smsgateway.Config{}
			result := config.WithClient(tt.client)

			if result.Client != tt.expected {
				t.Errorf("WithClient() client = %v, want %v", result.Client, tt.expected)
			}
		})
	}
}

func TestConfig_WithBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		expected string
	}{
		{
			name:     "with custom base URL",
			baseURL:  "https://custom.example.com/api",
			expected: "https://custom.example.com/api",
		},
		{
			name:     "with empty base URL should use default",
			baseURL:  "",
			expected: smsgateway.BaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := smsgateway.Config{}
			result := config.WithBaseURL(tt.baseURL)

			if result.BaseURL != tt.expected {
				t.Errorf("WithBaseURL() baseURL = %v, want %v", result.BaseURL, tt.expected)
			}
		})
	}
}

func TestConfig_WithJWTAuth(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "with JWT token",
			token:    "jwt.token.here",
			expected: "jwt.token.here",
		},
		{
			name:     "with empty token",
			token:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := smsgateway.Config{}
			result := config.WithJWTAuth(tt.token)

			if result.Token != tt.expected {
				t.Errorf("WithJWTAuth() token = %v, want %v", result.Token, tt.expected)
			}
		})
	}
}

func TestConfig_WithBasicAuth(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		expected struct {
			user     string
			password string
		}
	}{
		{
			name:     "with basic auth credentials",
			user:     "testuser",
			password: "testpass",
			expected: struct {
				user     string
				password string
			}{user: "testuser", password: "testpass"},
		},
		{
			name:     "with empty credentials",
			user:     "",
			password: "",
			expected: struct {
				user     string
				password string
			}{user: "", password: ""},
		},
		{
			name:     "with user only",
			user:     "testuser",
			password: "",
			expected: struct {
				user     string
				password string
			}{user: "testuser", password: ""},
		},
		{
			name:     "with password only",
			user:     "",
			password: "testpass",
			expected: struct {
				user     string
				password string
			}{user: "", password: "testpass"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := smsgateway.Config{}
			result := config.WithBasicAuth(tt.user, tt.password)

			if result.User != tt.expected.user {
				t.Errorf("WithBasicAuth() user = %v, want %v", result.User, tt.expected.user)
			}
			if result.Password != tt.expected.password {
				t.Errorf("WithBasicAuth() password = %v, want %v", result.Password, tt.expected.password)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      smsgateway.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with basic auth",
			config: smsgateway.Config{
				User:     "testuser",
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name: "valid config with JWT token",
			config: smsgateway.Config{
				Token: "jwt.token.here",
			},
			expectError: false,
		},
		{
			name: "valid config with both basic auth and token",
			config: smsgateway.Config{
				User:     "testuser",
				Password: "testpass",
				Token:    "jwt.token.here",
			},
			expectError: false,
		},
		{
			name: "valid config with user only",
			config: smsgateway.Config{
				User: "testuser",
			},
			expectError: false,
		},
		{
			name: "valid config with password only",
			config: smsgateway.Config{
				Password: "testpass",
			},
			expectError: false,
		},
		{
			name:        "invalid config with no auth credentials",
			config:      smsgateway.Config{},
			expectError: true,
			errorMsg:    "invalid config: missing auth credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestConfig_Chaining(t *testing.T) {
	// Test that methods can be chained together
	customClient := &http.Client{}
	config := smsgateway.Config{}.
		WithClient(customClient).
		WithBaseURL("https://custom.example.com/api").
		WithJWTAuth("jwt.token.here").
		WithBasicAuth("user", "pass")

	if config.Client != customClient {
		t.Errorf("Chained WithClient() failed, got %v, want %v", config.Client, customClient)
	}
	if config.BaseURL != "https://custom.example.com/api" {
		t.Errorf("Chained WithBaseURL() failed, got %v, want %v", config.BaseURL, "https://custom.example.com/api")
	}
	if config.Token != "jwt.token.here" {
		t.Errorf("Chained WithJWTAuth() failed, got %v, want %v", config.Token, "jwt.token.here")
	}
	if config.User != "user" {
		t.Errorf("Chained WithBasicAuth() user failed, got %v, want %v", config.User, "user")
	}
	if config.Password != "pass" {
		t.Errorf("Chained WithBasicAuth() password failed, got %v, want %v", config.Password, "pass")
	}
}
