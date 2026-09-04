package smsgateway_test

import (
	"errors"
	"math"
	"net/url"
	"reflect"
	"testing"
	"time"

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
		{
			name: "deviceActiveWithin at MaxInt32",
			options: []smsgateway.SendOption{
				smsgateway.WithDeviceActiveWithin(math.MaxInt32),
			},
			expected: url.Values{
				"deviceActiveWithin": []string{"2147483647"},
			},
		},
		{
			name: "deviceActiveWithin above MaxInt32 is clamped",
			options: []smsgateway.SendOption{
				smsgateway.WithDeviceActiveWithin(math.MaxInt32 + 1),
			},
			expected: url.Values{
				"deviceActiveWithin": []string{"2147483647"},
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

func TestSendOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options *smsgateway.SendOptions
		wantErr bool
	}{
		{
			name:    "No options",
			options: &smsgateway.SendOptions{},
			wantErr: false,
		},
		{
			name: "deviceActiveWithin at minimum",
			options: &smsgateway.SendOptions{
				DeviceActiveWithin: ptr(1),
			},
			wantErr: false,
		},
		{
			name: "deviceActiveWithin at maximum",
			options: &smsgateway.SendOptions{
				DeviceActiveWithin: ptr(math.MaxInt32),
			},
			wantErr: false,
		},
		{
			name: "deviceActiveWithin above MaxInt32",
			options: &smsgateway.SendOptions{
				DeviceActiveWithin: ptr(math.MaxInt32 + 1),
			},
			wantErr: true,
		},
		{
			name: "deviceActiveWithin below minimum",
			options: &smsgateway.SendOptions{
				DeviceActiveWithin: ptr(0),
			},
			wantErr: true,
		},
		{
			name: "deviceActiveWithin negative",
			options: &smsgateway.SendOptions{
				DeviceActiveWithin: ptr(-1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListInboxOptions_ToURLValues(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	deviceID := "abc123def456ghi789j"

	tests := []struct {
		name     string
		options  smsgateway.ListInboxOptions
		expected url.Values
	}{
		{
			name:     "No options",
			options:  smsgateway.ListInboxOptions{},
			expected: url.Values{},
		},
		{
			name: "Type set",
			options: smsgateway.ListInboxOptions{
				Type: ptr(smsgateway.IncomingMessageTypeSMS),
			},
			expected: url.Values{
				"type": {"SMS"},
			},
		},
		{
			name: "Limit set",
			options: smsgateway.ListInboxOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Limit: ptr(10),
				},
			},
			expected: url.Values{
				"limit": {"10"},
			},
		},
		{
			name: "Offset set",
			options: smsgateway.ListInboxOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Offset: ptr(20),
				},
			},
			expected: url.Values{
				"offset": {"20"},
			},
		},
		{
			name: "From set",
			options: smsgateway.ListInboxOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
				},
			},
			expected: url.Values{
				"from": {from.Format(time.RFC3339)},
			},
		},
		{
			name: "To set",
			options: smsgateway.ListInboxOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					To: &to,
				},
			},
			expected: url.Values{
				"to": {to.Format(time.RFC3339)},
			},
		},
		{
			name: "DeviceID set",
			options: smsgateway.ListInboxOptions{
				DeviceID: &deviceID,
			},
			expected: url.Values{
				"deviceId": {deviceID},
			},
		},
		{
			name: "All fields set",
			options: smsgateway.ListInboxOptions{
				Type: ptr(smsgateway.IncomingMessageTypeSMS),
				PaginationOptions: smsgateway.PaginationOptions{
					Limit:  ptr(50),
					Offset: ptr(10),
				},
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
					To:   &to,
				},
				DeviceID: &deviceID,
			},
			expected: url.Values{
				"type":     {"SMS"},
				"limit":    {"50"},
				"offset":   {"10"},
				"from":     {from.Format(time.RFC3339)},
				"to":       {to.Format(time.RFC3339)},
				"deviceId": {deviceID},
			},
		},
		{
			name: "Limit zero",
			options: smsgateway.ListInboxOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Limit: ptr(0),
				},
			},
			expected: url.Values{
				"limit": {"0"},
			},
		},
		{
			name: "Offset zero",
			options: smsgateway.ListInboxOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Offset: ptr(0),
				},
			},
			expected: url.Values{
				"offset": {"0"},
			},
		},
		{
			name: "DATA_SMS type",
			options: smsgateway.ListInboxOptions{
				Type: ptr(smsgateway.IncomingMessageTypeDataSMS),
			},
			expected: url.Values{
				"type": {"DATA_SMS"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ToURLValues()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToURLValues() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestListMessagesOptions_Validate(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name    string
		options smsgateway.ListMessagesOptions
		wantErr bool
	}{
		{
			name:    "Both From and To nil",
			options: smsgateway.ListMessagesOptions{},
			wantErr: false,
		},
		{
			name: "From before To",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
					To:   &to,
				},
			},
			wantErr: false,
		},
		{
			name: "From after To",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &to,
					To:   &from,
				},
			},
			wantErr: true,
		},
		{
			name: "From equal to To",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
					To:   &from,
				},
			},
			wantErr: false,
		},
		{
			name: "Only From set",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
				},
			},
			wantErr: false,
		},
		{
			name: "Only To set",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					To: &to,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}

				if !errors.Is(err, smsgateway.ErrValidationFailed) {
					t.Errorf("Validate() error = %v, want error type %v", err, smsgateway.ErrValidationFailed)
				}
			} else if err != nil {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListMessagesOptions_ToURLValues(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	state := "Sent"
	deviceID := "abc123def456ghi789j"

	tests := []struct {
		name     string
		options  smsgateway.ListMessagesOptions
		expected url.Values
	}{
		{
			name:     "No options",
			options:  smsgateway.ListMessagesOptions{},
			expected: url.Values{},
		},
		{
			name: "From set",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
				},
			},
			expected: url.Values{
				"from": {from.Format(time.RFC3339)},
			},
		},
		{
			name: "To set",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					To: &to,
				},
			},
			expected: url.Values{
				"to": {to.Format(time.RFC3339)},
			},
		},
		{
			name: "State set",
			options: smsgateway.ListMessagesOptions{
				State: &state,
			},
			expected: url.Values{
				"state": {"Sent"},
			},
		},
		{
			name: "DeviceID set",
			options: smsgateway.ListMessagesOptions{
				DeviceID: &deviceID,
			},
			expected: url.Values{
				"deviceId": {deviceID},
			},
		},
		{
			name: "Limit set",
			options: smsgateway.ListMessagesOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Limit: ptr(25),
				},
			},
			expected: url.Values{
				"limit": {"25"},
			},
		},
		{
			name: "Offset set",
			options: smsgateway.ListMessagesOptions{
				PaginationOptions: smsgateway.PaginationOptions{
					Offset: ptr(5),
				},
			},
			expected: url.Values{
				"offset": {"5"},
			},
		},
		{
			name: "IncludeContent true",
			options: smsgateway.ListMessagesOptions{
				IncludeContent: ptr(true),
			},
			expected: url.Values{
				"includeContent": {"true"},
			},
		},
		{
			name: "IncludeContent false",
			options: smsgateway.ListMessagesOptions{
				IncludeContent: ptr(false),
			},
			expected: url.Values{
				"includeContent": {"false"},
			},
		},
		{
			name: "Sort created_at ascending",
			options: smsgateway.ListMessagesOptions{
				Sort: ptr(smsgateway.CreatedAtAscending),
			},
			expected: url.Values{
				"sort": {"created_at"},
			},
		},
		{
			name: "Sort created_at descending",
			options: smsgateway.ListMessagesOptions{
				Sort: ptr(smsgateway.CreatedAtDescending),
			},
			expected: url.Values{
				"sort": {"-created_at"},
			},
		},
		{
			name: "All fields set",
			options: smsgateway.ListMessagesOptions{
				DatePeriodOptions: smsgateway.DatePeriodOptions{
					From: &from,
					To:   &to,
				},
				State:    &state,
				DeviceID: &deviceID,
				PaginationOptions: smsgateway.PaginationOptions{
					Limit:  ptr(100),
					Offset: ptr(0),
				},
				IncludeContent: ptr(true),
				Sort:           ptr(smsgateway.CreatedAtDescending),
			},
			expected: url.Values{
				"from":           {from.Format(time.RFC3339)},
				"to":             {to.Format(time.RFC3339)},
				"state":          {"Sent"},
				"deviceId":       {deviceID},
				"limit":          {"100"},
				"offset":         {"0"},
				"includeContent": {"true"},
				"sort":           {"-created_at"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.ToURLValues()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToURLValues() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
