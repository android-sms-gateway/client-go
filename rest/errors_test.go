package rest_test

import (
	"errors"
	"testing"

	liberr "github.com/android-sms-gateway/client-go/rest"
)

var errSomeOther = errors.New("some other error")

type testCase struct {
	name     string
	err      error
	expected bool
}

func TestIsAPIError(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		{"API error", liberr.ErrAPIError, true},
		{"Client error", liberr.ErrClient, true},
		{"Server error", liberr.ErrServer, true},
		{"Bad request", liberr.ErrBadRequest, true},
		{"Conflict", liberr.ErrConflict, true},
		{"Non-API error", errSomeOther, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := liberr.IsAPIError(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsClientError(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		{"Client error", liberr.ErrClient, true},
		{"Bad request", liberr.ErrBadRequest, true},
		{"Conflict", liberr.ErrConflict, true},
		{"API error", liberr.ErrAPIError, false},
		{"Server error", liberr.ErrServer, false},
		{"Non-client error", errSomeOther, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := liberr.IsClientError(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		{"Server error", liberr.ErrServer, true},
		{"API error", liberr.ErrAPIError, false},
		{"Client error", liberr.ErrClient, false},
		{"Bad request", liberr.ErrBadRequest, false},
		{"Conflict", liberr.ErrConflict, false},
		{"Non-server error", errSomeOther, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := liberr.IsServerError(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsConflict(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		{"Conflict", liberr.ErrConflict, true},
		{"Client error", liberr.ErrClient, false},
		{"API error", liberr.ErrAPIError, false},
		{"Server error", liberr.ErrServer, false},
		{"Bad request", liberr.ErrBadRequest, false},
		{"Non-conflict error", errSomeOther, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := liberr.IsConflict(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestIsBadRequest(t *testing.T) {
	t.Parallel()
	tests := []testCase{
		{"Bad request", liberr.ErrBadRequest, true},
		{"Client error", liberr.ErrClient, false},
		{"API error", liberr.ErrAPIError, false},
		{"Server error", liberr.ErrServer, false},
		{"Conflict", liberr.ErrConflict, false},
		{"Non-bad-request error", errSomeOther, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := liberr.IsBadRequest(tc.err)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
