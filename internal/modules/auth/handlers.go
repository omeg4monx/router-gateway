package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
)

const (
	minPasswordLength = 8
	maxAuthBodyBytes  = 64 << 10 // 64 KiB
)

type loginRequest struct {
	Password string `json:"password"`
}

type loginResponse struct {
	Authenticated bool   `json:"authenticated"`
	Subject       string `json:"subject"`
	ExpiresAt     string `json:"expires_at"`
}

type logoutResponse struct {
	Authenticated bool `json:"authenticated"`
}

type meResponse struct {
	Authenticated  bool   `json:"authenticated"`
	Subject        string `json:"subject"`
	RequireLogin   bool   `json:"require_login"`
	RemoteEnforced bool   `json:"remote_enforced"`
	Initialized    bool   `json:"initialized"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type changePasswordResponse struct {
	Updated bool `json:"updated"`
}

type updatePolicyRequest struct {
	RequireLogin *bool `json:"require_login"`
}

type policyResponse struct {
	RequireLogin   bool `json:"require_login"`
	RemoteEnforced bool `json:"remote_enforced"`
}

func (m *Module) handleLogin(w http.ResponseWriter, r *http.Request) {
	clientKey := clientRateLimitKey(r)
	now := time.Now().UTC()
	if allowed, retryAfter := m.loginLimiter.check(clientKey, now); !allowed {
		w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
		httpserver.WriteError(w, r, http.StatusTooManyRequests, "rate_limited", "too many failed login attempts; please try again later")
		return
	}

	var req loginRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if m.store == nil || m.sessions == nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "authentication module is not configured")
		return
	}

	hash, err := m.store.AdminPasswordHash(r.Context())
	if err != nil {
		if errors.Is(err, contracts.ErrNotInitialized) {
			httpserver.WriteError(w, r, http.StatusConflict, "not_initialized", "system is not initialized")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to read admin credentials")
		return
	}

	if !crypto.VerifyPassword(hash, req.Password) {
		m.loginLimiter.recordFailure(clientKey, now)
		httpserver.WriteError(w, r, http.StatusUnauthorized, "invalid_credentials", "invalid admin password")
		return
	}

	m.loginLimiter.reset(clientKey)

	expiresAt := time.Now().UTC().Add(m.sessionTTL)
	token, err := m.sessions.Sign(adminSubject, m.sessionTTL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to create session")
		return
	}

	setSessionCookie(w, r, token, m.sessionTTL, expiresAt)
	httpserver.WriteJSON(w, http.StatusOK, loginResponse{
		Authenticated: true,
		Subject:       adminSubject,
		ExpiresAt:     expiresAt.Format(time.RFC3339),
	})
}

func (m *Module) handleLogout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w, r)
	httpserver.WriteJSON(w, http.StatusOK, logoutResponse{
		Authenticated: false,
	})
}

func (m *Module) handleMe(w http.ResponseWriter, r *http.Request) {
	if m.store == nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "credential store is not configured")
		return
	}

	initialized := true
	if _, err := m.store.AdminPasswordHash(r.Context()); err != nil {
		if errors.Is(err, contracts.ErrNotInitialized) {
			initialized = false
		} else {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to read initialization state")
			return
		}
	}

	remoteEnforced := m.remoteEnforcedForRequest(r)
	requireLogin := true
	if !remoteEnforced {
		reqLogin, err := m.store.RequireLogin(r.Context())
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to read login policy")
			return
		}
		requireLogin = reqLogin
	}

	authenticated := false
	subject := ""
	if subj, err := m.AuthenticateAdmin(r); err == nil {
		authenticated = true
		subject = subj
	}

	httpserver.WriteJSON(w, http.StatusOK, meResponse{
		Authenticated:  authenticated,
		Subject:        subject,
		RequireLogin:   requireLogin,
		RemoteEnforced: remoteEnforced,
		Initialized:    initialized,
	})
}

func (m *Module) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req changePasswordRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	if m.store == nil || m.sessions == nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "authentication module is not configured")
		return
	}

	hash, err := m.store.AdminPasswordHash(r.Context())
	if err != nil {
		if errors.Is(err, contracts.ErrNotInitialized) {
			httpserver.WriteError(w, r, http.StatusConflict, "not_initialized", "system is not initialized")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to read admin credentials")
		return
	}

	if !crypto.VerifyPassword(hash, req.CurrentPassword) {
		httpserver.WriteError(w, r, http.StatusUnauthorized, "invalid_credentials", "current password is incorrect")
		return
	}

	if utf8.RuneCountInString(req.NewPassword) < minPasswordLength {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_password", "password must be at least 8 characters")
		return
	}

	newHash, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		if errors.Is(err, crypto.ErrPasswordTooLong) || errors.Is(err, crypto.ErrEmptyPassword) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_password", err.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to hash new password")
		return
	}

	if err := m.store.UpdateAdminPasswordHash(r.Context(), newHash); err != nil {
		if errors.Is(err, contracts.ErrNotInitialized) {
			httpserver.WriteError(w, r, http.StatusConflict, "not_initialized", "system is not initialized")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update admin password")
		return
	}

	subject := httpserver.AdminSubject(r.Context())
	if subject == "" {
		subject = adminSubject
	}
	expiresAt := time.Now().UTC().Add(m.sessionTTL)
	token, err := m.sessions.Sign(subject, m.sessionTTL)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to refresh session")
		return
	}

	setSessionCookie(w, r, token, m.sessionTTL, expiresAt)
	httpserver.WriteJSON(w, http.StatusOK, changePasswordResponse{
		Updated: true,
	})
}

func (m *Module) handleUpdatePolicy(w http.ResponseWriter, r *http.Request) {
	var req updatePolicyRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}
	if req.RequireLogin == nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "require_login is required")
		return
	}

	if m.store == nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "credential store is not configured")
		return
	}

	wantRequireLogin := *req.RequireLogin
	remoteEnforced := m.remoteEnforcedForRequest(r)
	if !wantRequireLogin && remoteEnforced {
		httpserver.WriteError(w, r, http.StatusForbidden, "remote_login_required", "login cannot be disabled on non-loopback bind addresses unless OMEG4MON_ALLOW_INSECURE_REMOTE is enabled")
		return
	}

	if err := m.store.SetRequireLogin(r.Context(), wantRequireLogin); err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to update login policy")
		return
	}

	if m.bus != nil {
		m.bus.Publish(TopicPolicyChanged, PolicyChangedEvent{
			RequireLogin:   wantRequireLogin,
			RemoteEnforced: remoteEnforced,
		})
	}

	httpserver.WriteJSON(w, http.StatusOK, policyResponse{
		RequireLogin:   wantRequireLogin,
		RemoteEnforced: remoteEnforced,
	})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return false
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return false
	}
	if dec.More() {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}
	return true
}
