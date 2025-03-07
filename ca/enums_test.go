package ca_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/ca"
)

func TestIsValidCSRType(t *testing.T) {
	tests := []struct {
		name    string
		csrType ca.CSRType
		want    bool
	}{
		{
			name:    "webhook type",
			csrType: ca.CSRTypeWebhook,
			want:    true,
		},
		{
			name:    "private server type",
			csrType: ca.CSRTypePrivateServer,
			want:    true,
		},
		{
			name:    "invalid type",
			csrType: "invalid_type",
			want:    false,
		},
		{
			name:    "empty type",
			csrType: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ca.IsValidCSRType(tt.csrType); got != tt.want {
				t.Errorf("IsValidCSRType() = %v, want %v", got, tt.want)
			}
		})
	}
}
