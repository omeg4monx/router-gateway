// Package importer implements a read-only, one-way legacy data.sqlite inspector
// and importer that migrates settings, combos, and API keys through narrow
// cross-module contracts without mutating the source database.
package importer

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "importer"

	// TopicImportCompleted is published when a legacy import run completes.
	TopicImportCompleted = "importer.completed"

	maxBodyBytes = 64 << 10 // 64 KiB
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Options configures an importer Module instance.
type Options struct {
	DB             *db.DB
	Bus            *eventbus.Bus
	Logger         *slog.Logger
	Settings       contracts.SettingsSnapshotter
	Combos         contracts.ComboSnapshotter
	KeyProvisioner contracts.KeyProvisioner
	Now            func() time.Time
}

// Module implements lifecycle.Module for optional one-way legacy data.sqlite import.
type Module struct {
	db       *db.DB
	bus      *eventbus.Bus
	logger   *slog.Logger
	settings contracts.SettingsSnapshotter
	combos   contracts.ComboSnapshotter
	keys     contracts.KeyProvisioner
	repo     *repository
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs an importer Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &Module{
		db:       opts.DB,
		bus:      opts.Bus,
		logger:   logger,
		settings: opts.Settings,
		combos:   opts.Combos,
		keys:     opts.KeyProvisioner,
		repo:     newRepository(opts.DB, opts.Now),
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares module startup dependencies.
func (m *Module) Dependencies() []string { return []string{"system", "settings", "combos", "keys"} }

// Migrations returns the embedded SQL migrations.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("importer: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts importer endpoints on AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Post("/importer/inspect", m.handleInspect)
	mounts.AdminAPI.Post("/importer/run", m.handleRun)
	mounts.AdminAPI.Get("/importer/history", m.handleListHistory)
	mounts.AdminAPI.Delete("/importer/history", m.handleClearHistory)
}

// Start validates required dependencies.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("importer: database is required")
	}
	return nil
}

// Stop completes graceful shutdown.
func (m *Module) Stop(context.Context) error { return nil }

type inspectRequest struct {
	SourcePath string `json:"source_path"`
}

func (m *Module) handleInspect(w http.ResponseWriter, r *http.Request) {
	var req inspectRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	cleanPath, err := validateSQLiteSource(req.SourcePath)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_source_path", err.Error())
		return
	}
	ext, err := readLegacySQLite(r.Context(), cleanPath)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "legacy_inspect_failed", err.Error())
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"source_path":     cleanPath,
		"tables_found":    ext.TablesFound,
		"settings_count":  len(ext.Settings),
		"combos_count":    len(ext.Combos),
		"keys_count":      len(ext.Keys),
		"providers_count": ext.ProvidersCount,
		"warnings":        ext.Warnings,
	})
}

type runImportRequest struct {
	SourcePath     string `json:"source_path"`
	DryRun         bool   `json:"dry_run"`
	ImportSettings bool   `json:"import_settings"`
	ImportCombos   bool   `json:"import_combos"`
	ImportKeys     bool   `json:"import_keys"`
	Overwrite      bool   `json:"overwrite"`
}

func (m *Module) handleRun(w http.ResponseWriter, r *http.Request) {
	var req runImportRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	cleanPath, err := validateSQLiteSource(req.SourcePath)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_source_path", err.Error())
		return
	}
	ext, err := readLegacySQLite(r.Context(), cleanPath)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "legacy_read_failed", err.Error())
		return
	}

	var (
		settingsCount int
		combosCount   int
		keysCount     int
		provisioned   = make([]contracts.ProvisionedKey, 0)
	)

	if req.DryRun {
		if req.ImportSettings {
			settingsCount = len(ext.Settings)
		}
		if req.ImportCombos {
			combosCount = len(ext.Combos)
		}
		if req.ImportKeys {
			keysCount = len(ext.Keys)
		}
	} else {
		if req.ImportSettings && m.settings != nil && len(ext.Settings) > 0 {
			n, err := m.settings.ImportSettings(r.Context(), ext.Settings, "importer", "legacy data.sqlite import")
			if err != nil {
				_, _ = m.repo.InsertRun(r.Context(), ImportRun{
					SourcePath: cleanPath,
					DryRun:     false,
					Status:     "error",
					Summary:    "Settings import failed",
					Error:      err.Error(),
				})
				httpserver.WriteError(w, r, http.StatusInternalServerError, "import_settings_failed", err.Error())
				return
			}
			settingsCount = n
		}
		if req.ImportCombos && m.combos != nil && len(ext.Combos) > 0 {
			cCount, _, err := m.combos.ImportCombos(r.Context(), ext.Combos, nil, req.Overwrite, "importer")
			if err != nil {
				_, _ = m.repo.InsertRun(r.Context(), ImportRun{
					SourcePath:       cleanPath,
					DryRun:           false,
					Status:           "error",
					SettingsImported: settingsCount,
					Summary:          "Combos import failed",
					Error:            err.Error(),
				})
				httpserver.WriteError(w, r, http.StatusInternalServerError, "import_combos_failed", err.Error())
				return
			}
			combosCount = cCount
		}
		if req.ImportKeys && m.keys != nil && len(ext.Keys) > 0 {
			for _, k := range ext.Keys {
				pk, err := m.keys.ProvisionKey(r.Context(), contracts.ProvisionKeyRequest{
					Name:          k.Name,
					AllowedModels: k.AllowedModels,
					Tags:          []string{"legacy-import"},
					RateLimitRPM:  k.RateLimitRPM,
					Actor:         "importer",
				})
				if err != nil {
					ext.Warnings = append(ext.Warnings, fmt.Sprintf("key %q skipped: %v", k.Name, err))
					continue
				}
				provisioned = append(provisioned, pk)
				keysCount++
			}
		}
	}

	modeLabel := "live"
	if req.DryRun {
		modeLabel = "dry-run"
	}
	summary := fmt.Sprintf("Completed %s import: %d settings, %d combos, %d keys", modeLabel, settingsCount, combosCount, keysCount)
	run, err := m.repo.InsertRun(r.Context(), ImportRun{
		SourcePath:       cleanPath,
		DryRun:           req.DryRun,
		Status:           "ok",
		SettingsImported: settingsCount,
		CombosImported:   combosCount,
		KeysImported:     keysCount,
		ProvidersFound:   ext.ProvidersCount,
		Summary:          summary,
	})
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "importer_history_failed", "Failed to record import run")
		return
	}
	if !req.DryRun && m.bus != nil {
		m.bus.Publish(TopicImportCompleted, run)
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"run":              run,
		"warnings":         ext.Warnings,
		"provisioned_keys": provisioned,
	})
}

func (m *Module) handleListHistory(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	runs, err := m.repo.ListRuns(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "importer_history_failed", "Failed to load import history")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"runs": runs})
}

func (m *Module) handleClearHistory(w http.ResponseWriter, r *http.Request) {
	deleted, err := m.repo.ClearRuns(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "importer_clear_failed", "Failed to clear import history")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		return errors.New("failed to read request body")
	}
	if len(raw) > maxBodyBytes {
		return errors.New("request body exceeds size limit")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}
