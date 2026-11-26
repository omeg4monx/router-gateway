package keys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// TokenPrefix is the canonical prefix for all gateway API key secrets.
	TokenPrefix = "sk-om4-"

	tokenRandomBytes = 36 // 36 bytes -> 48 base64url characters (~288 bits entropy)
	tokenSuffixLen   = 48
	tokenTotalLen    = len(TokenPrefix) + tokenSuffixLen
	keyPrefixLen     = 14
	keyIDRandomBytes = 12 // 12 bytes -> 24 hex characters
	keyIDPrefix      = "key_"
)

func generateKeyID() (string, error) {
	buf := make([]byte, keyIDRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("keys: generate key id: %w", err)
	}
	return keyIDPrefix + hex.EncodeToString(buf), nil
}

func generateToken() (raw, hash, prefix string, err error) {
	buf := make([]byte, tokenRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", "", fmt.Errorf("keys: generate key secret: %w", err)
	}
	raw = TokenPrefix + base64.RawURLEncoding.EncodeToString(buf)
	hash = hashToken(raw)
	prefix = raw[:keyPrefixLen]
	return raw, hash, prefix, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func isValidTokenFormat(raw string) bool {
	if len(raw) != tokenTotalLen || !strings.HasPrefix(raw, TokenPrefix) {
		return false
	}
	suffix := raw[len(TokenPrefix):]
	for i := 0; i < len(suffix); i++ {
		c := suffix[i]
		isUpper := c >= 'A' && c <= 'Z'
		isLower := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if !isUpper && !isLower && !isDigit && c != '-' && c != '_' {
			return false
		}
	}
	return true
}
