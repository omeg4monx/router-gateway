package tokensaver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoUTC = "2006-01-02T15:04:05.000000000Z"

// Config holds the persisted Token Saver + Headroom configuration.
type Config struct {
	Enabled            bool   `json:"enabled"`
	Strategy           string `json:"strategy"`
	MaxHistoryMessages int    `json:"max_history_messages"`
	MaxToolOutputChars int    `json:"max_tool_output_chars"`
	HeadroomEnabled    bool   `json:"headroom_enabled"`
	PrefixDepth        int    `json:"prefix_depth"`
	MaxCacheEntries    int    `json:"max_cache_entries"`
	CacheTTLSeconds    int    `json:"cache_ttl_seconds"`
	InjectCacheControl bool   `json:"inject_cache_control"`
	CavemanEnabled     bool   `json:"caveman_enabled"`
	CavemanLevel       string `json:"caveman_level"`
	PonytailEnabled    bool   `json:"ponytail_enabled"`
	PonytailLevel      string `json:"ponytail_level"`
	TrimHistoryEnabled bool   `json:"trim_history_enabled"`
	UpdatedAt          string `json:"updated_at"`
}

func defaultConfig(now string) Config {
	return Config{
		Enabled:            false,
		Strategy:           StrategyWhitespace,
		MaxHistoryMessages: 40,
		MaxToolOutputChars: 8000,
		HeadroomEnabled:    true,
		PrefixDepth:        2,
		MaxCacheEntries:    1024,
		CacheTTLSeconds:    1800,
		InjectCacheControl: true,
		CavemanEnabled:     false,
		CavemanLevel:       LevelFull,
		PonytailEnabled:    false,
		PonytailLevel:      LevelFull,
		TrimHistoryEnabled: false,
		UpdatedAt:          now,
	}
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) ensureDefaults(ctx context.Context) error {
	now := time.Now().UTC().Format(rfc3339NanoUTC)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT OR IGNORE INTO tokensaver_config
		(id, enabled, strategy, max_history_messages, max_tool_output_chars,
		 headroom_enabled, prefix_depth, max_cache_entries, cache_ttl_seconds, inject_cache_control,
		 caveman_enabled, caveman_level, ponytail_enabled, ponytail_level, trim_history_enabled, updated_at)
		VALUES (1, 0, ?, 40, 8000, 1, 2, 1024, 1800, 1, 0, ?, 0, ?, 0, ?)`,
		StrategyWhitespace, LevelFull, LevelFull, now,
	)
	if err != nil {
		return fmt.Errorf("seed tokensaver_config: %w", err)
	}
	return nil
}

func (r *repository) getConfig(ctx context.Context) (Config, error) {
	var (
		cfg             Config
		enabledInt      int
		headroomInt     int
		injectCacheCtrl int
		cavemanInt      int
		ponytailInt     int
		trimHistoryInt  int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT enabled, strategy, max_history_messages, max_tool_output_chars,
		       headroom_enabled, prefix_depth, max_cache_entries, cache_ttl_seconds, inject_cache_control,
		       caveman_enabled, caveman_level, ponytail_enabled, ponytail_level, trim_history_enabled, updated_at
		FROM tokensaver_config WHERE id = 1`).
		Scan(
			&enabledInt, &cfg.Strategy, &cfg.MaxHistoryMessages, &cfg.MaxToolOutputChars,
			&headroomInt, &cfg.PrefixDepth, &cfg.MaxCacheEntries, &cfg.CacheTTLSeconds, &injectCacheCtrl,
			&cavemanInt, &cfg.CavemanLevel, &ponytailInt, &cfg.PonytailLevel, &trimHistoryInt, &cfg.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return defaultConfig(time.Now().UTC().Format(rfc3339NanoUTC)), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read tokensaver_config: %w", err)
	}
	cfg.Enabled = enabledInt != 0
	cfg.HeadroomEnabled = headroomInt != 0
	cfg.InjectCacheControl = injectCacheCtrl != 0
	cfg.CavemanEnabled = cavemanInt != 0
	cfg.PonytailEnabled = ponytailInt != 0
	cfg.TrimHistoryEnabled = trimHistoryInt != 0
	if cfg.CavemanLevel == "" {
		cfg.CavemanLevel = LevelFull
	}
	if cfg.PonytailLevel == "" {
		cfg.PonytailLevel = LevelFull
	}
	return cfg, nil
}

func (r *repository) saveConfig(ctx context.Context, cfg Config) (Config, error) {
	cfg.UpdatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO tokensaver_config
		(id, enabled, strategy, max_history_messages, max_tool_output_chars,
		 headroom_enabled, prefix_depth, max_cache_entries, cache_ttl_seconds, inject_cache_control,
		 caveman_enabled, caveman_level, ponytail_enabled, ponytail_level, trim_history_enabled, updated_at)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			enabled = excluded.enabled,
			strategy = excluded.strategy,
			max_history_messages = excluded.max_history_messages,
			max_tool_output_chars = excluded.max_tool_output_chars,
			headroom_enabled = excluded.headroom_enabled,
			prefix_depth = excluded.prefix_depth,
			max_cache_entries = excluded.max_cache_entries,
			cache_ttl_seconds = excluded.cache_ttl_seconds,
			inject_cache_control = excluded.inject_cache_control,
			caveman_enabled = excluded.caveman_enabled,
			caveman_level = excluded.caveman_level,
			ponytail_enabled = excluded.ponytail_enabled,
			ponytail_level = excluded.ponytail_level,
			trim_history_enabled = excluded.trim_history_enabled,
			updated_at = excluded.updated_at`,
		boolToInt(cfg.Enabled), cfg.Strategy, cfg.MaxHistoryMessages, cfg.MaxToolOutputChars,
		boolToInt(cfg.HeadroomEnabled), cfg.PrefixDepth, cfg.MaxCacheEntries, cfg.CacheTTLSeconds,
		boolToInt(cfg.InjectCacheControl), boolToInt(cfg.CavemanEnabled), cfg.CavemanLevel,
		boolToInt(cfg.PonytailEnabled), cfg.PonytailLevel, boolToInt(cfg.TrimHistoryEnabled), cfg.UpdatedAt,
	)
	if err != nil {
		return Config{}, fmt.Errorf("save tokensaver_config: %w", err)
	}
	return cfg, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
