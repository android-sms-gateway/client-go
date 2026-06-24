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
		{
			name: "messages with validation error",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					SendIntervalMin: ptr(5),
					SendIntervalMax: ptr(1),
				},
			},
			wantErr: true,
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
		{
			name: "work hours enabled without start",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   nil,
				WorkHoursEnd:     ptr("18:00"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled without end",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     nil,
			},
			wantErr: true,
		},
		{
			name: "work hours enabled invalid start format",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("9:00"),
				WorkHoursEnd:     ptr("18:00"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled invalid end format",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     ptr("6:00 PM"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled invalid hour value",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("24:00"),
				WorkHoursEnd:     ptr("18:00"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled invalid minute value",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     ptr("18:60"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled seconds are not allowed",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     ptr("18:00:00"),
			},
			wantErr: true,
		},
		{
			name: "work hours enabled valid",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     ptr("18:00"),
			},
			wantErr: false,
		},
		{
			name: "work hours enabled overnight window",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(true),
				WorkHoursStart:   ptr("22:00"),
				WorkHoursEnd:     ptr("06:00"),
			},
			wantErr: false,
		},
		{
			name: "work hours disabled ignores start/end",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: ptr(false),
				WorkHoursStart:   nil,
				WorkHoursEnd:     nil,
			},
			wantErr: false,
		},
		{
			name: "send interval min only",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: ptr(1),
				SendIntervalMax: nil,
			},
			wantErr: false,
		},
		{
			name: "send interval max only",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: nil,
				SendIntervalMax: ptr(1),
			},
			wantErr: false,
		},
		{
			name: "equal send intervals",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: ptr(5),
				SendIntervalMax: ptr(5),
			},
			wantErr: false,
		},
		{
			name: "valid send interval range",
			messages: smsgateway.SettingsMessages{
				SendIntervalMin: ptr(1),
				SendIntervalMax: ptr(10),
			},
			wantErr: false,
		},
		{
			name: "work hours not enabled (nil)",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: nil,
				WorkHoursStart:   nil,
				WorkHoursEnd:     nil,
			},
			wantErr: false,
		},
		{
			name: "work hours not enabled but start set",
			messages: smsgateway.SettingsMessages{
				WorkHoursEnabled: nil,
				WorkHoursStart:   ptr("09:00"),
				WorkHoursEnd:     nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.messages.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SettingsMessages.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ptr[T any](i T) *T {
	return &i
}
