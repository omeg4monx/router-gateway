// Package system implements the core system lifecycle module exposing health,
// build version, first-run initialization, and guarded shutdown endpoints.
package system

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"runtime"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName        = "system"
	minPasswordLength = 8
	maxPasswordBytes  = 72
	maxInitBodyBytes  = 64 << 10 // 64 KiB
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// BuildInfo carries compile-time release metadata surfaced by GET /api/version.
type BuildInfo struct {
	Version   string
	Commit    string
	BuildTime string
}

// Options configures a new system Module instance.
type Options struct {
	DB            *db.DB
	Config        config.Config
	Build         BuildInfo
	ActiveModules func() []string
	Shutdown      func()
}

// Module implements lifecycle.Module for core system endpoints and state.
type Module struct {
	db            *db.DB
	repo          *repository
	cfg           config.Config
	build         BuildInfo
	activeModules func() []string
	shutdown      func()

	startedMu    sync.RWMutex
	startedAt    time.Time
	shutdownOnce sync.Once
	initMu       sync.Mutex
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a system Module with the provided dependencies and metadata.
func New(opts Options) *Module {
	return &Module{
		db:            opts.DB,
		repo:          newRepository(opts.DB),
		cfg:           opts.Config,
		build:         opts.Build,
		activeModules: opts.ActiveModules,
		shutdown:      opts.Shutdown,
		startedAt:     time.Now(),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies returns the module's lifecycle dependencies (none).
func (m *Module) Dependencies() []string {
	return nil
}

// Migrations returns the embedded SQL migration filesystem rooted at migrations/.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("system: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts the public and admin system endpoints onto the
// application's split API routers.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.PublicAPI != nil {
		mounts.PublicAPI.Get("/health", m.handleHealth)
		mounts.PublicAPI.Get("/version", m.handleVersion)
		mounts.PublicAPI.Get("/init", m.handleGetInit)
		mounts.PublicAPI.Post("/init", m.handlePostInit)
	}
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Post("/shutdown", m.handleShutdown)
	}
}

// Start ensures default system state is populated and, when Config.AdminPassword
// is configured, atomically bootstraps the initial admin password hash.
func (m *Module) Start(ctx context.Context) error {
	if m.db == nil {
		return errors.New("system: database is required")
	}

	m.startedMu.Lock()
	m.startedAt = time.Now()
	m.startedMu.Unlock()

	if err := m.repo.EnsureDefaultRequireLogin(ctx); err != nil {
		return err
	}

	if m.cfg.AdminPassword != "" {
		if err := validatePassword(m.cfg.AdminPassword); err != nil {
			return fmt.Errorf("system: invalid OMEG4MON_ADMIN_PASSWORD: %w", err)
		}

		initialized, err := m.repo.IsInitialized(ctx)
		if err != nil {
			return err
		}
		if initialized {
			return nil
		}

		hash, err := crypto.HashPassword(m.cfg.AdminPassword)
		if err != nil {
			return fmt.Errorf("system: hash bootstrap admin password: %w", err)
		}
		if err := m.repo.BootstrapAdminPassword(ctx, hash); err != nil {
			return err
		}
	}

	return nil
}

// Stop completes graceful shutdown for the system module.
func (m *Module) Stop(context.Context) error {
	return nil
}

type healthResponse struct {
	Status        string   `json:"status"`
	DB            string   `json:"db"`
	Modules       []string `json:"modules"`
	UptimeSeconds int64    `json:"uptimeSeconds"`
}

type versionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	GoVersion string `json:"goVersion"`
}

type initStatusResponse struct {
	Initialized  bool `json:"initialized"`
	RequireLogin bool `json:"requireLogin"`
}

type initRequest struct {
	Password string `json:"password"`
}

type initCompleteResponse struct {
	Initialized bool `json:"initialized"`
}

type shutdownResponse struct {
	Status string `json:"status"`
}

func (m *Module) handleHealth(w http.ResponseWriter, r *http.Request) {
	modules := m.snapshotModules()
	uptime := m.uptimeSeconds()

	if m.db == nil || m.db.Ping(r.Context()) != nil {
		httpserver.WriteJSON(w, http.StatusServiceUnavailable, healthResponse{
			Status:        "degraded",
			DB:            "error",
			Modules:       modules,
			UptimeSeconds: uptime,
		})
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, healthResponse{
		Status:        "ok",
		DB:            "ok",
		Modules:       modules,
		UptimeSeconds: uptime,
	})
}

func (m *Module) handleVersion(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, versionResponse{
		Version:   m.build.Version,
		Commit:    m.build.Commit,
		BuildTime: m.build.BuildTime,
		GoVersion: runtime.Version(),
	})
}

func (m *Module) handleGetInit(w http.ResponseWriter, r *http.Request) {
	initialized, requireLogin, err := m.repo.InitState(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to read initialization state")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, initStatusResponse{
		Initialized:  initialized,
		RequireLogin: requireLogin,
	})
}

func (m *Module) handlePostInit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxInitBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req initRequest
	if err := dec.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpserver.WriteError(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds maximum allowed size")
			return
		}
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return
	}

	if err := validatePassword(req.Password); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_password", err.Error())
		return
	}

	initialized, err := m.repo.IsInitialized(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to verify initialization state")
		return
	}
	if initialized {
		httpserver.WriteError(w, r, http.StatusConflict, "already_initialized", "system is already initialized")
		return
	}

	m.initMu.Lock()
	defer m.initMu.Unlock()

	initialized, err = m.repo.IsInitialized(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to verify initialization state")
		return
	}
	if initialized {
		httpserver.WriteError(w, r, http.StatusConflict, "already_initialized", "system is already initialized")
		return
	}

	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		if errors.Is(err, crypto.ErrPasswordTooLong) || errors.Is(err, crypto.ErrEmptyPassword) {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_password", err.Error())
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to hash password")
		return
	}

	if err := m.repo.InitAdminPassword(r.Context(), hash); err != nil {
		if errors.Is(err, ErrAlreadyInitialized) {
			httpserver.WriteError(w, r, http.StatusConflict, "already_initialized", "system is already initialized")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to save initialization state")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, initCompleteResponse{
		Initialized: true,
	})
}

func (m *Module) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if !m.cfg.IsLoopback() && !m.cfg.AllowInsecureRemote {
		httpserver.WriteError(w, r, http.StatusForbidden, "forbidden", "shutdown endpoint is disabled on non-loopback bind addresses")
		return
	}

	httpserver.WriteJSON(w, http.StatusAccepted, shutdownResponse{
		Status: "shutting_down",
	})

	m.shutdownOnce.Do(func() {
		if m.shutdown != nil {
			go m.shutdown()
		}
	})
}

func (m *Module) snapshotModules() []string {
	if m.activeModules == nil {
		return []string{}
	}
	raw := m.activeModules()
	if len(raw) == 0 {
		return []string{}
	}
	out := make([]string, len(raw))
	copy(out, raw)
	return out
}

func (m *Module) uptimeSeconds() int64 {
	m.startedMu.RLock()
	started := m.startedAt
	m.startedMu.RUnlock()

	secs := int64(time.Since(started).Seconds())
	if secs < 0 {
		return 0
	}
	return secs
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}
	if len(password) > maxPasswordBytes {
		return fmt.Errorf("password must not exceed %d bytes", maxPasswordBytes)
	}
	return nil
}
