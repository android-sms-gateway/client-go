package smsgateway_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const (
	testSigningKey = "test-signing-key"
	testBody       = `{"payload":{"sender":"6505551212","message":"hello"}}`

	// testSignature is the HMAC-SHA256 digest of testBody and testTimestamp, keyed with
	// testSigningKey. It pins the wire format produced by the Android application.
	testSignature = "8f3a0e4e49ef7cd0f4583ded4600c70ff3ad9b61442d842b95f9607a0a232fa3"

	// testTimestamp is the Unix timestamp, in seconds, used to produce testSignature.
	testTimestamp = "1700000000"
)

// sign computes the SMSGate payload signature: the hexadecimal HMAC-SHA256 digest of the
// raw body followed by the timestamp string.
func sign(key, body, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(body))
	mac.Write([]byte(timestamp))

	return hex.EncodeToString(mac.Sum(nil))
}

// timestampAt formats t as the Unix-seconds X-Timestamp value.
func timestampAt(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

// TestVerifyWebhookSignature tests that a signature is accepted only for the exact body,
// timestamp and signing key it was produced with.
func TestVerifyWebhookSignature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		signingKey string
		body       string
		timestamp  string
		signature  string
		want       bool
	}{
		{
			name:       "Valid signature",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  testSignature,
			want:       true,
		},
		{
			name:       "Signature of another body",
			signingKey: testSigningKey,
			body:       `{"payload":{"sender":"6505551212","message":"tampered"}}`,
			timestamp:  testTimestamp,
			signature:  testSignature,
			want:       false,
		},
		{
			name:       "Wrong signing key",
			signingKey: "other-key",
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  testSignature,
			want:       false,
		},
		{
			name:       "Signature of another timestamp",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  sign(testSigningKey, testBody, "1700000001"),
			want:       false,
		},
		{
			name:       "Empty signing key",
			signingKey: "",
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  sign("", testBody, testTimestamp),
			want:       false,
		},
		{
			name:       "Empty signature",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  "",
			want:       false,
		},
		{
			name:       "Empty timestamp",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  "",
			signature:  testSignature,
			want:       false,
		},
		{
			name:       "Uppercase signature",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  testTimestamp,
			signature:  strings.ToUpper(testSignature),
			want:       false,
		},
		{
			name:       "Malformed timestamp",
			signingKey: testSigningKey,
			body:       testBody,
			timestamp:  "not-a-unix-timestamp",
			signature:  sign(testSigningKey, testBody, "not-a-unix-timestamp"),
			want:       false,
		},
		{
			name:       "Empty body is authenticated when the signature matches",
			signingKey: testSigningKey,
			body:       "",
			timestamp:  testTimestamp,
			signature:  sign(testSigningKey, "", testTimestamp),
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := smsgateway.VerifyWebhookSignature(tt.signingKey, []byte(tt.body), tt.timestamp, tt.signature)
			if got != tt.want {
				t.Errorf("VerifyWebhookSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestVerifyWebhookSignatureWithin tests the freshness window of the signature check.
func TestVerifyWebhookSignatureWithin(t *testing.T) {
	t.Parallel()

	const maxAge = 5 * time.Minute

	now := time.Now()
	fresh := timestampAt(now)
	stale := timestampAt(now.Add(-10 * time.Minute))
	future := timestampAt(now.Add(10 * time.Minute))

	tests := []struct {
		name      string
		body      string
		timestamp string
		maxAge    time.Duration
		want      bool
	}{
		{
			name:      "Fresh timestamp",
			body:      testBody,
			timestamp: fresh,
			maxAge:    maxAge,
			want:      true,
		},
		{
			name:      "Stale timestamp",
			body:      testBody,
			timestamp: stale,
			maxAge:    maxAge,
			want:      false,
		},
		{
			name:      "Future timestamp",
			body:      testBody,
			timestamp: future,
			maxAge:    maxAge,
			want:      false,
		},
		{
			name:      "Non-positive max age disables the window",
			body:      testBody,
			timestamp: stale,
			maxAge:    0,
			want:      true,
		},
		{
			name:      "Negative max age disables the window",
			body:      testBody,
			timestamp: stale,
			maxAge:    -maxAge,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			signature := sign(testSigningKey, tt.body, tt.timestamp)

			got := smsgateway.VerifyWebhookSignatureWithin(
				testSigningKey,
				[]byte(tt.body),
				tt.timestamp,
				signature,
				tt.maxAge,
			)
			if got != tt.want {
				t.Errorf("VerifyWebhookSignatureWithin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestVerifyWebhookSignatureWithinRejectsInvalidSignature tests that the freshness window
// does not make an invalid signature acceptable.
func TestVerifyWebhookSignatureWithinRejectsInvalidSignature(t *testing.T) {
	t.Parallel()

	const maxAge = 5 * time.Minute

	fresh := timestampAt(time.Now())

	got := smsgateway.VerifyWebhookSignatureWithin(testSigningKey, []byte(testBody), fresh, testSignature, maxAge)
	if got {
		t.Error("VerifyWebhookSignatureWithin() = true, want false")
	}
}

// TestVerifyWebhookSignatureParityVector tests the signature of the body used by the Rust
// client tests, to keep the SDKs byte-compatible.
func TestVerifyWebhookSignatureParityVector(t *testing.T) {
	t.Parallel()

	const (
		key       = "my-secret-key"
		body      = `{"event":"sms:received","payload":{"message":"Hello"}}`
		timestamp = "1700000000"
		signature = "f7b97c201986d06d2ca1994d0c9802a691aa2ce6b44e11181c54b5b2c849e0ea"
	)

	if !smsgateway.VerifyWebhookSignature(key, []byte(body), timestamp, signature) {
		t.Error("VerifyWebhookSignature() = false, want true")
	}
}
