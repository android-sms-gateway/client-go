package smsgateway_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

// TestPushNotification_PushTokenTypeWire pins the optional push token type on
// upstream notifications: it tells the receiving server which FCM target field
// to use, and omission keeps legacy tokens working.
func TestPushNotification_PushTokenTypeWire(t *testing.T) {
	t.Parallel()

	tokenType := smsgateway.PushTokenTypeFID
	notification := smsgateway.PushNotification{
		Token:     "device-fid",
		Event:     smsgateway.PushMessageEnqueued,
		Data:      map[string]string{"event": "MessageEnqueued"},
		TokenType: &tokenType,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("failed to marshal notification: %v", err)
	}

	if want := `"pushTokenType":"fid"`; !strings.Contains(string(body), want) {
		t.Errorf("pushTokenType is not sent, want %s in %s", want, body)
	}

	var decoded smsgateway.PushNotification
	if err = json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to unmarshal notification: %v", err)
	}

	if decoded.TokenType == nil || *decoded.TokenType != smsgateway.PushTokenTypeFID {
		t.Errorf("pushTokenType round-trip = %v, want %q", decoded.TokenType, smsgateway.PushTokenTypeFID)
	}
}

// TestPushNotification_PushTokenTypeOmitted pins legacy compatibility: without
// the field the receiving server treats the token as a registration token.
func TestPushNotification_PushTokenTypeOmitted(t *testing.T) {
	t.Parallel()

	notification := smsgateway.PushNotification{
		Token:     "legacy-token",
		Event:     smsgateway.PushMessageEnqueued,
		Data:      nil,
		TokenType: nil,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		t.Fatalf("failed to marshal notification: %v", err)
	}

	if strings.Contains(string(body), "pushTokenType") {
		t.Errorf("pushTokenType must be omitted for legacy tokens, got %s", body)
	}
}
