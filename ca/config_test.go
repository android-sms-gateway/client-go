package ca_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/android-sms-gateway/client-go/ca"
)

//nolint:gochecknoglobals // constant
var customHTTPClient = &http.Client{}

func TestConfig_Client(t *testing.T) {
	tests := []struct {
		name   string
		option ca.Option
		want   *http.Client
	}{
		{
			name:   "With Client",
			option: ca.WithClient(customHTTPClient),
			want:   customHTTPClient,
		},
		{
			name:   "Without Client",
			option: ca.WithClient(nil),
			want:   http.DefaultClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ca.Config{}
			tt.option(&c)

			if got := c.Client(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.Client() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_BaseURL(t *testing.T) {
	tests := []struct {
		name   string
		option ca.Option
		want   string
	}{
		{
			name:   "With Base URL",
			option: ca.WithBaseURL("https://example.com"),
			want:   "https://example.com",
		},
		{
			name:   "Without Base URL",
			option: ca.WithBaseURL(""),
			want:   ca.BaseURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ca.Config{}
			tt.option(&c)

			if got := c.BaseURL(); got != tt.want {
				t.Errorf("Config.BaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
