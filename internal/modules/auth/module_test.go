package auth_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/modules/auth"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

type fakeCredentialStore struct {
	mu           sync.RWMutex
	passwordHash string
	requireLogin bool
	hashErr      error
	updateErr    error
	reqErr       error
	setReqErr    error
}

func newFakeStore(t *testing.T, initialPassword string, requireLogin bool) *fakeCredentialStore {
	t.Helper()
	store := &fakeCredentialStore{
		requireLogin: requireLogin,
	}
	if initialPassword != "" {
		hash, err := crypto.HashPassword(initialPassword)
		if err != nil {
			t.Fatalf("hash initial password: %v", err)
		}
		store.passwordHash = hash
	}
	return store
}

func (s *fakeCredentialStore) AdminPasswordHash(_ context.Context) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.hashErr != nil {
		return "", s.hashErr
	}
	if s.passwordHash == "" {
		return "", contracts.ErrNotInitialized
	}
	return s.passwordHash, nil
}

func (s *fakeCredentialStore) UpdateAdminPasswordHash(_ context.Context, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.updateErr != nil {
		return s.updateErr
	}
	if s.passwordHash == "" {
		return contracts.ErrNotInitialized
	}
	s.passwordHash = hash
	return nil
}

func (s *fakeCredentialStore) RequireLogin(_ context.Context) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.reqErr != nil {
		return false, s.reqErr
	}
	return s.requireLogin, nil
}

func (s *fakeCredentialStore) SetRequireLogin(_ context.Context, required bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.setReqErr != nil {
		return s.setReqErr
	}
	s.requireLogin = required
	return nil
}

func (s *fakeCredentialStore) setPassword(t *testing.T, password string) {
	t.Helper()
	hash, err := crypto.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	s.mu.Lock()
	s.passwordHash = hash
	s.mu.Unlock()
}

type authHarness struct {
	mod    *auth.Module
	store  *fakeCredentialStore
	signer *crypto.SessionSigner
	bus    *eventbus.Bus
	srv    *httpserver.Server
}

func newTestSigner(t *testing.T, seed byte) *crypto.SessionSigner {
	t.Helper()
	key := bytes.Repeat([]byte{seed}, crypto.MasterKeySize)
	signer, err := crypto.NewSessionSigner(key)
	if err != nil {
		t.Fatalf("NewSessionSigner: %v", err)
	}
	return signer
}

func newAuthHarness(t *testing.T, cfg config.Config, initialPassword string, requireLogin bool) *authHarness {
	t.Helper()
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == 0 {
		cfg.Port = 20128
	}
	store := newFakeStore(t, initialPassword, requireLogin)
	signer := newTestSigner(t, 0x42)
	bus := eventbus.New()
	mod := auth.New(auth.Options{
		Store:    store,
		Sessions: signer,
		Config:   cfg,
		Bus:      bus,
	})
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("mod.Start: %v", err)
	}

	srv := httpserver.New(httpserver.Options{
		Config: cfg,
		Auth:   mod,
	})
	mod.RegisterRoutes(srv.Mounts())

	return &authHarness{
		mod:    mod,
		store:  store,
		signer: signer,
		bus:    bus,
		srv:    srv,
	}
}

func doRequest(t *testing.T, h http.Handler, method, target string, body string, cookies []*http.Cookie, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeErrorCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var env httpserver.ErrorEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal error envelope %q: %v", rec.Body.String(), err)
	}
	return env.Error.Code
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestAuthModuleMetadataAndLifecycle(t *testing.T) {
	t.Parallel()

	store := newFakeStore(t, "secret-pass", true)
	signer := newTestSigner(t, 0x11)

	mod := auth.New(auth.Options{
		Store:    store,
		Sessions: signer,
		Config:   config.Config{Host: "127.0.0.1", Port: 20128},
	})

	if got := mod.Name(); got != "auth" {
		t.Fatalf("Name() = %q, want %q", got, "auth")
	}
	deps := mod.Dependencies()
	if len(deps) != 1 || deps[0] != "system" {
		t.Fatalf("Dependencies() = %v, want [system]", deps)
	}
	if mod.Migrations() != nil {
		t.Fatalf("Migrations() = %v, want nil", mod.Migrations())
	}
	if err := mod.Start(context.Background()); err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	if err := mod.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() unexpected error: %v", err)
	}

	if err := auth.New(auth.Options{Sessions: signer}).Start(context.Background()); err == nil {
		t.Fatal("Start() with nil Store: expected error, got nil")
	}
	if err := auth.New(auth.Options{Store: store}).Start(context.Background()); err == nil {
		t.Fatal("Start() with nil Sessions: expected error, got nil")
	}
}

func TestAuthLoginLifecycle(t *testing.T) {
	t.Parallel()

	h := newAuthHarness(t, config.Config{Host: "127.0.0.1"}, "", true)

	// 1. Uninitialized -> 409 not_initialized
	rec := doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"correct-horse-battery"}`, nil, nil)
	if rec.Code != http.StatusConflict {
		t.Fatalf("login uninitialized status = %d, want %d (body: %s)", rec.Code, http.StatusConflict, rec.Body.String())
	}
	if code := decodeErrorCode(t, rec); code != "not_initialized" {
		t.Fatalf("login uninitialized code = %q, want %q", code, "not_initialized")
	}

	// GET /api/auth/me when uninitialized -> 200 with initialized: false
	rec = doRequest(t, h.srv.Handler(), http.MethodGet, "/api/auth/me", "", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("me uninitialized status = %d, want 200", rec.Code)
	}
	var meResp struct {
		Authenticated  bool   `json:"authenticated"`
		Subject        string `json:"subject"`
		RequireLogin   bool   `json:"require_login"`
		RemoteEnforced bool   `json:"remote_enforced"`
		Initialized    bool   `json:"initialized"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("unmarshal me response: %v", err)
	}
	if meResp.Authenticated || meResp.Subject != "" || meResp.Initialized || !meResp.RequireLogin || meResp.RemoteEnforced {
		t.Fatalf("unexpected uninitialized me response: %+v", meResp)
	}

	// Initialize admin password.
	h.store.setPassword(t, "correct-horse-battery")

	// 2. Malformed / unknown field / multi-object JSON -> 400; oversized body -> 413
	for _, badBody := range []string{
		`{invalid`,
		`{"password":"correct-horse-battery","extra":true}`,
		`{"password":"correct-horse-battery"}{"password":"other"}`,
	} {
		rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", badBody, nil, nil)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("login bad body %q status = %d, want 400", badBody, rec.Code)
		}
	}
	oversized := `{"password":"` + strings.Repeat("a", 70*1024) + `"}`
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", oversized, nil, nil)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("login oversized body status = %d, want 413", rec.Code)
	}

	// 3. Wrong password -> 401 invalid_credentials
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"wrong-password"}`, nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("login wrong password status = %d, want 401", rec.Code)
	}
	if code := decodeErrorCode(t, rec); code != "invalid_credentials" {
		t.Fatalf("login wrong password code = %q, want %q", code, "invalid_credentials")
	}

	// 4. Valid password over plain HTTP -> 200 + HttpOnly, SameSite=Lax, Path=/, Secure=false
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"correct-horse-battery"}`, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login valid status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var loginResp struct {
		Authenticated bool   `json:"authenticated"`
		Subject       string `json:"subject"`
		ExpiresAt     string `json:"expires_at"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("unmarshal login response: %v", err)
	}
	if !loginResp.Authenticated || loginResp.Subject != "admin" {
		t.Fatalf("unexpected login response: %+v", loginResp)
	}
	expTime, err := time.Parse(time.RFC3339, loginResp.ExpiresAt)
	if err != nil {
		t.Fatalf("invalid RFC3339 expires_at %q: %v", loginResp.ExpiresAt, err)
	}
	if !expTime.After(time.Now().UTC().Add(23 * time.Hour)) {
		t.Fatalf("expires_at %v is unexpectedly short for default 24h TTL", expTime)
	}

	sessionCookie := findCookie(rec.Result().Cookies(), auth.SessionCookieName)
	if sessionCookie == nil {
		t.Fatalf("expected Set-Cookie %q, got %v", auth.SessionCookieName, rec.Result().Cookies())
	}
	if !sessionCookie.HttpOnly {
		t.Fatal("session cookie must be HttpOnly")
	}
	if sessionCookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("session cookie SameSite = %v, want Lax", sessionCookie.SameSite)
	}
	if sessionCookie.Path != "/" {
		t.Fatalf("session cookie Path = %q, want /", sessionCookie.Path)
	}
	if sessionCookie.Secure {
		t.Fatal("session cookie over plain HTTP should not set Secure")
	}

	// 5. Valid password with X-Forwarded-Proto: https -> Secure=true
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"correct-horse-battery"}`, nil, map[string]string{
		"X-Forwarded-Proto": "https",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("login forwarded https status = %d, want 200", rec.Code)
	}
	if c := findCookie(rec.Result().Cookies(), auth.SessionCookieName); c == nil || !c.Secure {
		t.Fatalf("expected Secure=true cookie with X-Forwarded-Proto: https, got %+v", c)
	}

	// 6. Valid password with direct TLS -> Secure=true
	tlsReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"password":"correct-horse-battery"}`))
	tlsReq.Header.Set("Content-Type", "application/json")
	tlsReq.TLS = &tls.ConnectionState{}
	tlsRec := httptest.NewRecorder()
	h.srv.Handler().ServeHTTP(tlsRec, tlsReq)
	if tlsRec.Code != http.StatusOK {
		t.Fatalf("login direct TLS status = %d, want 200", tlsRec.Code)
	}
	if c := findCookie(tlsRec.Result().Cookies(), auth.SessionCookieName); c == nil || !c.Secure {
		t.Fatalf("expected Secure=true cookie with r.TLS != nil, got %+v", c)
	}

	// 7. GET /api/auth/me with cookie -> authenticated: true
	rec = doRequest(t, h.srv.Handler(), http.MethodGet, "/api/auth/me", "", []*http.Cookie{sessionCookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("me authenticated status = %d, want 200", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("unmarshal me response: %v", err)
	}
	if !meResp.Authenticated || meResp.Subject != "admin" || !meResp.Initialized || !meResp.RequireLogin || meResp.RemoteEnforced {
		t.Fatalf("unexpected authenticated me response: %+v", meResp)
	}

	// 8. POST /api/auth/logout clears cookie -> MaxAge: -1
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/logout", "", []*http.Cookie{sessionCookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", rec.Code)
	}
	var logoutResp struct {
		Authenticated bool `json:"authenticated"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &logoutResp); err != nil {
		t.Fatalf("unmarshal logout response: %v", err)
	}
	if logoutResp.Authenticated {
		t.Fatal("logout response authenticated = true, want false")
	}
	cleared := findCookie(rec.Result().Cookies(), auth.SessionCookieName)
	if cleared == nil || cleared.MaxAge != -1 || !cleared.HttpOnly || cleared.Path != "/" {
		t.Fatalf("expected cleared cookie with MaxAge=-1, got %+v", cleared)
	}

	// 9. GET /api/auth/me without cookie (or with cleared empty cookie) -> 200 authenticated: false
	rec = doRequest(t, h.srv.Handler(), http.MethodGet, "/api/auth/me", "", []*http.Cookie{cleared}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("me after logout status = %d, want 200", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("unmarshal me after logout: %v", err)
	}
	if meResp.Authenticated || meResp.Subject != "" {
		t.Fatalf("expected unauthenticated me after logout, got %+v", meResp)
	}
}

func TestAuthVerifierCookieValidation(t *testing.T) {
	t.Parallel()

	h := newAuthHarness(t, config.Config{Host: "127.0.0.1"}, "secret-pass", true)

	validToken, err := h.signer.Sign("admin", time.Hour)
	if err != nil {
		t.Fatalf("Sign valid: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: validToken})
	subj, err := h.mod.AuthenticateAdmin(req)
	if err != nil {
		t.Fatalf("AuthenticateAdmin valid cookie error: %v", err)
	}
	if subj != "admin" {
		t.Fatalf("AuthenticateAdmin subject = %q, want admin", subj)
	}

	expiredToken, err := h.signer.Sign("admin", -time.Second)
	if err != nil {
		t.Fatalf("Sign expired: %v", err)
	}
	otherSigner := newTestSigner(t, 0x99)
	otherKeyToken, err := otherSigner.Sign("admin", time.Hour)
	if err != nil {
		t.Fatalf("Sign other key: %v", err)
	}
	tamperedToken := validToken[:len(validToken)-1] + "X"
	if tamperedToken == validToken {
		tamperedToken = validToken[:len(validToken)-1] + "Y"
	}

	cases := []struct {
		name   string
		cookie *http.Cookie
	}{
		{name: "missing", cookie: nil},
		{name: "empty", cookie: &http.Cookie{Name: auth.SessionCookieName, Value: ""}},
		{name: "malformed", cookie: &http.Cookie{Name: auth.SessionCookieName, Value: "not-a-valid-token"}},
		{name: "tampered", cookie: &http.Cookie{Name: auth.SessionCookieName, Value: tamperedToken}},
		{name: "expired", cookie: &http.Cookie{Name: auth.SessionCookieName, Value: expiredToken}},
		{name: "wrong_signer", cookie: &http.Cookie{Name: auth.SessionCookieName, Value: otherKeyToken}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
			if tc.cookie != nil {
				r.AddCookie(tc.cookie)
			}
			gotSubj, authErr := h.mod.AuthenticateAdmin(r)
			if authErr == nil {
				t.Fatalf("AuthenticateAdmin(%s) succeeded with subject %q, want error", tc.name, gotSubj)
			}
		})
	}
}

func TestAuthRemoteBindEnforcement(t *testing.T) {
	t.Parallel()

	// Part 1: Non-loopback bind (0.0.0.0) with AllowInsecureRemote: false
	remote := newAuthHarness(t, config.Config{
		Host:                "0.0.0.0",
		Port:                20128,
		AllowInsecureRemote: false,
	}, "remote-admin-pass", false) // Note: underlying store has requireLogin=false!

	// Even though underlying store returns false, mod.RequireLogin MUST return true.
	if !remote.mod.RequireLogin(context.Background()) {
		t.Fatal("remote.mod.RequireLogin() = false, want true even when store returns false")
	}

	// GET /api/auth/me reflects enforced posture.
	rec := doRequest(t, remote.srv.Handler(), http.MethodGet, "/api/auth/me", "", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("remote me status = %d, want 200", rec.Code)
	}
	var meResp struct {
		RequireLogin   bool `json:"require_login"`
		RemoteEnforced bool `json:"remote_enforced"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &meResp); err != nil {
		t.Fatalf("unmarshal remote me: %v", err)
	}
	if !meResp.RequireLogin || !meResp.RemoteEnforced {
		t.Fatalf("remote me = %+v, want require_login=true remote_enforced=true", meResp)
	}

	// Unauthenticated PUT /api/auth/policy -> 401
	rec = doRequest(t, remote.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":false}`, nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated remote policy status = %d, want 401", rec.Code)
	}

	// Login to get session cookie.
	loginRec := doRequest(t, remote.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"remote-admin-pass"}`, nil, nil)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("remote login status = %d, want 200", loginRec.Code)
	}
	remoteCookie := findCookie(loginRec.Result().Cookies(), auth.SessionCookieName)
	if remoteCookie == nil {
		t.Fatal("expected session cookie on remote login")
	}

	// Authenticated PUT /api/auth/policy {"require_login":false} -> 403 remote_login_required
	rec = doRequest(t, remote.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":false}`, []*http.Cookie{remoteCookie}, nil)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("remote disable require_login status = %d, want 403 (body: %s)", rec.Code, rec.Body.String())
	}
	if code := decodeErrorCode(t, rec); code != "remote_login_required" {
		t.Fatalf("remote disable require_login code = %q, want %q", code, "remote_login_required")
	}

	// Authenticated PUT /api/auth/policy {"require_login":true} -> 200
	events, unsub := remote.bus.Subscribe(auth.TopicPolicyChanged, 4)
	defer unsub()
	rec = doRequest(t, remote.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":true}`, []*http.Cookie{remoteCookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("remote enable require_login status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	select {
	case ev := <-events:
		if ev.Topic != "auth.policy_changed" {
			t.Fatalf("event topic = %q, want auth.policy_changed", ev.Topic)
		}
	case <-time.After(time.Second):
		t.Fatal("expected auth.policy_changed event on eventbus")
	}

	// Part 2: Loopback bind (127.0.0.1) allows toggling require_login: false and unlocks AdminAPI without cookie.
	local := newAuthHarness(t, config.Config{
		Host: "127.0.0.1",
		Port: 20128,
	}, "local-admin-pass", true)

	// Before toggling, AdminAPI rejects unauthenticated request.
	rec = doRequest(t, local.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":false}`, nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("local unauthenticated policy before toggle status = %d, want 401", rec.Code)
	}

	loginRec = doRequest(t, local.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"local-admin-pass"}`, nil, nil)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("local login status = %d, want 200", loginRec.Code)
	}
	localCookie := findCookie(loginRec.Result().Cookies(), auth.SessionCookieName)

	localEvents, unsubLocal := local.bus.Subscribe(auth.TopicPolicyChanged, 4)
	defer unsubLocal()

	// Missing require_login field -> 400
	rec = doRequest(t, local.srv.Handler(), http.MethodPut, "/api/auth/policy", `{}`, []*http.Cookie{localCookie}, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("local empty policy body status = %d, want 400", rec.Code)
	}

	// Toggle require_login: false -> 200
	rec = doRequest(t, local.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":false}`, []*http.Cookie{localCookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("local disable require_login status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var policyResp struct {
		RequireLogin   bool `json:"require_login"`
		RemoteEnforced bool `json:"remote_enforced"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &policyResp); err != nil {
		t.Fatalf("unmarshal policy response: %v", err)
	}
	if policyResp.RequireLogin || policyResp.RemoteEnforced {
		t.Fatalf("unexpected local policy response: %+v", policyResp)
	}
	if local.mod.RequireLogin(context.Background()) {
		t.Fatal("local.mod.RequireLogin() = true after disabling, want false")
	}

	select {
	case ev := <-localEvents:
		if ev.Topic != "auth.policy_changed" {
			t.Fatalf("local event topic = %q, want auth.policy_changed", ev.Topic)
		}
	case <-time.After(time.Second):
		t.Fatal("expected auth.policy_changed event on local toggle")
	}

	// AdminAPI is now accessible without a session cookie.
	rec = doRequest(t, local.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":true}`, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("local unauthenticated AdminAPI after require_login=false status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}

	// Part 3: Non-loopback with AllowInsecureRemote: true permits disabling require_login.
	insecureRemote := newAuthHarness(t, config.Config{
		Host:                "0.0.0.0",
		Port:                20128,
		AllowInsecureRemote: true,
	}, "insecure-remote-pass", true)
	loginRec = doRequest(t, insecureRemote.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"insecure-remote-pass"}`, nil, nil)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("insecureRemote login status = %d, want 200", loginRec.Code)
	}
	insecureCookie := findCookie(loginRec.Result().Cookies(), auth.SessionCookieName)
	rec = doRequest(t, insecureRemote.srv.Handler(), http.MethodPut, "/api/auth/policy", `{"require_login":false}`, []*http.Cookie{insecureCookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("insecureRemote disable require_login status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	if insecureRemote.mod.RequireLogin(context.Background()) {
		t.Fatal("insecureRemote.mod.RequireLogin() = true after disabling with AllowInsecureRemote=true, want false")
	}
}

func TestAuthPasswordChange(t *testing.T) {
	t.Parallel()

	h := newAuthHarness(t, config.Config{Host: "127.0.0.1"}, "initial-password", true)

	// Unauthenticated -> 401
	rec := doRequest(t, h.srv.Handler(), http.MethodPut, "/api/auth/password",
		`{"current_password":"initial-password","new_password":"rotated-password"}`, nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated password change status = %d, want 401", rec.Code)
	}

	// Login to obtain cookie.
	loginRec := doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"initial-password"}`, nil, nil)
	if loginRec.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200", loginRec.Code)
	}
	cookie := findCookie(loginRec.Result().Cookies(), auth.SessionCookieName)

	// Wrong current_password -> 401
	rec = doRequest(t, h.srv.Handler(), http.MethodPut, "/api/auth/password",
		`{"current_password":"wrong-current","new_password":"rotated-password"}`, []*http.Cookie{cookie}, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong current_password status = %d, want 401 (body: %s)", rec.Code, rec.Body.String())
	}
	if code := decodeErrorCode(t, rec); code != "invalid_credentials" {
		t.Fatalf("wrong current_password code = %q, want invalid_credentials", code)
	}

	// Short / empty / >72-byte new_password -> 400 invalid_password
	badPasswords := []struct {
		name string
		pw   string
	}{
		{name: "empty", pw: ""},
		{name: "short_ascii", pw: "1234567"},
		{name: "short_multibyte_7_runes", pw: "あいうえおかき"}, // 7 runes, 21 bytes
		{name: "too_long_73_bytes", pw: strings.Repeat("a", 73)},
	}
	for _, tc := range badPasswords {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(map[string]string{
				"current_password": "initial-password",
				"new_password":     tc.pw,
			})
			r := doRequest(t, h.srv.Handler(), http.MethodPut, "/api/auth/password", string(payload), []*http.Cookie{cookie}, nil)
			if r.Code != http.StatusBadRequest {
				t.Fatalf("new_password %s status = %d, want 400 (body: %s)", tc.name, r.Code, r.Body.String())
			}
			if code := decodeErrorCode(t, r); code != "invalid_password" {
				t.Fatalf("new_password %s code = %q, want invalid_password", tc.name, code)
			}
		})
	}

	// Valid password rotation -> 200 + fresh session cookie
	rec = doRequest(t, h.srv.Handler(), http.MethodPut, "/api/auth/password",
		`{"current_password":"initial-password","new_password":"rotated-password-99"}`, []*http.Cookie{cookie}, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("valid password rotation status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var updatedResp struct {
		Updated bool `json:"updated"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &updatedResp); err != nil {
		t.Fatalf("unmarshal updated response: %v", err)
	}
	if !updatedResp.Updated {
		t.Fatal("updated = false, want true")
	}
	refreshedCookie := findCookie(rec.Result().Cookies(), auth.SessionCookieName)
	if refreshedCookie == nil || !refreshedCookie.HttpOnly || refreshedCookie.Value == "" {
		t.Fatalf("expected refreshed session cookie after password change, got %+v", refreshedCookie)
	}

	// Old password rejected on login.
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"initial-password"}`, nil, nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("login with old password status = %d, want 401", rec.Code)
	}

	// New password accepted on login.
	rec = doRequest(t, h.srv.Handler(), http.MethodPost, "/api/auth/login", `{"password":"rotated-password-99"}`, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("login with new password status = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
}
