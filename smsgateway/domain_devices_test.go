package smsgateway_test

import (
	"encoding/json"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func TestDeviceStateRoundTrip(t *testing.T) {
	t.Parallel()

	device := smsgateway.Device{
		ID:    "device-id",
		State: smsgateway.DeviceStateDisabled,
	}

	encoded, err := json.Marshal(device)
	if err != nil {
		t.Fatal(err)
	}

	var decoded map[string]any
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded["state"] != string(smsgateway.DeviceStateDisabled) {
		t.Fatalf("expected state %q, got %v", smsgateway.DeviceStateDisabled, decoded["state"])
	}

	var fromWire smsgateway.Device
	if err = json.Unmarshal([]byte(`{"state":"Active"}`), &fromWire); err != nil {
		t.Fatal(err)
	}

	if fromWire.State != smsgateway.DeviceStateActive {
		t.Fatalf("expected state %q, got %q", smsgateway.DeviceStateActive, fromWire.State)
	}
}
