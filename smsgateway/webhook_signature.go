package smsgateway

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

const (
	// WebhookSignatureHeader is the HTTP header carrying the hexadecimal HMAC-SHA256
	// signature of a webhook request body.
	WebhookSignatureHeader = "X-Signature"

	// WebhookTimestampHeader is the HTTP header carrying the Unix timestamp, in seconds,
	// used to sign a webhook request body.
	WebhookTimestampHeader = "X-Timestamp"
)

// VerifyWebhookSignature reports whether signature authenticates body as a webhook
// request signed by SMSGate.
//
// The signature is the hexadecimal HMAC-SHA256 digest of the raw request body followed
// by the timestamp string, keyed with the signing key. The body must be the bytes as
// received, before any JSON decoding, and the timestamp must be the value of the
// X-Timestamp header.
//
// Verification fails when the signing key, the timestamp or the signature is empty,
// when the timestamp is not a Unix timestamp in seconds, or when the signature does
// not match the body.
//
// Use VerifyWebhookSignatureWithin to also reject replayed requests.
func VerifyWebhookSignature(secretKey string, body []byte, timestamp, signature string) bool {
	if secretKey == "" || timestamp == "" || signature == "" {
		return false
	}

	if _, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
		return false
	}

	return hmac.Equal([]byte(webhookSignature(secretKey, body, timestamp)), []byte(signature))
}

// VerifyWebhookSignatureWithin reports whether signature authenticates body as a webhook
// request signed by SMSGate and whether the request is fresh enough to be accepted.
//
// It applies the same checks as VerifyWebhookSignature and additionally rejects a request
// whose timestamp is more than maxAge away from the current time, in either direction.
// A non-positive maxAge disables the freshness check.
func VerifyWebhookSignatureWithin(
	secretKey string,
	body []byte,
	timestamp, signature string,
	maxAge time.Duration,
) bool {
	if !VerifyWebhookSignature(secretKey, body, timestamp, signature) {
		return false
	}

	if maxAge <= 0 {
		return true
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	age := time.Since(time.Unix(ts, 0))

	return age <= maxAge && age >= -maxAge
}

// webhookSignature computes the hexadecimal HMAC-SHA256 digest of the body followed by
// the timestamp, keyed with the signing key.
func webhookSignature(secretKey string, body []byte, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(body)
	mac.Write([]byte(timestamp))

	return hex.EncodeToString(mac.Sum(nil))
}
