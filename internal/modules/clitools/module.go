// Package clitools implements stateless CLI coding-tool configuration
// generators and loopback-only atomic local apply/rollback operations.
package clitools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/omeg4monx/router-gateway/internal/platform/config"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName           = "clitools"
	maxCLIToolsBodyBytes = 64 << 10 // 64 KiB
)

// Options configures a new clitools Module instance.
type Options struct {
	Config  config.Config
	HomeDir string
}

// Module implements lifecycle.Module for CLI tool configuration generation
// and loopback-guarded local apply/rollback.
type Module struct {
	cfg     config.Config
	homeDir string
	homeErr error
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a clitools Module, defaulting HomeDir to os.UserHomeDir() when
// omitted.
func New(opts Options) *Module {
	homeDir := strings.TrimSpace(opts.HomeDir)
	var homeErr error
	if homeDir == "" {
		homeDir, homeErr = os.UserHomeDir()
	}
	return &Module{
		cfg:     opts.Config,
		homeDir: homeDir,
		homeErr: homeErr,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string {
	return moduleName
}

// Dependencies declares that clitools depends on the system module.
func (m *Module) Dependencies() []string {
	return []string{"system"}
}

// Migrations returns nil because clitools owns no SQL tables.
func (m *Module) Migrations() fs.FS {
	return nil
}

// RegisterRoutes mounts CLI tool endpoints onto AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI != nil {
		mounts.AdminAPI.Get("/cli-tools", m.handleListTools)
		mounts.AdminAPI.Post("/cli-tools/{tool}/generate", m.handleGenerate)
		mounts.AdminAPI.Post("/cli-tools/{tool}/apply", m.handleApply)
		mounts.AdminAPI.Post("/cli-tools/{tool}/rollback", m.handleRollback)
	}
}

// Start verifies the resolved home directory when local apply is enabled.
func (m *Module) Start(context.Context) error {
	if m.localApplyEnabled() && m.homeErr != nil {
		return fmt.Errorf("clitools: resolve home directory: %w", m.homeErr)
	}
	return nil
}

// Stop completes graceful shutdown for the clitools module.
func (m *Module) Stop(context.Context) error {
	return nil
}

func (m *Module) localApplyEnabled() bool {
	return m.cfg.IsLoopback()
}

type catalogResponse struct {
	LocalApplyEnabled bool             `json:"local_apply_enabled"`
	Tools             []ToolDescriptor `json:"tools"`
}

type configRequest struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
}

type applyResponse struct {
	Applied    bool   `json:"applied"`
	Tool       string `json:"tool"`
	TargetPath string `json:"target_path"`
	BackupPath string `json:"backup_path"`
}

type rollbackResponse struct {
	Restored   bool   `json:"restored"`
	Tool       string `json:"tool"`
	TargetPath string `json:"target_path"`
	BackupPath string `json:"backup_path"`
}

func (m *Module) handleListTools(w http.ResponseWriter, _ *http.Request) {
	httpserver.WriteJSON(w, http.StatusOK, catalogResponse{
		LocalApplyEnabled: m.localApplyEnabled(),
		Tools:             listToolDescriptors(m.homeDir),
	})
}

func (m *Module) handleGenerate(w http.ResponseWriter, r *http.Request) {
	spec, ok := m.resolveToolParam(w, r)
	if !ok {
		return
	}

	var req configRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	res, err := generateToolConfig(spec, m.homeDir, generateInput(req))
	if err != nil {
		writeGenerateError(w, r, err)
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, res)
}

func (m *Module) handleApply(w http.ResponseWriter, r *http.Request) {
	if !m.localApplyEnabled() {
		httpserver.WriteError(w, r, http.StatusForbidden, "local_apply_disabled", "local configuration apply is only permitted on loopback binds")
		return
	}

	spec, ok := m.resolveToolParam(w, r)
	if !ok {
		return
	}

	if m.homeErr != nil || m.homeDir == "" {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "home directory is unavailable for local apply")
		return
	}

	var req configRequest
	if !decodeJSONBody(w, r, &req) {
		return
	}

	res, err := generateToolConfig(spec, m.homeDir, generateInput(req))
	if err != nil {
		writeGenerateError(w, r, err)
		return
	}

	backupPath, err := applyConfigFile(res.TargetPath, res.Content)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to apply tool configuration")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, applyResponse{
		Applied:    true,
		Tool:       spec.ID,
		TargetPath: res.TargetPath,
		BackupPath: backupPath,
	})
}

func (m *Module) handleRollback(w http.ResponseWriter, r *http.Request) {
	if !m.localApplyEnabled() {
		httpserver.WriteError(w, r, http.StatusForbidden, "local_apply_disabled", "local configuration rollback is only permitted on loopback binds")
		return
	}

	spec, ok := m.resolveToolParam(w, r)
	if !ok {
		return
	}

	if m.homeErr != nil || m.homeDir == "" {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "home directory is unavailable for local rollback")
		return
	}

	targetPath := spec.targetPath(m.homeDir)
	backupPath, err := rollbackConfigFile(targetPath)
	if err != nil {
		if errors.Is(err, ErrBackupNotFound) {
			httpserver.WriteError(w, r, http.StatusNotFound, "backup_not_found", "no backup file found for this tool")
			return
		}
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to rollback tool configuration")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, rollbackResponse{
		Restored:   true,
		Tool:       spec.ID,
		TargetPath: targetPath,
		BackupPath: backupPath,
	})
}

func (m *Module) resolveToolParam(w http.ResponseWriter, r *http.Request) (toolSpec, bool) {
	toolID := chi.URLParam(r, "tool")
	spec, err := lookupToolSpec(toolID)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusNotFound, "tool_not_found", fmt.Sprintf("unknown CLI tool %q", toolID))
		return toolSpec{}, false
	}
	return spec, true
}

func writeGenerateError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidBaseURL):
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_base_url", "base_url must be a valid http:// or https:// URL with a non-empty host")
	case errors.Is(err, ErrInvalidInput):
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_request", err.Error())
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "internal_error", "failed to generate tool configuration")
	}
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxCLIToolsBodyBytes)
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
