package mitm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	timestampLayout   = "2006-01-02T15:04:05.000000000Z"
	defaultListenAddr = "127.0.0.1:8443"
	maxStoredEvents   = 500
)

var (
	// ErrAliasConflict indicates an alias for (tool, source_model) already exists.
	ErrAliasConflict = errors.New("mitm: alias already exists for tool and source model")
	// ErrAliasNotFound indicates no alias matches the requested ID.
	ErrAliasNotFound = errors.New("mitm: alias not found")
)

// Config represents persisted MITM proxy configuration.
type Config struct {
	ListenAddr       string   `json:"listen_addr"`
	TargetGatewayURL string   `json:"target_gateway_url"`
	EnabledTools     []string `json:"enabled_tools"`
	AutoStart        bool     `json:"auto_start"`
	UpdatedAt        string   `json:"updated_at"`
}

// Alias maps a tool's requested model name to a target router model/combo.
type Alias struct {
	ID          string `json:"id"`
	Tool        string `json:"tool"`
	SourceModel string `json:"source_model"`
	TargetModel string `json:"target_model"`
	Enabled     bool   `json:"enabled"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Event records an intercepted MITM request.
type Event struct {
	ID          string `json:"id"`
	Tool        string `json:"tool"`
	Host        string `json:"host"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	SourceModel string `json:"source_model"`
	TargetModel string `json:"target_model"`
	StatusCode  int    `json:"status_code"`
	LatencyMs   int64  `json:"latency_ms"`
	CreatedAt   string `json:"created_at"`
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) ensureDefaults(ctx context.Context, defaultTarget string) error {
	now := time.Now().UTC().Format(timestampLayout)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT OR IGNORE INTO mitm_config (id, listen_addr, target_gateway_url, enabled_tools_json, auto_start, updated_at)
		VALUES (1, ?, ?, '["antigravity","copilot","kiro","cursor"]', 0, ?)
	`, defaultListenAddr, defaultTarget, now)
	if err != nil {
		return fmt.Errorf("ensure mitm_config defaults: %w", err)
	}
	return nil
}

func (r *repository) getConfig(ctx context.Context) (Config, error) {
	var (
		cfg       Config
		toolsJSON string
		autoStart int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT listen_addr, target_gateway_url, enabled_tools_json, auto_start, updated_at
		FROM mitm_config WHERE id = 1
	`).Scan(&cfg.ListenAddr, &cfg.TargetGatewayURL, &toolsJSON, &autoStart, &cfg.UpdatedAt)
	if err != nil {
		return Config{}, fmt.Errorf("query mitm_config: %w", err)
	}
	cfg.AutoStart = autoStart != 0
	if err := json.Unmarshal([]byte(toolsJSON), &cfg.EnabledTools); err != nil {
		cfg.EnabledTools = []string{}
	}
	if cfg.EnabledTools == nil {
		cfg.EnabledTools = []string{}
	}
	return cfg, nil
}

func (r *repository) updateConfig(ctx context.Context, cfg Config) (Config, error) {
	toolsJSON, err := json.Marshal(cfg.EnabledTools)
	if err != nil {
		return Config{}, fmt.Errorf("marshal enabled tools: %w", err)
	}
	cfg.UpdatedAt = time.Now().UTC().Format(timestampLayout)
	autoStart := 0
	if cfg.AutoStart {
		autoStart = 1
	}
	_, err = r.db.Writer().ExecContext(ctx, `
		UPDATE mitm_config
		SET listen_addr = ?, target_gateway_url = ?, enabled_tools_json = ?, auto_start = ?, updated_at = ?
		WHERE id = 1
	`, cfg.ListenAddr, cfg.TargetGatewayURL, string(toolsJSON), autoStart, cfg.UpdatedAt)
	if err != nil {
		return Config{}, fmt.Errorf("update mitm_config: %w", err)
	}
	return cfg, nil
}

func (r *repository) listAliases(ctx context.Context) ([]Alias, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, tool, source_model, target_model, enabled, created_at, updated_at
		FROM mitm_aliases
		ORDER BY tool ASC, source_model ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list mitm_aliases: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]Alias, 0)
	for rows.Next() {
		var (
			a       Alias
			enabled int
		)
		if err := rows.Scan(&a.ID, &a.Tool, &a.SourceModel, &a.TargetModel, &enabled, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan mitm_alias: %w", err)
		}
		a.Enabled = enabled != 0
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *repository) createAlias(ctx context.Context, a Alias) (Alias, error) {
	now := time.Now().UTC().Format(timestampLayout)
	a.CreatedAt = now
	a.UpdatedAt = now
	enabled := 0
	if a.Enabled {
		enabled = 1
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO mitm_aliases (id, tool, source_model, target_model, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, a.ID, a.Tool, a.SourceModel, a.TargetModel, enabled, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return Alias{}, ErrAliasConflict
		}
		return Alias{}, fmt.Errorf("insert mitm_alias: %w", err)
	}
	return a, nil
}

func (r *repository) deleteAlias(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM mitm_aliases WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete mitm_alias: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrAliasNotFound
	}
	return nil
}

func (r *repository) resolveAlias(ctx context.Context, tool, sourceModel string) (string, bool) {
	var target string
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT target_model
		FROM mitm_aliases
		WHERE enabled = 1 AND source_model = ? AND tool IN (?, '*')
		ORDER BY CASE WHEN tool = ? THEN 0 ELSE 1 END
		LIMIT 1
	`, sourceModel, tool, tool).Scan(&target)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", false
		}
		return "", false
	}
	return target, true
}

func (r *repository) recordEvent(ctx context.Context, ev Event) error {
	if ev.CreatedAt == "" {
		ev.CreatedAt = time.Now().UTC().Format(timestampLayout)
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO mitm_events (id, tool, host, method, path, source_model, target_model, status_code, latency_ms, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, ev.ID, ev.Tool, ev.Host, ev.Method, ev.Path, ev.SourceModel, ev.TargetModel, ev.StatusCode, ev.LatencyMs, ev.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert mitm_event: %w", err)
	}
	_, _ = r.db.Writer().ExecContext(ctx, `
		DELETE FROM mitm_events
		WHERE id NOT IN (SELECT id FROM mitm_events ORDER BY created_at DESC LIMIT ?)
	`, maxStoredEvents)
	return nil
}

func (r *repository) listEvents(ctx context.Context, limit int) ([]Event, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, tool, host, method, path, source_model, target_model, status_code, latency_ms, created_at
		FROM mitm_events
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list mitm_events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]Event, 0)
	for rows.Next() {
		var ev Event
		if err := rows.Scan(&ev.ID, &ev.Tool, &ev.Host, &ev.Method, &ev.Path, &ev.SourceModel, &ev.TargetModel, &ev.StatusCode, &ev.LatencyMs, &ev.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan mitm_event: %w", err)
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}

func (r *repository) clearEvents(ctx context.Context) (int64, error) {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM mitm_events`)
	if err != nil {
		return 0, fmt.Errorf("clear mitm_events: %w", err)
	}
	return res.RowsAffected()
}
