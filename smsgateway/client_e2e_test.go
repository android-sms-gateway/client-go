package smsgateway_test

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

// e2eMock is a small route-aware server for E2E Send/GetState tests.
type e2eMock struct {
	devices         string
	devicesCode     int
	postCode        int
	postBody        string
	stateBody       string
	mu              sync.Mutex
	listDevicesCall int
	postCall        int
}

func (m *e2eMock) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/devices":
			m.mu.Lock()
			m.listDevicesCall++
			devices := m.devices
			code := m.devicesCode
			m.mu.Unlock()
			if code == 0 {
				code = http.StatusOK
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			_, _ = io.WriteString(w, devices)
		case r.Method == http.MethodPost && r.URL.Path == "/messages":
			m.mu.Lock()
			m.postCall++
			body, _ := io.ReadAll(r.Body)
			m.postBody = string(body)
			code, state := m.postCode, m.stateBody
			m.mu.Unlock()
			w.WriteHeader(code)
			if state != "" {
				_, _ = io.WriteString(w, state)
			}
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/messages/"):
			m.mu.Lock()
			state := m.stateBody
			m.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, state)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func newE2EMock(t *testing.T, m *e2eMock) *smsgateway.Client {
	t.Helper()
	server := httptest.NewServer(m.handler())
	t.Cleanup(server.Close)
	return newClient(server.URL)
}

func newE2EMockEnc(t *testing.T, m *e2eMock, encryptor smsgateway.Encryptor) *smsgateway.Client {
	t.Helper()
	server := httptest.NewServer(m.handler())
	t.Cleanup(server.Close)
	return smsgateway.NewClient(smsgateway.Config{
		BaseURL:   server.URL,
		User:      username,
		Password:  password,
		Encryptor: encryptor,
	})
}

func newE2EMockEncTTL(t *testing.T, m *e2eMock, encryptor smsgateway.Encryptor, ttl time.Duration) *smsgateway.Client {
	t.Helper()
	server := httptest.NewServer(m.handler())
	t.Cleanup(server.Close)
	return smsgateway.NewClient(smsgateway.Config{
		BaseURL:        server.URL,
		User:           username,
		Password:       password,
		DeviceCacheTTL: ttl,
		Encryptor:      encryptor,
	})
}

func keyedDeviceJSON(t *testing.T, id string, keyVersion int) string {
	t.Helper()
	v := loadVector(t)
	return `[{"id":"` + id + `","name":"Test","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","lastSeen":"2025-01-01T00:00:00Z","publicKey":"` + v.PublicKeySpkiBase64 + `","keyVersion":` + strconv.Itoa(
		keyVersion,
	) + `}]`
}

func unkeyedDeviceJSON(id string) string {
	return `[{"id":"` + id + `","name":"Test","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","lastSeen":"2025-01-01T00:00:00Z"}]`
}

// decryptE2E decrypts a 7-chunk E2E value using the vector private key and
// returns the plaintext. Mirrors Android's decryptE2E.
func decryptE2E(t *testing.T, value string) string {
	t.Helper()
	v := loadVector(t)
	priv := parseVectorPrivateKey(t, v)

	chunks := strings.Split(value, "$")
	if len(chunks) != 7 {
		t.Fatalf("encrypted value has %d chunks, want 7", len(chunks))
	}
	if chunks[1] != "rsa-oaep-aes-256-gcm" || chunks[2] != "v=1" {
		t.Fatalf("unexpected format chunks: %q %q", chunks[1], chunks[2])
	}
	encKey, err := base64.StdEncoding.DecodeString(chunks[4])
	if err != nil {
		t.Fatalf("decode chunk 4: %v", err)
	}
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)
	if err != nil {
		t.Fatalf("RSA-OAEP decrypt chunk 4: %v", err)
	}
	iv, err := base64.StdEncoding.DecodeString(chunks[5])
	if err != nil {
		t.Fatalf("decode chunk 5: %v", err)
	}
	ctTag, err := base64.StdEncoding.DecodeString(chunks[6])
	if err != nil {
		t.Fatalf("decode chunk 6: %v", err)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("cipher.NewGCM: %v", err)
	}
	plain, err := gcm.Open(nil, iv, ctTag, nil)
	if err != nil {
		t.Fatalf("AES-GCM decrypt chunk 6: %v", err)
	}
	return string(plain)
}

func TestClient_Send_E2E_TextMessage(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 3), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID: "dev-123",
		TextMessage: &smsgateway.TextMessage{
			Text: "The quick brown fox",
		},
		PhoneNumbers: []string{"+79990001234", "999000"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var sent struct {
		IsEncrypted bool `json:"isEncrypted"`
		TextMessage *struct {
			Text string `json:"text"`
		} `json:"textMessage"`
		PhoneNumbers []string `json:"phoneNumbers"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v (body: %s)", err, mock.postBody)
	}

	if !sent.IsEncrypted {
		t.Fatal("isEncrypted must be true for an E2E message")
	}
	if sent.TextMessage == nil || !strings.HasPrefix(sent.TextMessage.Text, "$rsa-oaep-aes-256-gcm$v=1$k=3$") {
		t.Fatalf("text must be E2E-encrypted with k=3, got %q", sent.TextMessage.Text)
	}
	for i, phone := range sent.PhoneNumbers {
		if !strings.HasPrefix(phone, "$rsa-oaep-aes-256-gcm$v=1$k=3$") {
			t.Fatalf("phone %d must be E2E-encrypted with k=3, got %q", i, phone)
		}
	}

	// Round-trip: the device private key must recover the original values.
	if got := decryptE2E(t, sent.TextMessage.Text); got != "The quick brown fox" {
		t.Fatalf("decrypted text = %q", got)
	}
	wantPhones := []string{"+79990001234", "999000"}
	for i, want := range wantPhones {
		if got := decryptE2E(t, sent.PhoneNumbers[i]); got != want {
			t.Fatalf("decrypted phone %d = %q, want %q", i, got, want)
		}
	}

	// Fresh IVs: body and each phone number must use distinct IVs.
	ivs := map[string]struct{}{
		strings.Split(sent.TextMessage.Text, "$")[5]: {},
		strings.Split(sent.PhoneNumbers[0], "$")[5]:  {},
		strings.Split(sent.PhoneNumbers[1], "$")[5]:  {},
	}
	if len(ivs) != 3 {
		t.Fatalf("body and each phone number must have distinct IVs, got %d unique", len(ivs))
	}
}

func TestClient_Send_E2E_DataMessage(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-9", 2), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID: "dev-9",
		DataMessage: &smsgateway.DataMessage{
			Data: "SGVsbG8gV29ybGQh",
			Port: 53739,
		},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var sent struct {
		IsEncrypted bool `json:"isEncrypted"`
		DataMessage *struct {
			Data string `json:"data"`
			Port uint16 `json:"port"`
		} `json:"dataMessage"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v (body: %s)", err, mock.postBody)
	}

	if !sent.IsEncrypted {
		t.Fatal("isEncrypted must be true for an E2E data message")
	}
	if sent.DataMessage == nil || !strings.HasPrefix(sent.DataMessage.Data, "$rsa-oaep-aes-256-gcm$v=1$k=2$") {
		t.Fatalf("data must be E2E-encrypted with k=2, got %+v", sent.DataMessage)
	}
	if sent.DataMessage.Port != 53739 {
		t.Fatalf("port must not be encrypted, got %d", sent.DataMessage.Port)
	}
	if got := decryptE2E(t, sent.DataMessage.Data); got != "SGVsbG8gV29ybGQh" {
		t.Fatalf("decrypted data = %q, want the base64 payload string", got)
	}
}

func TestClient_Send_E2E_DeviceIDRequired(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		TextMessage:  &smsgateway.TextMessage{Text: "secret"},
		PhoneNumbers: []string{"+79990001234"},
	}

	_, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(true))
	if !errors.Is(err, smsgateway.ErrDeviceIDRequired) {
		t.Fatalf("Send error = %v, want ErrDeviceIDRequired", err)
	}
	if mock.postCall != 0 {
		t.Fatal("no request must be sent when DeviceID is missing for an E2E message")
	}
}

func TestClient_Send_E2E_MissingPublicKey_Required(t *testing.T) {
	mock := &e2eMock{devices: unkeyedDeviceJSON("dev-unkeyed"), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-unkeyed",
		TextMessage:  &smsgateway.TextMessage{Text: "secret"},
		PhoneNumbers: []string{"+79990001234"},
	}

	_, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(true))
	if !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("Send error = %v, want ErrE2ENotConfigured", err)
	}
	if mock.postCall != 0 {
		t.Fatal("no request must be sent when the target device has no public key")
	}
}

func TestClient_Send_E2E_MissingKeyVersion_Required(t *testing.T) {
	devices := `[{"id":"dev-nov","name":"Test","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","lastSeen":"2025-01-01T00:00:00Z","publicKey":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8Z9jSr8dYDpZOnnepBCSUBji/3NC2k+sdjWK3pupuXJWN4YZUKDGqt7mk7Gx6kllZ4tNlwASRY7mkTbbFWIhZr2OJnOUZc/9lywTcB3Bof2haEmRhRfEElvlQGdEZof5+RuecZ2c5zIY6Vln8JcYe0KKILcp9huAv3aCE+1uUsiCmTtk7ABPI+HF7oHwsSJPCf0fQ/vCUYBSth0QSMck9jqLaAo1S1I5zN3CAeTSfn/Hk5U+jjKvJf1STWoR78AtRhBvuAEnsfFAbcJHaQ77N7tdUojP6LFMqWBXJFfZmOfKCKNowuYgpZFCyqi8i7TxE+EVdXFOE1sZFxCS6n9HAwIDAQAB"}]`
	mock := &e2eMock{devices: devices, postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-nov",
		TextMessage:  &smsgateway.TextMessage{Text: "secret"},
		PhoneNumbers: []string{"+79990001234"},
	}

	_, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(true))
	if !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("Send error = %v, want ErrE2ENotConfigured", err)
	}
}

func TestClient_Send_BackwardCompatible_NoDeviceID(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMock(t, mock)

	message := smsgateway.Message{
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World!"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if mock.listDevicesCall != 0 {
		t.Fatal("no device listing lookup must happen without a DeviceID")
	}

	// Note: TextMessage is a nested object, so decode into the raw map path.
	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("message without DeviceID must remain plaintext")
	}
}

func TestClient_Send_BackwardCompatible_UnkeyedDevice(t *testing.T) {
	mock := &e2eMock{devices: unkeyedDeviceJSON("dev-unkeyed"), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-unkeyed",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World!"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("unkeyed device must receive plaintext as before")
	}
	if phones, ok := raw["phoneNumbers"].([]any); !ok || len(phones) != 1 || phones[0] != "+1234567890" {
		t.Fatalf("phone numbers must remain plaintext, got %v", raw["phoneNumbers"])
	}
}

func TestClient_Send_E2E_Automatic_KeyedDevice(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 5), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World!"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	// E2E-by-default: with an encryptor installed and no explicit option, a
	// keyed device is encrypted (k=5) and the message flagged.
	var sent struct {
		IsEncrypted bool `json:"isEncrypted"`
		TextMessage *struct {
			Text string `json:"text"`
		} `json:"textMessage"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if !sent.IsEncrypted {
		t.Fatal("keyed device must be E2E-encrypted by default when an encryptor is installed")
	}
	if sent.TextMessage == nil || !strings.HasPrefix(sent.TextMessage.Text, "$rsa-oaep-aes-256-gcm$v=1$k=5$") {
		t.Fatalf("message must be E2E-encrypted with k=5 by default, got %q", sent.TextMessage.Text)
	}
}

func TestClient_Send_DeviceNotFound(t *testing.T) {
	mock := &e2eMock{devices: `[]`, postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-missing",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello"},
		PhoneNumbers: []string{"+1234567890"},
	}

	_, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(true))
	if !errors.Is(err, smsgateway.ErrDeviceNotFound) {
		t.Fatalf("Send error = %v, want ErrDeviceNotFound", err)
	}
}

func TestClient_Send_AutoMode_DeviceListingError(t *testing.T) {
	mock := &e2eMock{devicesCode: http.StatusInternalServerError, postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello"},
		PhoneNumbers: []string{"+1234567890"},
	}

	// In auto mode only ErrDeviceNotFound falls back to plaintext: a
	// device-listing transport/5xx error must propagate and prevent the send.
	_, err := client.Send(context.Background(), message)
	if err == nil {
		t.Fatal("device-listing failure must fail the send in auto mode")
	}
	if mock.postCall != 0 {
		t.Fatal("no request must be sent when the device listing fails")
	}
}

func TestClient_Send_DoesNotMutateCallerMessage(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send 1: %v", err)
	}

	// Send must not write ciphertext through the caller's pointers or slice:
	// the same value must stay re-sendable and plaintext on the caller side.
	if message.TextMessage.Text != "Hello World" {
		t.Fatalf("caller TextMessage mutated to %q", message.TextMessage.Text)
	}
	if message.PhoneNumbers[0] != "+1234567890" {
		t.Fatalf("caller PhoneNumbers mutated to %v", message.PhoneNumbers)
	}
	if message.IsEncrypted {
		t.Fatal("caller IsEncrypted must not be set by Send")
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send 2 with the same value: %v", err)
	}
	assertEncryptedWithKeyVersion(t, mock, 1)
}

func TestClient_Send_CacheAllDevicesFromListing(t *testing.T) {
	d1 := strings.TrimSuffix(keyedDeviceJSON(t, "dev-1", 1), "]")
	d2 := strings.TrimPrefix(keyedDeviceJSON(t, "dev-2", 2), "[")
	mock := &e2eMock{devices: d1 + "," + d2, postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	send := func(deviceID string) {
		t.Helper()
		message := smsgateway.Message{
			DeviceID:     deviceID,
			TextMessage:  &smsgateway.TextMessage{Text: "Hello"},
			PhoneNumbers: []string{"+1234567890"},
		}
		if _, err := client.Send(context.Background(), message); err != nil {
			t.Fatalf("Send to %s: %v", deviceID, err)
		}
	}

	send("dev-1")
	send("dev-2")
	send("dev-1")

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 1 {
		t.Fatalf(
			"one listing must populate the cache for every returned device, listing calls = %d",
			mock.listDevicesCall,
		)
	}
}

func TestClient_Send_PassphrasePreserved(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMock(t, mock)

	passphraseBody := "$aes-256-cbc/pbkdf2-sha1$i=300000$c2FsdA==$Y2lwaGVydGV4dA=="
	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: passphraseBody},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("passphrase-encrypted body must not be re-encrypted with E2E")
	}
	textMessage, ok := raw["textMessage"].(map[string]any)
	if !ok || textMessage["text"] != passphraseBody {
		t.Fatalf("passphrase body must be preserved verbatim, got %v", raw["textMessage"])
	}
	if phones, phonesOK := raw["phoneNumbers"].([]any); !phonesOK || phones[0] != "+1234567890" {
		t.Fatalf("passphrase message phones must be preserved, got %v", raw["phoneNumbers"])
	}
}

func TestClient_Send_IsEncryptedPreserved(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMock(t, mock)

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		IsEncrypted:  true,
		TextMessage:  &smsgateway.TextMessage{Text: "$some-other-format"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] != true {
		t.Fatal("caller-set isEncrypted must be preserved")
	}
	if phones, ok := raw["phoneNumbers"].([]any); !ok || phones[0] != "+1234567890" {
		t.Fatalf("caller-marked encrypted phones must be preserved, got %v", raw["phoneNumbers"])
	}
}

// stubEncryptor is a fake [smsgateway.Encryptor] seam: it returns a fixed value
// regardless of the device or plaintext, so tests can assert where Encrypt is
// applied without depending on the real key material.
type stubEncryptor struct {
	value string
}

func (s stubEncryptor) Encrypt(_ smsgateway.Device, _ string) (string, error) {
	return s.value, nil
}

func TestClient_Send_NoEncryptor_PlaintextShortCircuit(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 3), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMock(t, mock)

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World!"},
		PhoneNumbers: []string{"+1234567890", "+79990001234"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 0 {
		t.Fatalf("no encryptor must short-circuit with zero device listing calls, got %d", mock.listDevicesCall)
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("message without an encryptor must remain plaintext")
	}
	textMessage, ok := raw["textMessage"].(map[string]any)
	if !ok || textMessage["text"] != "Hello World!" {
		t.Fatalf("body must be sent verbatim as plaintext, got %v", raw["textMessage"])
	}
	if phones, phoneOK := raw["phoneNumbers"].([]any); !phoneOK || len(phones) != 2 ||
		phones[0] != "+1234567890" || phones[1] != "+79990001234" {
		t.Fatalf("phones must be sent verbatim as plaintext, got %v", raw["phoneNumbers"])
	}
}

func TestClient_Send_E2E_RequiredNilEncryptor(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 3), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMock(t, mock)

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "secret"},
		PhoneNumbers: []string{"+1234567890"},
	}

	// WithE2EEncryption(true) with no encryptor must fail fast with the typed
	// error, before any device listing or post happens: an explicit encryption
	// demand is never silently downgraded to plaintext.
	_, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(true))
	if !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("Send error = %v, want ErrE2ENotConfigured", err)
	}
	if mock.postCall != 0 {
		t.Fatal("no request must be sent when E2E is required without an encryptor")
	}
	if mock.listDevicesCall != 0 {
		t.Fatalf("required+nil encryptor must not list devices, got %d calls", mock.listDevicesCall)
	}
}

func TestClient_Send_E2E_PrefixedBodyVerbatimWithEncryptor(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 3), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	// Even with an encryptor installed, an already-prefixed body is never
	// double-encrypted: it must be sent verbatim.
	prefixed := "$rsa-oaep-aes-256-gcm$v=1$k=1$c2FsdA==$c2FsdA==$Y2lwaGVydGV4dA=="
	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: prefixed},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 0 {
		t.Fatalf("prefixed body must be sent verbatim without device listing, got %d calls", mock.listDevicesCall)
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("prefixed body must not be flagged as newly encrypted")
	}
	textMessage, ok := raw["textMessage"].(map[string]any)
	if !ok || textMessage["text"] != prefixed {
		t.Fatalf("prefixed body must be preserved verbatim, got %v", raw["textMessage"])
	}
	if phones, phoneOK := raw["phoneNumbers"].([]any); !phoneOK || len(phones) != 1 || phones[0] != "+1234567890" {
		t.Fatalf("prefixed-message phones must be preserved, got %v", raw["phoneNumbers"])
	}
}

func TestClient_Send_E2E_OffWithEncryptor(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 5), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello World!"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), message, smsgateway.WithE2EEncryption(false)); err != nil {
		t.Fatalf("Send: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 0 {
		t.Fatalf("disabled E2E must not list devices even with an encryptor, got %d calls", mock.listDevicesCall)
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(mock.postBody), &raw); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if raw["isEncrypted"] == true {
		t.Fatal("WithE2EEncryption(false) must send plaintext even with an encryptor installed")
	}
	if raw["deviceId"] != "dev-123" {
		t.Fatalf("deviceId must be preserved when E2E is disabled, got %v", raw["deviceId"])
	}
}

func TestClient_Send_E2E_FakeEncryptorSeam(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 3), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, stubEncryptor{value: "ENC:FIXED"})

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "body"},
		PhoneNumbers: []string{"+111", "+222"},
	}

	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var sent struct {
		IsEncrypted bool `json:"isEncrypted"`
		TextMessage *struct {
			Text string `json:"text"`
		} `json:"textMessage"`
		PhoneNumbers []string `json:"phoneNumbers"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if !sent.IsEncrypted {
		t.Fatal("fake-encryptor send must set IsEncrypted")
	}
	if sent.TextMessage.Text != "ENC:FIXED" {
		t.Fatalf("body = %q, want the fake encryptor's fixed value", sent.TextMessage.Text)
	}
	for i, phone := range sent.PhoneNumbers {
		if phone != "ENC:FIXED" {
			t.Fatalf("phone %d = %q, want the fake encryptor's fixed value", i, phone)
		}
	}
}

func TestClient_GetState_EchoesEncryptedPhone(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	message := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "secret"},
		PhoneNumbers: []string{"+79990001234"},
	}
	if _, err := client.Send(context.Background(), message); err != nil {
		t.Fatalf("Send: %v", err)
	}

	var sent struct {
		PhoneNumbers []string `json:"phoneNumbers"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	encryptedPhone := sent.PhoneNumbers[0]
	if !strings.HasPrefix(encryptedPhone, "$rsa-oaep-aes-256-gcm$") {
		t.Fatalf("phone must be E2E-encrypted, got %q", encryptedPhone)
	}

	// The server stores the encrypted phone verbatim and returns it in status.
	mock.mu.Lock()
	mock.stateBody = `{"id":"msg-1","deviceId":"dev-123","state":"Processed","isEncrypted":true,"recipients":[{"phoneNumber":"` + encryptedPhone + `","state":"Processed"}]}`
	mock.mu.Unlock()

	state, err := client.GetState(context.Background(), "msg-1")
	if err != nil {
		t.Fatalf("GetState: %v", err)
	}
	if len(state.Recipients) != 1 {
		t.Fatalf("recipients = %d, want 1", len(state.Recipients))
	}
	if state.Recipients[0].PhoneNumber != encryptedPhone {
		t.Fatalf(
			"status polling must echo the exact encrypted phone string, got %q want %q",
			state.Recipients[0].PhoneNumber,
			encryptedPhone,
		)
	}
}

func TestClient_Send_CacheReusesDeviceListing(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEnc(t, mock, smsgateway.NewEncryptor())

	base := smsgateway.Message{
		DeviceID:     "dev-123",
		TextMessage:  &smsgateway.TextMessage{Text: "Hello"},
		PhoneNumbers: []string{"+1234567890"},
	}

	if _, err := client.Send(context.Background(), base); err != nil {
		t.Fatalf("Send 1: %v", err)
	}
	if _, err := client.Send(context.Background(), base); err != nil {
		t.Fatalf("Send 2: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 1 {
		t.Fatalf("device listing must be cached per deviceId, listing calls = %d", mock.listDevicesCall)
	}
}

func TestClient_Send_DeviceCacheExpiresAfterTTL(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 1), postCode: http.StatusCreated, stateBody: `{}`}
	client := newE2EMockEncTTL(t, mock, smsgateway.NewEncryptor(), 50*time.Millisecond)

	// A fresh message per send: Send encrypts in place, and the text message
	// is a pointer, so reusing one value would already carry the E2E prefix.
	newMessage := func() smsgateway.Message {
		return smsgateway.Message{
			DeviceID:     "dev-123",
			TextMessage:  &smsgateway.TextMessage{Text: "Hello"},
			PhoneNumbers: []string{"+1234567890"},
		}
	}

	if _, err := client.Send(context.Background(), newMessage()); err != nil {
		t.Fatalf("Send 1: %v", err)
	}
	assertEncryptedWithKeyVersion(t, mock, 1)

	// Device key rotation: the server listing now carries keyVersion 2.
	mock.mu.Lock()
	mock.devices = keyedDeviceJSON(t, "dev-123", 2)
	mock.mu.Unlock()

	time.Sleep(80 * time.Millisecond)

	if _, err := client.Send(context.Background(), newMessage()); err != nil {
		t.Fatalf("Send 2: %v", err)
	}
	assertEncryptedWithKeyVersion(t, mock, 2)

	mock.mu.Lock()
	defer mock.mu.Unlock()
	if mock.listDevicesCall != 2 {
		t.Fatalf("listing must be re-fetched after TTL expiry, listing calls = %d", mock.listDevicesCall)
	}
}

// assertEncryptedWithKeyVersion checks that the last posted message body is
// E2E-encrypted with the given keyVersion.
func assertEncryptedWithKeyVersion(t *testing.T, mock *e2eMock, version int) {
	t.Helper()
	mock.mu.Lock()
	defer mock.mu.Unlock()

	var sent struct {
		TextMessage *struct {
			Text string `json:"text"`
		} `json:"textMessage"`
	}
	if err := json.Unmarshal([]byte(mock.postBody), &sent); err != nil {
		t.Fatalf("unmarshal posted body: %v", err)
	}
	if sent.TextMessage == nil {
		t.Fatal("textMessage missing from posted body")
	}

	want := fmt.Sprintf("$rsa-oaep-aes-256-gcm$v=1$k=%d$", version)
	if !strings.HasPrefix(sent.TextMessage.Text, want) {
		got := sent.TextMessage.Text
		if len(got) > 40 {
			got = got[:40]
		}
		t.Fatalf("message must be encrypted with keyVersion %d, got %q", version, got)
	}
}

func TestClient_ListDevices_ParsesPublicKey(t *testing.T) {
	mock := &e2eMock{devices: keyedDeviceJSON(t, "dev-123", 7)}
	client := newE2EMock(t, mock)

	devices, err := client.ListDevices(context.Background())
	if err != nil {
		t.Fatalf("ListDevices: %v", err)
	}
	if len(devices) != 1 {
		t.Fatalf("devices = %d, want 1", len(devices))
	}
	v := loadVector(t)
	if devices[0].PublicKey == nil || *devices[0].PublicKey != v.PublicKeySpkiBase64 {
		t.Fatalf("publicKey must be parsed from the listing, got %v", devices[0].PublicKey)
	}
	if devices[0].KeyVersion == nil || *devices[0].KeyVersion != 7 {
		t.Fatalf("keyVersion must be parsed from the listing, got %v", devices[0].KeyVersion)
	}
}
