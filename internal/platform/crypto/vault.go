package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	// MasterKeySize is the required byte length for AES-256 master keys.
	MasterKeySize = 32

	masterKeyFileName  = "master.key"
	vaultVersionPrefix = "v1:"
)

var (
	// ErrInvalidMasterKey indicates a missing, malformed, or wrong-length master key.
	ErrInvalidMasterKey = errors.New("crypto: invalid master key")
	// ErrInvalidDataDir indicates an empty or unusable data directory path.
	ErrInvalidDataDir = errors.New("crypto: invalid data directory")
	// ErrInvalidCiphertext indicates malformed, truncated, or unversioned ciphertext.
	ErrInvalidCiphertext = errors.New("crypto: invalid ciphertext")
	// ErrDecryptionFailed indicates authenticated decryption failed (wrong key or tampered data).
	ErrDecryptionFailed = errors.New("crypto: decryption failed")
)

// LoadOrCreateMasterKey resolves the 32-byte master key.
// If envHex is non-empty, it must be an exact 64-character hexadecimal encoding
// of 32 bytes and takes precedence without touching disk.
// Otherwise, it loads <dataDir>/master.key (validating exact 32-byte length) or
// generates a new random 32-byte key and persists it atomically with 0600 permissions.
func LoadOrCreateMasterKey(dataDir, envHex string) ([]byte, error) {
	if envHex != "" {
		if len(envHex) != MasterKeySize*2 {
			return nil, fmt.Errorf("%w: hex key must be %d characters, got %d", ErrInvalidMasterKey, MasterKeySize*2, len(envHex))
		}
		key, err := hex.DecodeString(envHex)
		if err != nil || len(key) != MasterKeySize {
			return nil, fmt.Errorf("%w: invalid hex encoding", ErrInvalidMasterKey)
		}
		return key, nil
	}

	if strings.TrimSpace(dataDir) == "" {
		return nil, ErrInvalidDataDir
	}

	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, fmt.Errorf("crypto: create data directory: %w", err)
	}

	keyPath := filepath.Join(dataDir, masterKeyFileName)
	existing, err := readMasterKeyFile(keyPath)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	key := make([]byte, MasterKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("crypto: generate master key: %w", err)
	}

	if err := writeKeyAtomically(dataDir, keyPath, key); err != nil {
		if errors.Is(err, os.ErrExist) {
			return readMasterKeyFile(keyPath)
		}
		return nil, err
	}
	return key, nil
}

func readMasterKeyFile(keyPath string) ([]byte, error) {
	existing, err := os.ReadFile(keyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, fmt.Errorf("crypto: read master key: %w", err)
	}
	if len(existing) != MasterKeySize {
		return nil, fmt.Errorf("%w: master.key length %d, want %d", ErrInvalidMasterKey, len(existing), MasterKeySize)
	}
	return existing, nil
}

func writeKeyAtomically(dataDir, keyPath string, key []byte) error {
	tmpFile, err := os.CreateTemp(dataDir, ".master.key.*")
	if err != nil {
		return fmt.Errorf("crypto: create temp master key file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if err := os.Chmod(tmpPath, 0o600); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("crypto: chmod temp master key file: %w", err)
	}

	if _, err := tmpFile.Write(key); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("crypto: write temp master key file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("crypto: sync temp master key file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("crypto: close temp master key file: %w", err)
	}

	if err := os.Link(tmpPath, keyPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			return os.ErrExist
		}
		if isLinkUnsupported(err) {
			return writeKeyExclusive(keyPath, key)
		}
		return fmt.Errorf("crypto: link master key file: %w", err)
	}
	return nil
}

func isLinkUnsupported(err error) bool {
	return errors.Is(err, errors.ErrUnsupported) ||
		errors.Is(err, syscall.EPERM) ||
		errors.Is(err, syscall.EOPNOTSUPP) ||
		errors.Is(err, syscall.ENOTSUP)
}

func writeKeyExclusive(keyPath string, key []byte) error {
	f, err := os.OpenFile(keyPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return os.ErrExist
		}
		return fmt.Errorf("crypto: create master key file: %w", err)
	}
	if err := os.Chmod(keyPath, 0o600); err != nil {
		_ = f.Close()
		return fmt.Errorf("crypto: chmod master key file: %w", err)
	}
	if _, err := f.Write(key); err != nil {
		_ = f.Close()
		return fmt.Errorf("crypto: write master key file: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("crypto: sync master key file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("crypto: close master key file: %w", err)
	}
	return nil
}

// Vault provides authenticated encryption at rest using AES-256-GCM.
// A Vault instance is safe for concurrent use by multiple goroutines.
type Vault struct {
	aead cipher.AEAD
}

// NewVault constructs an AES-256-GCM Vault from a 32-byte key.
func NewVault(key []byte) (*Vault, error) {
	if len(key) != MasterKeySize {
		return nil, fmt.Errorf("%w: vault key length %d, want %d", ErrInvalidMasterKey, len(key), MasterKeySize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: init AES cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: init GCM: %w", err)
	}

	return &Vault{aead: aead}, nil
}

// Encrypt seals plaintext with AES-256-GCM using a fresh 12-byte random nonce
// and encodes the result as "v1:<base64url(nonce || ciphertext)>".
func (v *Vault) Encrypt(plaintext []byte) (string, error) {
	if v == nil || v.aead == nil {
		return "", ErrInvalidMasterKey
	}

	nonceSize := v.aead.NonceSize()
	buf := make([]byte, nonceSize, nonceSize+len(plaintext)+v.aead.Overhead())
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", fmt.Errorf("crypto: generate nonce: %w", err)
	}

	sealed := v.aead.Seal(buf, buf, plaintext, nil)
	return vaultVersionPrefix + base64.RawURLEncoding.EncodeToString(sealed), nil
}

// Decrypt verifies and decrypts a "v1:<base64url(nonce || ciphertext)>" string.
func (v *Vault) Decrypt(encoded string) ([]byte, error) {
	if v == nil || v.aead == nil {
		return nil, ErrInvalidMasterKey
	}

	payload, ok := strings.CutPrefix(encoded, vaultVersionPrefix)
	if !ok || payload == "" {
		return nil, ErrInvalidCiphertext
	}

	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid base64url", ErrInvalidCiphertext)
	}

	nonceSize := v.aead.NonceSize()
	if len(raw) < nonceSize+v.aead.Overhead() {
		return nil, fmt.Errorf("%w: truncated payload", ErrInvalidCiphertext)
	}

	plaintext, err := v.aead.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return nil, fmt.Errorf("%w: authentication failed", ErrDecryptionFailed)
	}
	if plaintext == nil {
		return []byte{}, nil
	}
	return plaintext, nil
}
