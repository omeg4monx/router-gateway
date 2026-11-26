package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const sessionKeyDomain = "session-key"

var (
	// ErrEmptySubject indicates a session token was requested without a subject.
	ErrEmptySubject = errors.New("crypto: session subject must not be empty")
	// ErrInvalidToken indicates a malformed session token.
	ErrInvalidToken = errors.New("crypto: invalid session token")
	// ErrInvalidSignature indicates the session token MAC did not match.
	ErrInvalidSignature = errors.New("crypto: invalid session signature")
	// ErrTokenExpired indicates a verified session token is past its ExpiresAt timestamp.
	ErrTokenExpired = errors.New("crypto: session token expired")
)

// Claims represents the authenticated payload carried by a signed session token.
type Claims struct {
	Subject   string    `json:"sub"`
	ExpiresAt time.Time `json:"exp"`
}

// SessionSigner signs and verifies stateless session tokens using a subkey
// derived from the 32-byte master key via HMAC-SHA256.
// A SessionSigner instance is safe for concurrent use by multiple goroutines.
type SessionSigner struct {
	signingKey []byte
}

// NewSessionSigner creates a SessionSigner with a domain-separated subkey
// derived from key via HMAC-SHA256.
func NewSessionSigner(key []byte) (*SessionSigner, error) {
	if len(key) != MasterKeySize {
		return nil, fmt.Errorf("%w: session signer key length %d, want %d", ErrInvalidMasterKey, len(key), MasterKeySize)
	}

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(sessionKeyDomain))
	subkey := mac.Sum(nil)

	return &SessionSigner{signingKey: subkey}, nil
}

// Sign encodes and signs Claims for subject expiring after ttl.
// Format: base64url(json(Claims)).base64url(hmac)
func (s *SessionSigner) Sign(subject string, ttl time.Duration) (string, error) {
	if s == nil || len(s.signingKey) == 0 {
		return "", ErrInvalidMasterKey
	}
	if subject == "" {
		return "", ErrEmptySubject
	}

	claims := Claims{
		Subject:   subject,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("crypto: marshal session claims: %w", err)
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)
	sig := s.computeMAC(payloadEncoded)
	sigEncoded := base64.RawURLEncoding.EncodeToString(sig)

	return payloadEncoded + "." + sigEncoded, nil
}

// Verify authenticates token using constant-time MAC comparison and validates
// that the token has not expired.
func (s *SessionSigner) Verify(token string) (Claims, error) {
	if s == nil || len(s.signingKey) == 0 {
		return Claims{}, ErrInvalidMasterKey
	}

	payloadEncoded, sigEncoded, ok := strings.Cut(token, ".")
	if !ok || payloadEncoded == "" || sigEncoded == "" || strings.Contains(sigEncoded, ".") {
		return Claims{}, ErrInvalidToken
	}

	sig, err := base64.RawURLEncoding.Strict().DecodeString(sigEncoded)
	if err != nil || len(sig) != sha256.Size {
		return Claims{}, ErrInvalidToken
	}

	expected := s.computeMAC(payloadEncoded)
	if !hmac.Equal(sig, expected) {
		return Claims{}, ErrInvalidSignature
	}

	payloadJSON, err := base64.RawURLEncoding.Strict().DecodeString(payloadEncoded)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return Claims{}, fmt.Errorf("%w: unmarshal claims", ErrInvalidToken)
	}
	if claims.Subject == "" || claims.ExpiresAt.IsZero() {
		return Claims{}, ErrInvalidToken
	}
	if !time.Now().UTC().Before(claims.ExpiresAt) {
		return Claims{}, ErrTokenExpired
	}

	return claims, nil
}

func (s *SessionSigner) computeMAC(payloadEncoded string) []byte {
	mac := hmac.New(sha256.New, s.signingKey)
	_, _ = mac.Write([]byte(payloadEncoded))
	return mac.Sum(nil)
}
