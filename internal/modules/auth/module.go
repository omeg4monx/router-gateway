// Package auth implements single-admin cookie authentication, session
// verification, password rotation, and login enforcement policy management.
package auth

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "auth"

	// DefaultSessionTTL is the default validity window for signed admin sessions.
	DefaultSessionTTL = 24 * time.Hour

	// TopicPolicyChanged is published on the event bus when login policy changes.
	TopicPolicyChanged = "auth.policy_changed"

	adminSubject = "admin"
)

// ErrUnauthenticated indicates that a request lacks a valid admin session cookie.
var ErrUnauthenticated = errors.New("auth: missing or invalid session")

// PolicyChangedEvent is published on TopicPolicyChanged after persisting a
// login enforcement policy update.
type PolicyChangedEvent struct {
	RequireLogin   bool `json:"require_login"`
	RemoteEnforced bool `json:"remote_enforced"`
}

// Options configures a new auth Module instance.
type Options struct {
	Store      contracts.AdminCredentialStore
	Sessions   *crypto.SessionSigner
	Config     config.Config
	Bus        *eventbus.Bus
	SessionTTL time.Duration
}

// Module implements lifecycle.Module and contracts.AuthVerifier for the
// single-admin dashboard control plane.
type Module struct {
	store        contracts.AdminCredentialStore
	sessions     *crypto.SessionSigner
	cfg          config.Config
	bus          *eventbus.Bus
	sessionTTL   time.Duration
	loginLimiter *loginRateLimiter
}

var (
	_ lifecycle.Module       = (*Module)(nil)
	_ contracts.AuthVerifier = (*Module)(nil)
)

// New constructs an auth Module with the provided options, defaulting
// SessionTTL to DefaultSessionTTL when non-positive.
func New(opts Options) *Module {
	ttl := opts.SessionTTL
	if ttl <= 0 {
		ttl = DefaultSessionTTL
	}
	return &Module{
		store:        opts.Store,
		sessions:     opts.Sessions,
		cfg:          opts.Config,
		bus:          opts.Bus,
		sessionTTL:   ttl,
		loginLimiter: newLoginRateLimiter(defaultLoginMaxFailures, defaultLoginWindow),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that auth depends on the system module for admin
// credential storage.
func (m *Module) Dependencies() []string {
	return []string{"system"}
}

// Migrations returns nil because auth owns no SQL tables.
func (m *Module) Migrations() fs.FS {
	return nil
}

// RegisterRoutes mounts public and admin authentication endpoints onto the
// application's split API routers.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.PublicAPI != nil {
		mounts.PublicAPI.Post("/auth/login", m.handleLogin)
		mounts.PublicAPI.Post("/auth/logout", m.handleLogout)
		mounts.PublicAPI.Get("/auth/me", m.handleMe)
	}
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Put("/auth/password", m.handleChangePassword)
		mounts.AdminAPI.Put("/auth/policy", m.handleUpdatePolicy)
	}
}

// Start validates that required auth dependencies are wired.
func (m *Module) Start(context.Context) error {
	if m.store == nil {
		return errors.New("auth: admin credential store is required")
	}
	if m.sessions == nil {
		return errors.New("auth: session signer is required")
	}
	return nil
}

// Stop completes graceful shutdown for the auth module.
func (m *Module) Stop(context.Context) error {
	return nil
}

// AuthenticateAdmin verifies the signed admin session cookie on r and returns
// its authenticated subject.
func (m *Module) AuthenticateAdmin(r *http.Request) (string, error) {
	if m.sessions == nil {
		return "", errors.New("auth: session signer is not configured")
	}
	token, err := readSessionCookie(r)
	if err != nil {
		return "", err
	}
	claims, err := m.sessions.Verify(token)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUnauthenticated, err)
	}
	return claims.Subject, nil
}

// RequireLogin reports whether dashboard authentication is currently enforced.
// Non-loopback binds without AllowInsecureRemote always enforce login even if
// the underlying store reports false. Store lookup errors fail closed.
func (m *Module) RequireLogin(ctx context.Context) bool {
	if m.remoteEnforced() {
		return true
	}
	if m.store == nil {
		return true
	}
	required, err := m.store.RequireLogin(ctx)
	if err != nil {
		return true
	}
	return required
}

func (m *Module) remoteEnforced() bool {
	return !m.cfg.AllowInsecureRemote && (!m.cfg.IsLoopback() || m.cfg.TLSEnabled())
}

func (m *Module) remoteEnforcedForRequest(r *http.Request) bool {
	return m.remoteEnforced() || (!m.cfg.AllowInsecureRemote && httpserver.IsHTTPS(r))
}
