package smsgateway_test

import (
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestDeviceSettings_Validate(t *testing.T) {
	tests := []struct {
		name     string
		settings smsgateway.DeviceSettings
		wantErr  bool
	}{
		{
			name: "nil messages",
			settings: smsgateway.DeviceSettings{
				Messages: nil,
			},
			wantErr: false,
		},
		{
			name: "valid messages",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					SendIntervalMin:  nil,
					SendIntervalMax:  nil,
					LimitPeriod:      nil,
					LimitValue:       nil,
					SimSelectionMode: nil,
					LogLifetimeDays:  nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.settings.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DeviceSettings.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSettingsMessages_Validate(t *testing.T) {
	tests := []struct {
		name     string
		messages smsgateway.SettingsMessages
		wantErr  bool
	}{
		{
			name: "valid settings",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin:  nil,
				SendIntervalMax:  nil,
				LimitPeriod:      nil,
				LimitValue:       nil,
				SimSelectionMode: nil,
				LogLifetimeDays:  nil,
			},
			wantErr: false,
		},
		{
			name: "invalid send interval",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: ptr(1),
				SendIntervalMax: ptr(0), // Invalid: should be >= 1
			},
			wantErr: true,
		},
		{
			name: "invalid send interval max < min",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: ptr(2),
				SendIntervalMax: ptr(1), // Invalid: max < min
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.messages.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SettingsMessages.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create a pointer to an int.
func ptr(i int) *int {
	return &i
}
