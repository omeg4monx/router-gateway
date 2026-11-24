// Package cloudsync implements optional remote state synchronization and
// portable JSON snapshot export/restore across control-plane modules.
package cloudsync

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
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
	"github.com/omeg4monx/router-gateway/internal/platform/crypto"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
	"github.com/omeg4monx/router-gateway/internal/platform/eventbus"
	"github.com/omeg4monx/router-gateway/internal/platform/httpserver"
	"github.com/omeg4monx/router-gateway/internal/platform/lifecycle"
)

const (
	moduleName = "cloudsync"

	// TopicCloudSyncCompleted is published after a push, pull, or restore completes.
	TopicCloudSyncCompleted = "cloudsync.completed"

	maxRequestBodyBytes  = 2 << 20 // 2 MiB
	maxRemoteBodyBytes   = 4 << 20 // 4 MiB
	defaultSyncTimeout   = 20 * time.Second
	snapshotSchemaFormat = 1
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// SnapshotBundle is the portable JSON state document exchanged with remote sync
// endpoints or exported/restored directly via the dashboard.
type SnapshotBundle struct {
	Version    int                       `json:"version"`
	ExportedAt string                    `json:"exported_at"`
	Checksum   string                    `json:"checksum"`
	Settings   map[string]any            `json:"settings,omitempty"`
	Combos     []contracts.ComboSnapshot `json:"combos,omitempty"`
	Aliases    []contracts.AliasSnapshot `json:"aliases,omitempty"`
}

// Options configures a cloudsync Module instance.
type Options struct {
	DB                    *db.DB
	Vault                 *crypto.Vault
	HTTPClient            *http.Client
	Bus                   *eventbus.Bus
	Logger                *slog.Logger
	Settings              contracts.SettingsSnapshotter
	Combos                contracts.ComboSnapshotter
	Now                   func() time.Time
	AllowPrivateEndpoints bool
}

// Module implements lifecycle.Module for optional remote state sync.
type Module struct {
	db           *db.DB
	vault        *crypto.Vault
	httpClient   *http.Client
	bus          *eventbus.Bus
	logger       *slog.Logger
	settings     contracts.SettingsSnapshotter
	combos       contracts.ComboSnapshotter
	repo         *repository
	allowPrivate bool
}

var _ lifecycle.Module = (*Module)(nil)

// New constructs a cloudsync Module.
func New(opts Options) *Module {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: defaultSyncTimeout}
	}
	return &Module{
		db:           opts.DB,
		vault:        opts.Vault,
		httpClient:   client,
		bus:          opts.Bus,
		logger:       logger,
		settings:     opts.Settings,
		combos:       opts.Combos,
		repo:         newRepository(opts.DB, opts.Now),
		allowPrivate: opts.AllowPrivateEndpoints,
	}
}

// Name returns the canonical module name.
func (m *Module) Name() string { return moduleName }

// Dependencies declares module startup dependencies.
func (m *Module) Dependencies() []string { return []string{"system", "settings", "combos"} }

// Migrations returns the embedded SQL migrations.
func (m *Module) Migrations() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(fmt.Sprintf("cloudsync: load embedded migrations: %v", err))
	}
	return sub
}

// RegisterRoutes mounts cloudsync endpoints on AdminAPI.
func (m *Module) RegisterRoutes(mounts lifecycle.RouteMounts) {
	if mounts.AdminAPI == nil {
		return
	}
	mounts.AdminAPI.Get("/cloud-sync/config", m.handleGetConfig)
	mounts.AdminAPI.Patch("/cloud-sync/config", m.handlePatchConfig)
	mounts.AdminAPI.Post("/cloud-sync/push", m.handlePush)
	mounts.AdminAPI.Post("/cloud-sync/pull", m.handlePull)
	mounts.AdminAPI.Get("/cloud-sync/snapshot", m.handleExportSnapshot)
	mounts.AdminAPI.Post("/cloud-sync/restore", m.handleRestoreSnapshot)
	mounts.AdminAPI.Get("/cloud-sync/history", m.handleListHistory)
	mounts.AdminAPI.Delete("/cloud-sync/history", m.handleClearHistory)
}

// Start validates required dependencies.
func (m *Module) Start(context.Context) error {
	if m.db == nil {
		return errors.New("cloudsync: database is required")
	}
	if m.vault == nil {
		return errors.New("cloudsync: vault is required")
	}
	return nil
}

// Stop completes graceful shutdown.
func (m *Module) Stop(context.Context) error { return nil }

func (m *Module) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	sc, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load cloud sync config")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, sc.Config)
}

type patchConfigRequest struct {
	Enabled         *bool   `json:"enabled"`
	EndpointURL     *string `json:"endpoint_url"`
	AuthToken       *string `json:"auth_token"`
	ClearAuthToken  bool    `json:"clear_auth_token"`
	SyncIntervalMin *int    `json:"sync_interval_min"`
	ConflictPolicy  *string `json:"conflict_policy"`
	IncludeSettings *bool   `json:"include_settings"`
	IncludeCombos   *bool   `json:"include_combos"`
}

func (m *Module) handlePatchConfig(w http.ResponseWriter, r *http.Request) {
	var req patchConfigRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	cur, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load current config")
		return
	}
	if req.Enabled != nil {
		cur.Enabled = *req.Enabled
	}
	if req.EndpointURL != nil {
		trimmed := strings.TrimSpace(*req.EndpointURL)
		if trimmed != "" {
			if _, err := validateEndpointURL(r.Context(), trimmed, m.allowPrivate); err != nil {
				httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_endpoint_url", err.Error())
				return
			}
		}
		cur.EndpointURL = trimmed
	}
	if req.ClearAuthToken {
		cur.TokenEnc = ""
		cur.AuthTokenMask = ""
	} else if req.AuthToken != nil && strings.TrimSpace(*req.AuthToken) != "" {
		rawTok := strings.TrimSpace(*req.AuthToken)
		enc, err := m.vault.Encrypt([]byte(rawTok))
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_encrypt_failed", "Failed to encrypt auth token")
			return
		}
		cur.TokenEnc = enc
		cur.AuthTokenMask = maskSecret(rawTok)
	}
	if req.SyncIntervalMin != nil {
		if *req.SyncIntervalMin < 1 || *req.SyncIntervalMin > 1440 {
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_interval", "sync_interval_min must be between 1 and 1440")
			return
		}
		cur.SyncIntervalMin = *req.SyncIntervalMin
	}
	if req.ConflictPolicy != nil {
		switch strings.ToLower(strings.TrimSpace(*req.ConflictPolicy)) {
		case policyMerge, policyRemoteWins, policyLocalWins:
			cur.ConflictPolicy = strings.ToLower(strings.TrimSpace(*req.ConflictPolicy))
		default:
			httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_conflict_policy", "conflict_policy must be merge, remote_wins, or local_wins")
			return
		}
	}
	if req.IncludeSettings != nil {
		cur.IncludeSettings = *req.IncludeSettings
	}
	if req.IncludeCombos != nil {
		cur.IncludeCombos = *req.IncludeCombos
	}

	saved, err := m.repo.SaveConfig(r.Context(), cur)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_save_failed", "Failed to save cloud sync config")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, saved.Config)
}

func (m *Module) buildSnapshot(ctx context.Context, cfg Config) (SnapshotBundle, int, error) {
	bundle := SnapshotBundle{
		Version:    snapshotSchemaFormat,
		ExportedAt: m.repo.nowUTC(),
	}
	items := 0
	if cfg.IncludeSettings && m.settings != nil {
		vals, err := m.settings.ExportSettings(ctx)
		if err != nil {
			return SnapshotBundle{}, 0, fmt.Errorf("export settings: %w", err)
		}
		bundle.Settings = vals
		items += len(vals)
	}
	if cfg.IncludeCombos && m.combos != nil {
		combos, aliases, err := m.combos.ExportCombos(ctx)
		if err != nil {
			return SnapshotBundle{}, 0, fmt.Errorf("export combos: %w", err)
		}
		bundle.Combos = combos
		bundle.Aliases = aliases
		items += len(combos) + len(aliases)
	}
	bundle.Checksum = computeBundleChecksum(bundle)
	return bundle, items, nil
}

func (m *Module) applySnapshot(ctx context.Context, bundle SnapshotBundle, policy string, cfg Config) (int, error) {
	applied := 0
	overwrite := policy == policyRemoteWins

	if cfg.IncludeSettings && m.settings != nil && len(bundle.Settings) > 0 {
		toApply := bundle.Settings
		if policy == policyLocalWins {
			existing, err := m.settings.ExportSettings(ctx)
			if err != nil {
				return applied, fmt.Errorf("read local settings: %w", err)
			}
			filtered := make(map[string]any)
			for k, v := range bundle.Settings {
				if _, exists := existing[k]; !exists {
					filtered[k] = v
				}
			}
			toApply = filtered
		}
		if len(toApply) > 0 {
			n, err := m.settings.ImportSettings(ctx, toApply, "cloudsync", "cloudsync "+policy)
			if err != nil {
				return applied, fmt.Errorf("import settings: %w", err)
			}
			applied += n
		}
	}

	if cfg.IncludeCombos && m.combos != nil && (len(bundle.Combos) > 0 || len(bundle.Aliases) > 0) {
		cCount, aCount, err := m.combos.ImportCombos(ctx, bundle.Combos, bundle.Aliases, overwrite, "cloudsync")
		if err != nil {
			return applied, fmt.Errorf("import combos: %w", err)
		}
		applied += cCount + aCount
	}
	return applied, nil
}

func (m *Module) handlePush(w http.ResponseWriter, r *http.Request) {
	sc, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load config")
		return
	}
	if strings.TrimSpace(sc.EndpointURL) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "endpoint_not_configured", "Configure endpoint_url before pushing")
		return
	}
	if _, err := validateEndpointURL(r.Context(), sc.EndpointURL, m.allowPrivate); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_endpoint_url", err.Error())
		return
	}

	bundle, itemCount, err := m.buildSnapshot(r.Context(), sc.Config)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "snapshot_build_failed", err.Error())
		return
	}
	payload, err := json.Marshal(bundle)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "snapshot_encode_failed", "Failed to encode snapshot")
		return
	}

	reqCtx, cancel := context.WithTimeout(r.Context(), defaultSyncTimeout)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPut, sc.EndpointURL, bytes.NewReader(payload))
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request_build_failed", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if sc.TokenEnc != "" {
		tokBytes, err := m.vault.Decrypt(sc.TokenEnc)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "Failed to decrypt auth token")
			return
		}
		httpReq.Header.Set("Authorization", "Bearer "+string(tokBytes))
	}

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		_ = m.repo.RecordSyncStatus(r.Context(), "error", err.Error())
		_, _ = m.repo.InsertHistory(r.Context(), "push", "error", "Remote push failed", bundle.Checksum, itemCount, err.Error())
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_push_failed", "Failed to reach remote sync endpoint")
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("remote returned HTTP %d", resp.StatusCode)
		_ = m.repo.RecordSyncStatus(r.Context(), "error", errMsg)
		_, _ = m.repo.InsertHistory(r.Context(), "push", "error", "Remote rejected snapshot", bundle.Checksum, itemCount, errMsg)
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_rejected", errMsg)
		return
	}

	_ = m.repo.RecordSyncStatus(r.Context(), "ok", "")
	ev, _ := m.repo.InsertHistory(r.Context(), "push", "ok", fmt.Sprintf("Pushed %d items to remote", itemCount), bundle.Checksum, itemCount, "")
	if m.bus != nil {
		m.bus.Publish(TopicCloudSyncCompleted, ev)
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"direction":  "push",
		"checksum":   bundle.Checksum,
		"item_count": itemCount,
		"event":      ev,
	})
}

func (m *Module) handlePull(w http.ResponseWriter, r *http.Request) {
	sc, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load config")
		return
	}
	if strings.TrimSpace(sc.EndpointURL) == "" {
		httpserver.WriteError(w, r, http.StatusBadRequest, "endpoint_not_configured", "Configure endpoint_url before pulling")
		return
	}
	if _, err := validateEndpointURL(r.Context(), sc.EndpointURL, m.allowPrivate); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_endpoint_url", err.Error())
		return
	}

	reqCtx, cancel := context.WithTimeout(r.Context(), defaultSyncTimeout)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, sc.EndpointURL, http.NoBody)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "request_build_failed", err.Error())
		return
	}
	if sc.TokenEnc != "" {
		tokBytes, err := m.vault.Decrypt(sc.TokenEnc)
		if err != nil {
			httpserver.WriteError(w, r, http.StatusInternalServerError, "vault_decrypt_failed", "Failed to decrypt auth token")
			return
		}
		httpReq.Header.Set("Authorization", "Bearer "+string(tokBytes))
	}

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		_ = m.repo.RecordSyncStatus(r.Context(), "error", err.Error())
		_, _ = m.repo.InsertHistory(r.Context(), "pull", "error", "Remote pull failed", "", 0, err.Error())
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_pull_failed", "Failed to fetch remote snapshot")
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxRemoteBodyBytes+1))
	if err != nil || len(raw) > maxRemoteBodyBytes {
		errMsg := "remote snapshot exceeds size limit or failed to read"
		_ = m.repo.RecordSyncStatus(r.Context(), "error", errMsg)
		_, _ = m.repo.InsertHistory(r.Context(), "pull", "error", "Remote pull failed", "", 0, errMsg)
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_payload_invalid", errMsg)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("remote returned HTTP %d", resp.StatusCode)
		_ = m.repo.RecordSyncStatus(r.Context(), "error", errMsg)
		_, _ = m.repo.InsertHistory(r.Context(), "pull", "error", "Remote rejected pull", "", 0, errMsg)
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_rejected", errMsg)
		return
	}

	var bundle SnapshotBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		errMsg := "remote payload is not a valid snapshot JSON"
		_ = m.repo.RecordSyncStatus(r.Context(), "error", errMsg)
		_, _ = m.repo.InsertHistory(r.Context(), "pull", "error", "Invalid remote snapshot", "", 0, errMsg)
		httpserver.WriteError(w, r, http.StatusBadGateway, "remote_snapshot_invalid", errMsg)
		return
	}
	if bundle.Checksum == "" {
		bundle.Checksum = computeBundleChecksum(bundle)
	}

	applied, err := m.applySnapshot(r.Context(), bundle, sc.ConflictPolicy, sc.Config)
	if err != nil {
		_ = m.repo.RecordSyncStatus(r.Context(), "error", err.Error())
		_, _ = m.repo.InsertHistory(r.Context(), "pull", "error", "Apply snapshot failed", bundle.Checksum, applied, err.Error())
		httpserver.WriteError(w, r, http.StatusInternalServerError, "snapshot_apply_failed", err.Error())
		return
	}

	_ = m.repo.RecordSyncStatus(r.Context(), "ok", "")
	ev, _ := m.repo.InsertHistory(r.Context(), "pull", "ok", fmt.Sprintf("Pulled and applied %d items (%s)", applied, sc.ConflictPolicy), bundle.Checksum, applied, "")
	if m.bus != nil {
		m.bus.Publish(TopicCloudSyncCompleted, ev)
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"direction":  "pull",
		"checksum":   bundle.Checksum,
		"item_count": applied,
		"event":      ev,
	})
}

func (m *Module) handleExportSnapshot(w http.ResponseWriter, r *http.Request) {
	sc, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load config")
		return
	}
	bundle, itemCount, err := m.buildSnapshot(r.Context(), sc.Config)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "snapshot_build_failed", err.Error())
		return
	}
	_, _ = m.repo.InsertHistory(r.Context(), "export", "ok", fmt.Sprintf("Exported %d items", itemCount), bundle.Checksum, itemCount, "")
	httpserver.WriteJSON(w, http.StatusOK, bundle)
}

type restoreSnapshotRequest struct {
	ConflictPolicy string         `json:"conflict_policy"`
	Snapshot       SnapshotBundle `json:"snapshot"`
}

func (m *Module) handleRestoreSnapshot(w http.ResponseWriter, r *http.Request) {
	var req restoreSnapshotRequest
	if err := decodeJSON(r, &req); err != nil {
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	sc, err := m.repo.GetConfig(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_config_failed", "Failed to load config")
		return
	}
	policy := strings.ToLower(strings.TrimSpace(req.ConflictPolicy))
	if policy == "" {
		policy = sc.ConflictPolicy
	}
	switch policy {
	case policyMerge, policyRemoteWins, policyLocalWins:
	default:
		httpserver.WriteError(w, r, http.StatusBadRequest, "invalid_conflict_policy", "conflict_policy must be merge, remote_wins, or local_wins")
		return
	}
	if req.Snapshot.Checksum == "" {
		req.Snapshot.Checksum = computeBundleChecksum(req.Snapshot)
	}
	applied, err := m.applySnapshot(r.Context(), req.Snapshot, policy, sc.Config)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "snapshot_restore_failed", err.Error())
		return
	}
	ev, _ := m.repo.InsertHistory(r.Context(), "restore", "ok", fmt.Sprintf("Restored %d items (%s)", applied, policy), req.Snapshot.Checksum, applied, "")
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":         true,
		"direction":  "restore",
		"checksum":   req.Snapshot.Checksum,
		"item_count": applied,
		"event":      ev,
	})
}

func (m *Module) handleListHistory(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	events, err := m.repo.ListHistory(r.Context(), limit)
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_history_failed", "Failed to load sync history")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"events": events})
}

func (m *Module) handleClearHistory(w http.ResponseWriter, r *http.Request) {
	deleted, err := m.repo.ClearHistory(r.Context())
	if err != nil {
		httpserver.WriteError(w, r, http.StatusInternalServerError, "cloudsync_clear_failed", "Failed to clear sync history")
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, map[string]any{"deleted": deleted})
}

func computeBundleChecksum(b SnapshotBundle) string {
	copyBundle := SnapshotBundle{
		Version:  b.Version,
		Settings: b.Settings,
		Combos:   b.Combos,
		Aliases:  b.Aliases,
	}
	raw, _ := json.Marshal(copyBundle)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:8])
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBodyBytes+1))
	if err != nil {
		return errors.New("failed to read request body")
	}
	if len(raw) > maxRequestBodyBytes {
		return errors.New("request body exceeds size limit")
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}
