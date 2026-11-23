package pxpipe

import (
	"context"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const timestampLayout = "2006-01-02T15:04:05.000000000Z"

// Config holds persisted transform pipeline settings.
type Config struct {
	Enabled             bool   `json:"enabled"`
	AutoStart           bool   `json:"auto_start"`
	StripNullFields     bool   `json:"strip_null_fields"`
	CollapseWhitespace  bool   `json:"collapse_whitespace"`
	DedupeSystemPrompts bool   `json:"dedupe_system_prompts"`
	MaxHistoryMessages  int    `json:"max_history_messages"`
	UpdatedAt           string `json:"updated_at"`
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) ensureDefaults(ctx context.Context) error {
	now := time.Now().UTC().Format(timestampLayout)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT OR IGNORE INTO pxpipe_config (
			id, enabled, auto_start, strip_null_fields, collapse_whitespace,
			dedupe_system_prompts, max_history_messages, updated_at
		) VALUES (1, 1, 1, 1, 1, 1, 0, ?)
	`, now)
	if err != nil {
		return fmt.Errorf("ensure pxpipe_config defaults: %w", err)
	}
	return nil
}

func (r *repository) getConfig(ctx context.Context) (Config, error) {
	var (
		cfg       Config
		enabled   int
		autoStart int
		stripNull int
		collapse  int
		dedupe    int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT enabled, auto_start, strip_null_fields, collapse_whitespace,
		       dedupe_system_prompts, max_history_messages, updated_at
		FROM pxpipe_config WHERE id = 1
	`).Scan(&enabled, &autoStart, &stripNull, &collapse, &dedupe, &cfg.MaxHistoryMessages, &cfg.UpdatedAt)
	if err != nil {
		return Config{}, fmt.Errorf("query pxpipe_config: %w", err)
	}
	cfg.Enabled = enabled != 0
	cfg.AutoStart = autoStart != 0
	cfg.StripNullFields = stripNull != 0
	cfg.CollapseWhitespace = collapse != 0
	cfg.DedupeSystemPrompts = dedupe != 0
	return cfg, nil
}

func (r *repository) updateConfig(ctx context.Context, cfg Config) (Config, error) {
	cfg.UpdatedAt = time.Now().UTC().Format(timestampLayout)
	_, err := r.db.Writer().ExecContext(ctx, `
		UPDATE pxpipe_config
		SET enabled = ?, auto_start = ?, strip_null_fields = ?, collapse_whitespace = ?,
		    dedupe_system_prompts = ?, max_history_messages = ?, updated_at = ?
		WHERE id = 1
	`, boolToInt(cfg.Enabled), boolToInt(cfg.AutoStart), boolToInt(cfg.StripNullFields),
		boolToInt(cfg.CollapseWhitespace), boolToInt(cfg.DedupeSystemPrompts),
		cfg.MaxHistoryMessages, cfg.UpdatedAt)
	if err != nil {
		return Config{}, fmt.Errorf("update pxpipe_config: %w", err)
	}
	return cfg, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
