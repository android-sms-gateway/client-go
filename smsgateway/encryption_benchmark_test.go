package smsgateway_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

// Host-side benchmarks for the E2E hybrid scheme (RSA-2048 OAEP-SHA256 +
// AES-256-GCM). These give a desktop-JCA baseline only; the authoritative
// numbers for the acceptance criteria come from the Android instrumented
// benchmark (android-app .../E2EEncryptionBenchmarkTest.kt) run on API 21/22
// and API 28+ emulators (see docs/plan/e2e-encryption/benchmarks.md).
//
// Run: go test -run=^$ -bench=E2E -benchmem ./smsgateway/

// BenchmarkE2EKeyGen measures RSA-2048 keypair generation (Android
// E2EKeyService.rotateKey: 2048-bit Keystore key on API 23+, software on
// API 21-22).
func BenchmarkE2EKeyGen(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			b.Fatalf("rsa.GenerateKey: %v", err)
		}
		_ = key
	}
}

// BenchmarkE2EEncryptValue measures one full hybrid encryption producing the
// 7-chunk wire-format value (32-byte AES key + 12-byte IV + RSA-OAEP wrap +
// AES-GCM seal).
func BenchmarkE2EEncryptValue(b *testing.B) {
	v := loadVector(b)
	device := vectorDevice(b, v, v.KeyVersion)
	encryptor := smsgateway.NewEncryptor()
	plaintext := "The quick brown fox jumps over the lazy dog 0123456789"

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		value, err := encryptor.Encrypt(device, plaintext)
		if err != nil {
			b.Fatalf("Encrypt: %v", err)
		}
		_ = value
	}
}

// BenchmarkE2EDecryptValue measures one full hybrid decryption mirroring the
// Android decryptE2E path: RSA-OAEP-SHA256 decrypt of chunk 4, then AES-GCM
// open of chunk 6.
func BenchmarkE2EDecryptValue(b *testing.B) {
	v := loadVector(b)
	device := vectorDevice(b, v, v.KeyVersion)
	priv := parseVectorPrivateKey(b, v)
	encryptor := smsgateway.NewEncryptor()
	plaintext := "The quick brown fox jumps over the lazy dog 0123456789"
	value, err := encryptor.Encrypt(device, plaintext)
	if err != nil {
		b.Fatalf("Encrypt: %v", err)
	}

	chunks := strings.Split(value, "$")
	encKey, err := base64.StdEncoding.DecodeString(chunks[4])
	if err != nil {
		b.Fatalf("decode chunk 4: %v", err)
	}
	iv, err := base64.StdEncoding.DecodeString(chunks[5])
	if err != nil {
		b.Fatalf("decode chunk 5: %v", err)
	}
	ctTag, err := base64.StdEncoding.DecodeString(chunks[6])
	if err != nil {
		b.Fatalf("decode chunk 6: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		decryptE2EOne(b, priv, encKey, iv, ctTag)
	}
}

// decryptE2EOne performs one full hybrid decryption mirroring the Android
// decryptE2E path (RSA-OAEP-SHA256 decrypt of chunk 4, then AES-GCM open of
// chunk 6).
func decryptE2EOne(b *testing.B, priv *rsa.PrivateKey, encKey, iv, ctTag []byte) {
	b.Helper()
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)
	if err != nil {
		b.Fatalf("RSA-OAEP decrypt: %v", err)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		b.Fatalf("aes.NewCipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		b.Fatalf("cipher.NewGCM: %v", err)
	}
	plain, err := gcm.Open(nil, iv, ctTag, nil)
	if err != nil {
		b.Fatalf("AES-GCM open: %v", err)
	}
	_ = plain
}

// BenchmarkE2EDecryptBatch100 measures 100 sequential message decryptions -
// the acceptance scenario "100 messages decrypted in < 15s" (Android
// decryptE2E per message; host baseline only).
func BenchmarkE2EDecryptBatch100(b *testing.B) {
	v := loadVector(b)
	device := vectorDevice(b, v, v.KeyVersion)
	priv := parseVectorPrivateKey(b, v)
	encryptor := smsgateway.NewEncryptor()
	plaintext := "The quick brown fox jumps over the lazy dog 0123456789"

	type sealed struct {
		encKey []byte
		iv     []byte
		ctTag  []byte
	}
	values := make([]sealed, 100)
	for i := range values {
		value, err := encryptor.Encrypt(device, plaintext)
		if err != nil {
			b.Fatalf("Encrypt: %v", err)
		}
		chunks := strings.Split(value, "$")
		encKey, err := base64.StdEncoding.DecodeString(chunks[4])
		if err != nil {
			b.Fatalf("decode chunk 4: %v", err)
		}
		iv, err := base64.StdEncoding.DecodeString(chunks[5])
		if err != nil {
			b.Fatalf("decode chunk 5: %v", err)
		}
		ctTag, err := base64.StdEncoding.DecodeString(chunks[6])
		if err != nil {
			b.Fatalf("decode chunk 6: %v", err)
		}
		values[i] = sealed{encKey: encKey, iv: iv, ctTag: ctTag}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		for _, s := range values {
			decryptE2EOne(b, priv, s.encKey, s.iv, s.ctTag)
		}
	}
}

// loadVector/vectorDevice/parseVectorPrivateKey accept *testing.B by sharing
// the *testing.T helpers through the vectorTestHelper interface (defined in
// encryption_test.go).
