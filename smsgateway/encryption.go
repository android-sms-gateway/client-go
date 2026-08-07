package smsgateway

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
)

// E2E wire-format constants. The format is defined in
// docs/plan/e2e-encryption/e2e-crypto-spec.md and matches Android's
// RSA/ECB/OAEPWithSHA-256AndMGF1Padding + AES/GCM/NoPadding decryption.
const (
	e2eFormatID      = "rsa-oaep-aes-256-gcm"
	e2eFormatVersion = "1"
	e2ePrefix        = "$" + e2eFormatID + "$"

	// passphrasePrefix is the legacy passphrase format
	// ($aes-256-cbc/pbkdf2-sha1$). Messages already in that format are sent
	// as-is, never re-encrypted with E2E.
	//
	//nolint:gosec // format prefix constant, not a credential
	passphrasePrefix = "$aes-256-cbc/pbkdf2-sha1$"

	aesKeyLen        = 32  // bytes (AES-256)
	ivLen            = 12  // bytes (GCM nonce)
	gcmTagLen        = 16  // bytes (128 bits)
	rsaCiphertextLen = 256 // bytes (RSA-2048 modulus)
)

// Encryptor produces E2E-encrypted values for the resolved target [Device].
// Implementations own key retrieval: the device carries the key material
// (PublicKey + KeyVersion) that the encryptor parses and validates.
//
// The reference implementation is [MessageEncryptor], which emits the hybrid
// RSA-OAEP + AES-256-GCM wire format:
//
//	$rsa-oaep-aes-256-gcm$v=1$k={keyVersion}${b64(encrypted_aes_key)}${b64(iv)}${b64(ciphertext || tag)}
//
// A custom implementation can plug in another encryption method (for example,
// a future passphrase encryptor). A shared Client calls Encrypt concurrently,
// so implementations must be safe for concurrent use.
type Encryptor interface {
	// Encrypt returns the E2E wire-format value for plaintext, keyed by the
	// target device's key material.
	Encrypt(device Device, plaintext string) (string, error)
}

// MessageEncryptor implements [Encryptor] with the hybrid scheme. Every call
// generates a fresh 32-byte AES key and a fresh 12-byte IV, so each encrypted
// value is unique. Production use MUST NOT share a key or IV across values
// (GCM nonce reuse is catastrophic).
type MessageEncryptor struct {
	rand io.Reader

	aesKey []byte // test-only fixed material, empty in production
	iv     []byte // test-only fixed material, empty in production
}

// EncryptorOption configures a [MessageEncryptor].
type EncryptorOption func(*MessageEncryptor)

// WithRandom sets the random source used for AES key/IV generation and RSA-OAEP
// padding. Defaults to [crypto/rand]. Exposed for reproducible test vectors.
func WithRandom(r io.Reader) EncryptorOption {
	return func(e *MessageEncryptor) {
		e.rand = r
	}
}

// WithEncryptionMaterial pins the AES key and IV used by
// [MessageEncryptor.Encrypt].
//
// TEST ONLY - do not use in production. Reusing a key or IV across values
// breaks AES-GCM security (nonce reuse is catastrophic). Production always
// derives fresh material from the CSPRNG.
func WithEncryptionMaterial(aesKey, iv []byte) EncryptorOption {
	return func(e *MessageEncryptor) {
		e.aesKey = aesKey
		e.iv = iv
	}
}

// NewEncryptor creates a [MessageEncryptor] backed by a CSPRNG.
func NewEncryptor(options ...EncryptorOption) *MessageEncryptor {
	e := &MessageEncryptor{
		rand:   rand.Reader,
		aesKey: nil,
		iv:     nil,
	}
	for _, option := range options {
		option(e)
	}

	return e
}

// Compile-time conformance assertion: MessageEncryptor implements Encryptor.
// (Blank-identifier declaration; gochecknoglobals permits it.)
var _ Encryptor = (*MessageEncryptor)(nil)

// Encrypt produces the 7-chunk E2E wire-format string for the given value.
//
// Key retrieval is owned here: the device's PublicKey (base64 NO_WRAP of the
// X.509 SPKI DER encoding) is imported and its KeyVersion used verbatim; the
// device is the source of truth for the key version. A nil PublicKey or
// KeyVersion, or a non-RSA key, returns ErrE2ENotConfigured.
func (e *MessageEncryptor) Encrypt(device Device, plaintext string) (string, error) {
	if device.PublicKey == nil || device.KeyVersion == nil {
		return "", fmt.Errorf("%w: target device has no public key or key version", ErrE2ENotConfigured)
	}

	pub, err := parsePublicKey(*device.PublicKey)
	if err != nil {
		return "", err
	}
	keyVersion := *device.KeyVersion

	aesKey, iv, err := e.material()
	if err != nil {
		return "", err
	}

	// Chunk 4: RSA-OAEP with SHA-256 (hash and MGF1) and an empty label over
	// the 32-byte AES key. Randomization per RFC 8017 7.1.2 makes chunk 4
	// unpredictable; the device decrypts it with its private key.
	encAesKey, err := rsa.EncryptOAEP(sha256.New(), e.rand, pub, aesKey, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt AES key: %w", err)
	}

	// Chunk 6: AES-256-GCM with a 128-bit tag and empty AAD. gcm.Seal returns
	// ciphertext followed by the 16-byte tag, exactly what Android expects.
	ctTag, err := seal(aesKey, iv, plaintext)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("$%s$v=%s$k=%d$%s$%s$%s",
		e2eFormatID,
		e2eFormatVersion,
		keyVersion,
		base64.StdEncoding.EncodeToString(encAesKey),
		base64.StdEncoding.EncodeToString(iv),
		ctTag,
	), nil
}

// material returns the AES key and IV for this encryption. Test-pinned
// material wins; otherwise fresh CSPRNG bytes are generated (32-byte key,
// 12-byte IV - never 12 bytes for the key).
func (e *MessageEncryptor) material() ([]byte, []byte, error) {
	if len(e.aesKey) == aesKeyLen && len(e.iv) == ivLen {
		return e.aesKey, e.iv, nil
	}

	aesKey := make([]byte, aesKeyLen)
	if _, err := io.ReadFull(e.rand, aesKey); err != nil {
		return nil, nil, fmt.Errorf("failed to read AES key: %w", err)
	}

	iv := make([]byte, ivLen)
	if _, err := io.ReadFull(e.rand, iv); err != nil {
		return nil, nil, fmt.Errorf("failed to read IV: %w", err)
	}

	return aesKey, iv, nil
}

// seal encrypts plaintext with AES-256-GCM (128-bit tag, empty AAD) and
// returns base64(ciphertext || 16-byte tag) for chunk 6.
func seal(aesKey, iv []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	ctTag := gcm.Seal(nil, iv, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ctTag), nil
}

// parsePublicKey imports a device listing publicKey (base64 NO_WRAP of the
// X.509 SPKI DER encoding) as an RSA public key.
func parsePublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%w: public key is not an RSA key", ErrE2ENotConfigured)
	}

	return rsaPub, nil
}

// hasEncryptedFormatPrefix reports whether value carries an encrypted-format
// prefix (E2E or passphrase). Prefix comparison is constant-time as
// recommended by the E2E spec (section 10).
func hasEncryptedFormatPrefix(value string) bool {
	return hasPrefixConstantTime(value, e2ePrefix) || hasPrefixConstantTime(value, passphrasePrefix)
}

func hasPrefixConstantTime(value, prefix string) bool {
	if len(value) < len(prefix) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(value[:len(prefix)]), []byte(prefix)) == 1
}
