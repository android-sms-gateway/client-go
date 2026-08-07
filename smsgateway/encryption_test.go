package smsgateway_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const vectorPath = "../../docs/plan/e2e-encryption/test-vectors/e2e-vector-v1.json"

// vectorTestHelper is the subset of [testing.TB] used by the vector helpers so
// both test and benchmark cases can use them (benchmarks in
// encryption_benchmark_test.go).
type vectorTestHelper interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Skipf(format string, args ...any)
}

// vector is the JSON schema of e2e-vector-v1.json (subset of the generator's Vector).
type vector struct {
	SchemaVersion            int    `json:"schemaVersion"`
	Format                   string `json:"format"`
	Version                  string `json:"version"`
	KeyVersion               int    `json:"keyVersion"`
	PrivateKeyPEM            string `json:"privateKeyPem"`
	PublicKeyPEM             string `json:"publicKeyPem"`
	PublicKeySpkiBase64      string `json:"publicKeySpkiBase64"`
	AesKeyHex                string `json:"aesKeyHex"`
	IvHex                    string `json:"ivHex"`
	Plaintext                string `json:"plaintext"`
	Chunk4EncryptedAesKeyB64 string `json:"chunk4EncryptedAesKeyB64"`
	CtTagChunkB64            string `json:"ctTagChunkB64"`
	FullFormatSample         string `json:"fullFormatSample"`
}

func loadVector(t vectorTestHelper) *vector {
	t.Helper()
	data, err := os.ReadFile(vectorPath)
	if err != nil {
		t.Skipf("test vector not available: %v", err)
	}

	v := new(vector)
	err = json.Unmarshal(data, v)
	if err != nil {
		t.Fatalf("unmarshal vector: %v", err)
	}

	return v
}

// vectorDevice builds a Device carrying the vector's public key as the E2E key
// material, mirroring what the client passes to [smsgateway.Encryptor.Encrypt].
func vectorDevice(t vectorTestHelper, v *vector, keyVersion int) smsgateway.Device {
	t.Helper()
	kv := keyVersion
	return smsgateway.Device{
		VersionedPublicKey: smsgateway.VersionedPublicKey{
			PublicKey:  &v.PublicKeySpkiBase64,
			KeyVersion: &kv,
		},
	}
}

func parseVectorPrivateKey(t vectorTestHelper, v *vector) *rsa.PrivateKey {
	t.Helper()
	block, _ := pem.Decode([]byte(v.PrivateKeyPEM))
	if block == nil {
		t.Fatalf("failed to parse private key PEM")
		return nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		t.Fatalf("private key is not RSA: %T", key)
	}
	return priv
}

// seedReader is an infinite deterministic byte stream: SHA-256(seed || LE64(counter)).
// Mirrors the generator's seedReader so SDK output is reproducible in-process.
type seedReader struct {
	seed    []byte
	counter uint64
	buf     [32]byte
	pos     int
}

func newSeedReader(seed string) *seedReader {
	return &seedReader{seed: []byte(seed), pos: len(seedReader{}.buf)}
}

func (r *seedReader) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		if r.pos == len(r.buf) {
			h := sha256.New()
			h.Write(r.seed)
			var c [8]byte
			binary.LittleEndian.PutUint64(c[:], r.counter)
			h.Write(c[:])
			copy(r.buf[:], h.Sum(nil))
			r.counter++
			r.pos = 0
		}
		m := copy(p[n:], r.buf[r.pos:])
		r.pos += m
		n += m
	}
	return n, nil
}

// Happy path: an Encryptor pinned to the vector's AES key + IV must produce a
// full 7-chunk value whose chunk 6 byte-equals the committed ctTagChunkB64 and
// whose chunk 4 RSA-decrypts to the fixed 32-byte AES key.
func TestEncryptorVectorFullFormat(t *testing.T) {
	v := loadVector(t)
	device := vectorDevice(t, v, v.KeyVersion)
	priv := parseVectorPrivateKey(t, v)

	aesKey, err := hex.DecodeString(v.AesKeyHex)
	if err != nil {
		t.Fatalf("decode aesKeyHex: %v", err)
	}
	iv, err := hex.DecodeString(v.IvHex)
	if err != nil {
		t.Fatalf("decode ivHex: %v", err)
	}
	if len(aesKey) != 32 {
		t.Fatalf("AES key length = %d, want 32 (12 bytes is the IV length)", len(aesKey))
	}
	if len(iv) != 12 {
		t.Fatalf("IV length = %d, want 12", len(iv))
	}

	e := smsgateway.NewEncryptor(
		smsgateway.WithRandom(newSeedReader("vector")),
		smsgateway.WithEncryptionMaterial(aesKey, iv),
	)
	value, err := e.Encrypt(device, v.Plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	chunks := strings.Split(value, "$")
	wantChunks := []string{
		"",
		"rsa-oaep-aes-256-gcm",
		"v=1",
		"k=1",
	}
	if len(chunks) != 7 {
		t.Fatalf("Encrypt produced %d chunks, want exactly 7", len(chunks))
	}
	for i, want := range wantChunks {
		if chunks[i] != want {
			t.Fatalf("chunk %d = %q, want %q", i, chunks[i], want)
		}
	}
	if chunks[5] != base64.StdEncoding.EncodeToString(iv) {
		t.Fatalf("chunk 5 = %q, want base64 of vector IV %q", chunks[5], base64.StdEncoding.EncodeToString(iv))
	}
	if chunks[6] != v.CtTagChunkB64 {
		t.Fatalf("chunk 6 = %q, want committed ctTagChunkB64 %q", chunks[6], v.CtTagChunkB64)
	}

	// Chunk 4 must RSA-OAEP-decrypt (SHA-256/MGF1-SHA-256, empty label) to the
	// fixed 32-byte AES key.
	encAesKey, err := base64.StdEncoding.DecodeString(chunks[4])
	if err != nil {
		t.Fatalf("decode chunk 4: %v", err)
	}
	if len(encAesKey) != 256 {
		t.Fatalf("chunk 4 length = %d, want 256 (RSA-2048)", len(encAesKey))
	}
	gotKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encAesKey, nil)
	if err != nil {
		t.Fatalf("RSA-OAEP decrypt chunk 4: %v", err)
	}
	if !bytesEqual(gotKey, aesKey) {
		t.Fatalf("chunk 4 decrypts to %x, want %x", gotKey, aesKey)
	}

	// Chunk 6 must AES-GCM-decrypt (empty AAD, 128-bit tag) to the plaintext.
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
	gotPlain, err := gcm.Open(nil, iv, ctTag, nil)
	if err != nil {
		t.Fatalf("AES-GCM decrypt chunk 6: %v", err)
	}
	if string(gotPlain) != v.Plaintext {
		t.Fatalf("chunk 6 decrypts to %q, want %q", gotPlain, v.Plaintext)
	}
}

// Invariant: the committed vector's chunk 4 must RSA-OAEP-decrypt to the fixed
// 32-byte AES key (verifies the committed artifact, not just our encryptor).
func TestVectorChunk4DecryptsToAesKey(t *testing.T) {
	v := loadVector(t)
	priv := parseVectorPrivateKey(t, v)
	aesKey, err := hex.DecodeString(v.AesKeyHex)
	if err != nil {
		t.Fatalf("decode aesKeyHex: %v", err)
	}

	encKey, err := base64.StdEncoding.DecodeString(v.Chunk4EncryptedAesKeyB64)
	if err != nil {
		t.Fatalf("decode chunk4EncryptedAesKeyB64: %v", err)
	}
	got, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)
	if err != nil {
		t.Fatalf("RSA-OAEP decrypt vector chunk4: %v", err)
	}
	if !bytesEqual(got, aesKey) {
		t.Fatalf("vector chunk4 decrypts to %x, want %x", got, aesKey)
	}
}

// Invariant: the committed vector's full format sample is exactly 7 chunks.
func TestVectorFullFormatChunks(t *testing.T) {
	v := loadVector(t)
	iv, err := hex.DecodeString(v.IvHex)
	if err != nil {
		t.Fatalf("decode ivHex: %v", err)
	}

	chunks := strings.Split(v.FullFormatSample, "$")
	want := []string{
		"",
		"rsa-oaep-aes-256-gcm",
		"v=1",
		"k=1",
		v.Chunk4EncryptedAesKeyB64,
		base64.StdEncoding.EncodeToString(iv),
		v.CtTagChunkB64,
	}
	if len(chunks) != len(want) {
		t.Fatalf("fullFormatSample splits into %d chunks, want %d", len(chunks), len(want))
	}
	for i := range want {
		if chunks[i] != want[i] {
			t.Fatalf("fullFormatSample chunk %d = %q, want %q", i, chunks[i], want[i])
		}
	}
}

// Byte determinism within one process run: the same injected rand stream must
// reproduce a byte-identical full value; a different stream must produce a
// different RSA chunk 4 but the same plaintext after decryption.
func TestEncryptorByteDeterminism(t *testing.T) {
	v := loadVector(t)
	device := vectorDevice(t, v, v.KeyVersion)
	priv := parseVectorPrivateKey(t, v)

	encrypt := func(seed string) string {
		t.Helper()
		e := smsgateway.NewEncryptor(smsgateway.WithRandom(newSeedReader(seed)))
		value, err := e.Encrypt(device, v.Plaintext)
		if err != nil {
			t.Fatalf("Encrypt(%q): %v", seed, err)
		}
		return value
	}

	a := encrypt("same-seed")
	b := encrypt("same-seed")
	if a != b {
		t.Fatal("same rand stream must produce byte-identical encrypted value")
	}

	c := encrypt("other-seed")
	if c == a {
		t.Fatal("different rand stream must produce a different encrypted value")
	}

	// Both values must be decryptable to the plaintext by the private key.
	for _, value := range []string{a, c} {
		chunks := strings.Split(value, "$")
		encKey, err := base64.StdEncoding.DecodeString(chunks[4])
		if err != nil {
			t.Fatalf("decode chunk 4: %v", err)
		}
		aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encKey, nil)
		if err != nil {
			t.Fatalf("RSA-OAEP decrypt chunk 4: %v", err)
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
		iv, err := base64.StdEncoding.DecodeString(chunks[5])
		if err != nil {
			t.Fatalf("decode chunk 5: %v", err)
		}
		gotPlain, err := gcm.Open(nil, iv, ctTag, nil)
		if err != nil {
			t.Fatalf("AES-GCM decrypt chunk 6: %v", err)
		}
		if string(gotPlain) != v.Plaintext {
			t.Fatalf("decrypted %q, want %q", gotPlain, v.Plaintext)
		}
	}
}

// Fresh IVs: two separately-encrypted values must never share an IV.
func TestEncryptorDistinctIVs(t *testing.T) {
	v := loadVector(t)
	device := vectorDevice(t, v, v.KeyVersion)

	e1 := smsgateway.NewEncryptor(smsgateway.WithRandom(newSeedReader("iv-1")))
	e2 := smsgateway.NewEncryptor(smsgateway.WithRandom(newSeedReader("iv-2")))

	first, err := e1.Encrypt(device, "phone 1")
	if err != nil {
		t.Fatalf("Encrypt 1: %v", err)
	}
	second, err := e2.Encrypt(device, "phone 2")
	if err != nil {
		t.Fatalf("Encrypt 2: %v", err)
	}

	iv1 := strings.Split(first, "$")[5]
	iv2 := strings.Split(second, "$")[5]
	if iv1 == iv2 {
		t.Fatalf("IVs must be fresh per value, got shared IV %q", iv1)
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// MessageEncryptor owns key retrieval: missing or unusable key material on the
// resolved Device must be validated inside Encrypt, before any material is
// generated.
func TestMessageEncryptor_Encrypt_KeyNegative(t *testing.T) {
	v := loadVector(t)
	e := smsgateway.NewEncryptor(smsgateway.WithRandom(newSeedReader("negative-key")))

	// Nil PublicKey => ErrE2ENotConfigured.
	noKey := vectorDevice(t, v, 1)
	noKey.PublicKey = nil
	if _, err := e.Encrypt(noKey, "secret"); !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("nil PublicKey error = %v, want ErrE2ENotConfigured", err)
	}

	// Nil KeyVersion => ErrE2ENotConfigured.
	noVersion := vectorDevice(t, v, 1)
	noVersion.KeyVersion = nil
	if _, err := e.Encrypt(noVersion, "secret"); !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("nil KeyVersion error = %v, want ErrE2ENotConfigured", err)
	}

	// Both PublicKey and KeyVersion nil (zero Device) => ErrE2ENotConfigured.
	if _, err := e.Encrypt(smsgateway.Device{}, "secret"); !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("zero Device error = %v, want ErrE2ENotConfigured", err)
	}

	// Malformed base64 public key => error, but not ErrE2ENotConfigured (the
	// key is present but undecodable).
	malformed := vectorDevice(t, v, 1)
	badB64 := "this is not base64!!!"
	malformed.PublicKey = &badB64
	if _, err := e.Encrypt(malformed, "secret"); err == nil || errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("malformed key error = %v, want a non-ErrE2ENotConfigured error", err)
	}

	// Non-RSA public key (EC) => ErrE2ENotConfigured.
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate EC key: %v", err)
	}
	ecDer, err := x509.MarshalPKIXPublicKey(&ecKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal EC public key: %v", err)
	}
	nonRSA := vectorDevice(t, v, 1)
	nonRSAKey := base64.StdEncoding.EncodeToString(ecDer)
	nonRSA.PublicKey = &nonRSAKey
	if _, err = e.Encrypt(nonRSA, "secret"); !errors.Is(err, smsgateway.ErrE2ENotConfigured) {
		t.Fatalf("non-RSA key error = %v, want ErrE2ENotConfigured", err)
	}
}
