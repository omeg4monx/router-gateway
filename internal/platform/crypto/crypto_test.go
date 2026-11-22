package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLoadOrCreateMasterKey(t *testing.T) {
	t.Run("creates master.key with 0600 permissions and reloads identical 32 bytes", func(t *testing.T) {
		dataDir := filepath.Join(t.TempDir(), "vault-data")

		key1, err := LoadOrCreateMasterKey(dataDir, "")
		if err != nil {
			t.Fatalf("first LoadOrCreateMasterKey failed: %v", err)
		}
		if len(key1) != MasterKeySize {
			t.Fatalf("got key length %d, want %d", len(key1), MasterKeySize)
		}

		keyPath := filepath.Join(dataDir, "master.key")
		info, err := os.Stat(keyPath)
		if err != nil {
			t.Fatalf("stat master.key failed: %v", err)
		}
		if perm := info.Mode().Perm(); perm != 0o600 {
			t.Fatalf("master.key permissions = %04o, want 0600", perm)
		}

		key2, err := LoadOrCreateMasterKey(dataDir, "")
		if err != nil {
			t.Fatalf("second LoadOrCreateMasterKey failed: %v", err)
		}
		if !bytes.Equal(key1, key2) {
			t.Fatal("reloaded master key does not match created key")
		}
	})

	t.Run("valid 64-char hex env override takes precedence without creating file", func(t *testing.T) {
		dataDir := filepath.Join(t.TempDir(), "unused-dir")
		raw := bytes.Repeat([]byte{0xab}, MasterKeySize)
		envHex := hex.EncodeToString(raw)

		key, err := LoadOrCreateMasterKey(dataDir, envHex)
		if err != nil {
			t.Fatalf("LoadOrCreateMasterKey with env override failed: %v", err)
		}
		if !bytes.Equal(key, raw) {
			t.Fatalf("got key %x, want %x", key, raw)
		}

		if _, err := os.Stat(filepath.Join(dataDir, "master.key")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected master.key not to be created when env override is set, stat err = %v", err)
		}
	})

	t.Run("malformed hex or wrong length env override returns error", func(t *testing.T) {
		cases := []struct {
			name   string
			envHex string
		}{
			{name: "malformed hex", envHex: strings.Repeat("zz", MasterKeySize)},
			{name: "too short", envHex: hex.EncodeToString(make([]byte, 16))},
			{name: "too long", envHex: hex.EncodeToString(make([]byte, 33))},
			{name: "odd length hex", envHex: "abc"},
			{name: "whitespace around valid hex", envHex: " " + hex.EncodeToString(make([]byte, MasterKeySize)) + " "},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := LoadOrCreateMasterKey(t.TempDir(), tc.envHex)
				if err == nil {
					t.Fatalf("expected error for %s, got nil", tc.name)
				}
				if !errors.Is(err, ErrInvalidMasterKey) {
					t.Fatalf("got error %v, want errors.Is(_, ErrInvalidMasterKey)", err)
				}
			})
		}
	})

	t.Run("existing master.key with wrong length returns error", func(t *testing.T) {
		dataDir := t.TempDir()
		keyPath := filepath.Join(dataDir, "master.key")
		if err := os.WriteFile(keyPath, []byte("short-key"), 0o600); err != nil {
			t.Fatalf("write invalid master.key: %v", err)
		}

		_, err := LoadOrCreateMasterKey(dataDir, "")
		if err == nil {
			t.Fatal("expected error for truncated master.key, got nil")
		}
		if !errors.Is(err, ErrInvalidMasterKey) {
			t.Fatalf("got error %v, want errors.Is(_, ErrInvalidMasterKey)", err)
		}
	})

	t.Run("empty dataDir without env override returns error", func(t *testing.T) {
		_, err := LoadOrCreateMasterKey("", "")
		if err == nil {
			t.Fatal("expected error for empty dataDir, got nil")
		}
		if !errors.Is(err, ErrInvalidDataDir) {
			t.Fatalf("got error %v, want errors.Is(_, ErrInvalidDataDir)", err)
		}
	})
}

func TestVaultRoundTripAndTamperRejection(t *testing.T) {
	keyA := bytes.Repeat([]byte{0x11}, MasterKeySize)
	keyB := bytes.Repeat([]byte{0x22}, MasterKeySize)

	t.Run("NewVault validates key length", func(t *testing.T) {
		for _, length := range []int{0, 16, 31, 33, 64} {
			if _, err := NewVault(make([]byte, length)); !errors.Is(err, ErrInvalidMasterKey) {
				t.Fatalf("NewVault(len=%d) err = %v, want ErrInvalidMasterKey", length, err)
			}
		}
	})

	vaultA, err := NewVault(keyA)
	if err != nil {
		t.Fatalf("NewVault(keyA) failed: %v", err)
	}
	vaultB, err := NewVault(keyB)
	if err != nil {
		t.Fatalf("NewVault(keyB) failed: %v", err)
	}

	t.Run("distinct nonces and clean round-trip", func(t *testing.T) {
		plaintext := []byte("sk-secret-provider-api-key-123456")
		ct1, err := vaultA.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt #1 failed: %v", err)
		}
		ct2, err := vaultA.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encrypt #2 failed: %v", err)
		}

		if !strings.HasPrefix(ct1, "v1:") || !strings.HasPrefix(ct2, "v1:") {
			t.Fatalf("expected v1: prefix, got %q and %q", ct1, ct2)
		}
		if ct1 == ct2 {
			t.Fatal("encrypting identical plaintext twice produced identical ciphertexts; expected random nonce")
		}

		got1, err := vaultA.Decrypt(ct1)
		if err != nil {
			t.Fatalf("Decrypt #1 failed: %v", err)
		}
		if !bytes.Equal(got1, plaintext) {
			t.Fatalf("Decrypt #1 = %q, want %q", got1, plaintext)
		}

		got2, err := vaultA.Decrypt(ct2)
		if err != nil {
			t.Fatalf("Decrypt #2 failed: %v", err)
		}
		if !bytes.Equal(got2, plaintext) {
			t.Fatalf("Decrypt #2 = %q, want %q", got2, plaintext)
		}

		emptyCT, err := vaultA.Encrypt([]byte{})
		if err != nil {
			t.Fatalf("Encrypt empty failed: %v", err)
		}
		emptyPT, err := vaultA.Decrypt(emptyCT)
		if err != nil {
			t.Fatalf("Decrypt empty failed: %v", err)
		}
		if len(emptyPT) != 0 {
			t.Fatalf("Decrypt empty = %x, want empty slice", emptyPT)
		}
	})

	t.Run("bit-flipped ciphertext and wrong key fail", func(t *testing.T) {
		ct, err := vaultA.Encrypt([]byte("sensitive-secret"))
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		// Flip a bit in the decoded ciphertext body.
		raw, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(ct, "v1:"))
		if err != nil {
			t.Fatalf("decode ciphertext: %v", err)
		}
		raw[len(raw)-1] ^= 0x01
		tampered := "v1:" + base64.RawURLEncoding.EncodeToString(raw)

		if _, err := vaultA.Decrypt(tampered); !errors.Is(err, ErrDecryptionFailed) {
			t.Fatalf("Decrypt(tampered) err = %v, want ErrDecryptionFailed", err)
		}

		if _, err := vaultB.Decrypt(ct); !errors.Is(err, ErrDecryptionFailed) {
			t.Fatalf("Decrypt(wrong key) err = %v, want ErrDecryptionFailed", err)
		}
	})

	t.Run("malformed and truncated ciphertexts fail without panic", func(t *testing.T) {
		malformed := []string{
			"",
			"v1:",
			"v2:AAAA",
			"no-prefix",
			"v1:!!!not-base64!!!",
			"v1:" + base64.RawURLEncoding.EncodeToString(make([]byte, 10)),
			"v1:" + base64.RawURLEncoding.EncodeToString(make([]byte, 27)),
		}
		for _, input := range malformed {
			if _, err := vaultA.Decrypt(input); !errors.Is(err, ErrInvalidCiphertext) {
				t.Fatalf("Decrypt(%q) err = %v, want ErrInvalidCiphertext", input, err)
			}
		}
	})

	t.Run("concurrent encrypt and decrypt", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := range 16 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				msg := []byte("secret-payload-" + string(rune('a'+id)))
				ct, err := vaultA.Encrypt(msg)
				if err != nil {
					t.Errorf("goroutine %d Encrypt failed: %v", id, err)
					return
				}
				pt, err := vaultA.Decrypt(ct)
				if err != nil {
					t.Errorf("goroutine %d Decrypt failed: %v", id, err)
					return
				}
				if !bytes.Equal(pt, msg) {
					t.Errorf("goroutine %d round-trip mismatch", id)
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestPasswordHashAndVerify(t *testing.T) {
	t.Run("rejects empty password", func(t *testing.T) {
		if _, err := HashPassword(""); !errors.Is(err, ErrEmptyPassword) {
			t.Fatalf("HashPassword(\"\") err = %v, want ErrEmptyPassword", err)
		}
	})

	t.Run("rejects password longer than 72 bytes", func(t *testing.T) {
		tooLong := strings.Repeat("a", 73)
		if _, err := HashPassword(tooLong); !errors.Is(err, ErrPasswordTooLong) {
			t.Fatalf("HashPassword(73 bytes) err = %v, want ErrPasswordTooLong", err)
		}
	})

	t.Run("hashes and verifies matching and non-matching passwords", func(t *testing.T) {
		pw := "correct-horse-battery-staple"
		h1, err := HashPassword(pw)
		if err != nil {
			t.Fatalf("HashPassword #1 failed: %v", err)
		}
		h2, err := HashPassword(pw)
		if err != nil {
			t.Fatalf("HashPassword #2 failed: %v", err)
		}
		if h1 == h2 {
			t.Fatal("expected distinct bcrypt hashes due to random salt")
		}
		if !VerifyPassword(h1, pw) {
			t.Fatal("VerifyPassword returned false for matching password")
		}
		if VerifyPassword(h1, "wrong-password") {
			t.Fatal("VerifyPassword returned true for wrong password")
		}
		if VerifyPassword(h1, "") {
			t.Fatal("VerifyPassword returned true for empty password")
		}
		if VerifyPassword("", pw) {
			t.Fatal("VerifyPassword returned true for empty hash")
		}
		if VerifyPassword("not-a-bcrypt-hash", pw) {
			t.Fatal("VerifyPassword returned true for malformed hash")
		}
	})
}

func TestSessionSignVerifyAndExpiry(t *testing.T) {
	keyA := bytes.Repeat([]byte{0x33}, MasterKeySize)
	keyB := bytes.Repeat([]byte{0x44}, MasterKeySize)

	t.Run("NewSessionSigner validates key length", func(t *testing.T) {
		for _, length := range []int{0, 16, 31, 33} {
			if _, err := NewSessionSigner(make([]byte, length)); !errors.Is(err, ErrInvalidMasterKey) {
				t.Fatalf("NewSessionSigner(len=%d) err = %v, want ErrInvalidMasterKey", length, err)
			}
		}
	})

	signerA, err := NewSessionSigner(keyA)
	if err != nil {
		t.Fatalf("NewSessionSigner(keyA) failed: %v", err)
	}
	signerB, err := NewSessionSigner(keyB)
	if err != nil {
		t.Fatalf("NewSessionSigner(keyB) failed: %v", err)
	}

	t.Run("sign and verify valid session token", func(t *testing.T) {
		before := time.Now().UTC()
		token, err := signerA.Sign("admin", time.Hour)
		if err != nil {
			t.Fatalf("Sign failed: %v", err)
		}

		claims, err := signerA.Verify(token)
		if err != nil {
			t.Fatalf("Verify failed: %v", err)
		}
		if claims.Subject != "admin" {
			t.Fatalf("claims.Subject = %q, want %q", claims.Subject, "admin")
		}
		if claims.ExpiresAt.Before(before.Add(59 * time.Minute)) {
			t.Fatalf("claims.ExpiresAt = %v, expected roughly +1h from %v", claims.ExpiresAt, before)
		}
	})

	t.Run("rejects empty subject", func(t *testing.T) {
		if _, err := signerA.Sign("", time.Hour); !errors.Is(err, ErrEmptySubject) {
			t.Fatalf("Sign(\"\") err = %v, want ErrEmptySubject", err)
		}
	})

	t.Run("tampered signature, tampered payload, and wrong key fail", func(t *testing.T) {
		token, err := signerA.Sign("admin", time.Hour)
		if err != nil {
			t.Fatalf("Sign failed: %v", err)
		}
		parts := strings.Split(token, ".")
		if len(parts) != 2 {
			t.Fatalf("expected 2 dot-separated token parts, got %d", len(parts))
		}

		// Tamper signature byte
		sigRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			t.Fatalf("decode signature: %v", err)
		}
		sigRaw[0] ^= 0xff
		tamperedSigToken := parts[0] + "." + base64.RawURLEncoding.EncodeToString(sigRaw)
		if _, err := signerA.Verify(tamperedSigToken); !errors.Is(err, ErrInvalidSignature) {
			t.Fatalf("Verify(tampered signature) err = %v, want ErrInvalidSignature", err)
		}

		// Tamper payload
		payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
		if err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		tamperedPayload := bytes.Replace(payloadRaw, []byte("admin"), []byte("root!"), 1)
		tamperedPayloadToken := base64.RawURLEncoding.EncodeToString(tamperedPayload) + "." + parts[1]
		if _, err := signerA.Verify(tamperedPayloadToken); !errors.Is(err, ErrInvalidSignature) {
			t.Fatalf("Verify(tampered payload) err = %v, want ErrInvalidSignature", err)
		}

		// Wrong signer key
		if _, err := signerB.Verify(token); !errors.Is(err, ErrInvalidSignature) {
			t.Fatalf("Verify(wrong key) err = %v, want ErrInvalidSignature", err)
		}
	})

	t.Run("negative TTL produces expired token error", func(t *testing.T) {
		token, err := signerA.Sign("admin", -time.Second)
		if err != nil {
			t.Fatalf("Sign(negative TTL) failed: %v", err)
		}
		if _, err := signerA.Verify(token); !errors.Is(err, ErrTokenExpired) {
			t.Fatalf("Verify(expired token) err = %v, want ErrTokenExpired", err)
		}
	})

	t.Run("malformed tokens fail with ErrInvalidToken", func(t *testing.T) {
		malformed := []string{
			"",
			"no-dot-at-all",
			".only-sig",
			"only-payload.",
			"a.b.c",
			"!!!.!!!",
			"abc.short",
		}
		for _, tok := range malformed {
			if _, err := signerA.Verify(tok); !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("Verify(%q) err = %v, want ErrInvalidToken", tok, err)
			}
		}
	})

	t.Run("concurrent sign and verify", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := range 16 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				sub := "user-" + string(rune('a'+id))
				tok, err := signerA.Sign(sub, time.Minute)
				if err != nil {
					t.Errorf("goroutine %d Sign failed: %v", id, err)
					return
				}
				claims, err := signerA.Verify(tok)
				if err != nil {
					t.Errorf("goroutine %d Verify failed: %v", id, err)
					return
				}
				if claims.Subject != sub {
					t.Errorf("goroutine %d subject = %q, want %q", id, claims.Subject, sub)
				}
			}(i)
		}
		wg.Wait()
	})
}
