package smsgateway_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestSendOptions_ToURLValues(t *testing.T) {
	tests := []struct {
		name     string
		options  []smsgateway.SendOption
		expected url.Values
	}{
		{
			name:     "No options",
			options:  []smsgateway.SendOption{},
			expected: url.Values{},
		},
		{
			name: "Only skipPhoneValidation=true",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(true),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"true"},
			},
		},
		{
			name: "Only skipPhoneValidation=false",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(false),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"false"},
			},
		},
		{
			name: "Only deviceActiveWithin",
			options: []smsgateway.SendOption{
				smsgateway.WithDeviceActiveWithin(24),
			},
			expected: url.Values{
				"deviceActiveWithin": []string{"24"},
			},
		},
		{
			name: "Both options",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(true),
				smsgateway.WithDeviceActiveWithin(48),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"true"},
				"deviceActiveWithin":  []string{"48"},
			},
		},
		{
			name: "Different order",
			options: []smsgateway.SendOption{
				smsgateway.WithDeviceActiveWithin(72),
				smsgateway.WithSkipPhoneValidation(false),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"false"},
				"deviceActiveWithin":  []string{"72"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &smsgateway.SendOptions{}
			options.Apply(tt.options...)

			result := options.ToURLValues()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToURLValues() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSendOptionFunctions(t *testing.T) {
	tests := []struct {
		name     string
		options  []smsgateway.SendOption
		expected url.Values
	}{
		{
			name:     "No options",
			options:  []smsgateway.SendOption{},
			expected: url.Values{},
		},
		{
			name: "WithSkipPhoneValidation=true",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(true),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"true"},
			},
		},
		{
			name: "WithSkipPhoneValidation=false",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(false),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"false"},
			},
		},
		{
			name: "WithDeviceActiveWithin",
			options: []smsgateway.SendOption{
				smsgateway.WithDeviceActiveWithin(24),
			},
			expected: url.Values{
				"deviceActiveWithin": []string{"24"},
			},
		},
		{
			name: "Multiple options",
			options: []smsgateway.SendOption{
				smsgateway.WithSkipPhoneValidation(true),
				smsgateway.WithDeviceActiveWithin(48),
			},
			expected: url.Values{
				"skipPhoneValidation": []string{"true"},
				"deviceActiveWithin":  []string{"48"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &smsgateway.SendOptions{}
			options.Apply(tt.options...)

			result := options.ToURLValues()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToURLValues() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
