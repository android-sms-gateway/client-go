package smsgateway_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const testAttachmentEnvelope = "$rsa-oaep-aes-256-gcm$v=1$k=1$QUFBQQ==$QUFBQQ==$QUFBQQ=="

// TestMobilePostInboxRequest_AttachmentPayloadIsOpaqueString pins the wire
// format of the inbox attachment payload: the encrypted envelope is sent
// verbatim and is never base64-encoded.
func TestMobilePostInboxRequest_AttachmentPayloadIsOpaqueString(t *testing.T) {
	t.Parallel()

	size := int64(102400)

	req := smsgateway.MobilePostInboxRequest{
		{
			ID:        "msg-1",
			Type:      smsgateway.IncomingMessageTypeMmsDownloaded,
			Sender:    "+79990001234",
			Content:   testAttachmentEnvelope,
			CreatedAt: time.Date(2026, 10, 21, 10, 30, 0, 0, time.UTC),
			Attachments: []smsgateway.MobilePostInboxRequestItemAttachment{
				{
					PartID:      1,
					ContentType: "image/jpeg",
					Name:        "$aes-256-cbc/pbkdf2-sha1$i=75000$c2FsdA==$Y2lwaGVy",
					Size:        &size,
					Data:        testAttachmentEnvelope,
				},
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	if want := `"data":"` + testAttachmentEnvelope + `"`; !strings.Contains(string(body), want) {
		t.Errorf("attachment payload is not sent verbatim, want %s in %s", want, body)
	}

	var decoded smsgateway.MobilePostInboxRequest
	if err = json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if got := decoded[0].Attachments[0].Data; got != testAttachmentEnvelope {
		t.Errorf("attachment payload round-trip = %q, want %q", got, testAttachmentEnvelope)
	}
}

// TestMobilePostInboxRequest_AttachmentPayloadIsNotDecoded pins that the
// payload is opaque: a base64-looking value is kept as-is and never decoded.
func TestMobilePostInboxRequest_AttachmentPayloadIsNotDecoded(t *testing.T) {
	t.Parallel()

	const payload = "SGVsbG8gV29ybGQh"

	body := []byte(`[{"id":"msg-1","type":"MMS_DOWNLOADED","sender":"+79990001234",` +
		`"content":"encrypted","isEncrypted":true,"createdAt":"2026-10-21T10:30:00Z",` +
		`"attachments":[{"partId":1,"contentType":"image/jpeg","name":"part001","data":"` + payload + `"}]}]`)

	var decoded smsgateway.MobilePostInboxRequest
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if got := decoded[0].Attachments[0].Data; got != payload {
		t.Errorf("attachment payload = %q, want %q", got, payload)
	}
}
