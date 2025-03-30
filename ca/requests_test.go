package ca_test

import (
	"errors"
	"testing"

	"github.com/android-sms-gateway/client-go/ca"
)

func TestPostCSRRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request ca.PostCSRRequest
		wantErr bool
		err     error
	}{
		{
			name: "empty type should be valid",
			request: ca.PostCSRRequest{
				Type:    "",
				Content: "-----BEGIN CERTIFICATE REQUEST-----TEST",
			},
			wantErr: false,
		},
		{
			name: "valid webhook type should be valid",
			request: ca.PostCSRRequest{
				Type:    ca.CSRTypeWebhook,
				Content: "-----BEGIN CERTIFICATE REQUEST-----TEST",
			},
			wantErr: false,
		},
		{
			name: "valid private_server type should be valid",
			request: ca.PostCSRRequest{
				Type:    ca.CSRTypePrivateServer,
				Content: "-----BEGIN CERTIFICATE REQUEST-----TEST",
			},
			wantErr: false,
		},
		{
			name: "invalid type should return error",
			request: ca.PostCSRRequest{
				Type:    "invalid_type",
				Content: "-----BEGIN CERTIFICATE REQUEST-----TEST",
			},
			wantErr: true,
			err:     ca.ErrValidationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("Expected error message '%s', got '%s'", tt.err, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
