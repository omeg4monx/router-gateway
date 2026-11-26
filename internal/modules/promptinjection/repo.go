package promptinjection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	rfc3339NanoUTC  = "2006-01-02T15:04:05.000000000Z"
	maxStoredEvents = 500
)

var errNotFound = errors.New("promptinjection: not found")

// Config holds the active prompt guard policy.
type Config struct {
	Enabled         bool   `json:"enabled"`
	Mode            string `json:"mode"`
	Sensitivity     string `json:"sensitivity"`
	BlockStatusCode int    `json:"block_status_code"`
	UpdatedAt       string `json:"updated_at"`
}

// Rule describes a single prompt injection detection rule.
type Rule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	MatchType   string `json:"match_type"`
	Action      string `json:"action"`
	Severity    string `json:"severity"`
	Enabled     bool   `json:"enabled"`
	Builtin     bool   `json:"builtin"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GuardEvent records a single prompt guard trigger in the bounded audit log.
type GuardEvent struct {
	ID        string `json:"id"`
	RequestID string `json:"request_id"`
	Model     string `json:"model"`
	KeyID     string `json:"key_id"`
	RuleID    string `json:"rule_id"`
	RuleName  string `json:"rule_name"`
	Action    string `json:"action"`
	Severity  string `json:"severity"`
	Snippet   string `json:"snippet"`
	CreatedAt string `json:"created_at"`
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func defaultConfig(now string) Config {
	return Config{
		Enabled:         false,
		Mode:            ModeAudit,
		Sensitivity:     SensitivityMedium,
		BlockStatusCode: 400,
		UpdatedAt:       now,
	}
}

func (r *repository) ensureDefaults(ctx context.Context) error {
	now := time.Now().UTC().Format(rfc3339NanoUTC)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT OR IGNORE INTO prompt_guard_config (id, enabled, mode, sensitivity, block_status_code, updated_at)
		VALUES (1, 0, ?, ?, 400, ?)`,
		ModeAudit, SensitivityMedium, now,
	)
	if err != nil {
		return fmt.Errorf("seed prompt_guard_config: %w", err)
	}

	for _, rule := range builtinRules(now) {
		_, err := r.db.Writer().ExecContext(ctx, `
			INSERT OR IGNORE INTO prompt_guard_rules
			(id, name, pattern, match_type, action, severity, enabled, builtin, description, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?)`,
			rule.ID, rule.Name, rule.Pattern, rule.MatchType, rule.Action, rule.Severity,
			boolToInt(rule.Enabled), rule.Description, rule.CreatedAt, rule.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("seed builtin rule %s: %w", rule.ID, err)
		}
	}
	return nil
}

func (r *repository) getConfig(ctx context.Context) (Config, error) {
	var (
		cfg        Config
		enabledInt int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT enabled, mode, sensitivity, block_status_code, updated_at
		FROM prompt_guard_config WHERE id = 1`).
		Scan(&enabledInt, &cfg.Mode, &cfg.Sensitivity, &cfg.BlockStatusCode, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return defaultConfig(time.Now().UTC().Format(rfc3339NanoUTC)), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read prompt_guard_config: %w", err)
	}
	cfg.Enabled = enabledInt != 0
	return cfg, nil
}

func (r *repository) saveConfig(ctx context.Context, cfg Config) (Config, error) {
	cfg.UpdatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO prompt_guard_config (id, enabled, mode, sensitivity, block_status_code, updated_at)
		VALUES (1, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			enabled = excluded.enabled,
			mode = excluded.mode,
			sensitivity = excluded.sensitivity,
			block_status_code = excluded.block_status_code,
			updated_at = excluded.updated_at`,
		boolToInt(cfg.Enabled), cfg.Mode, cfg.Sensitivity, cfg.BlockStatusCode, cfg.UpdatedAt,
	)
	if err != nil {
		return Config{}, fmt.Errorf("save prompt_guard_config: %w", err)
	}
	return cfg, nil
}

func (r *repository) listRules(ctx context.Context) ([]Rule, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, pattern, match_type, action, severity, enabled, builtin, description, created_at, updated_at
		FROM prompt_guard_rules
		ORDER BY builtin DESC, created_at ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list prompt_guard_rules: %w", err)
	}
	defer rows.Close()

	rules := make([]Rule, 0)
	for rows.Next() {
		var (
			item       Rule
			enabledInt int
			builtinInt int
		)
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Pattern, &item.MatchType, &item.Action, &item.Severity,
			&enabledInt, &builtinInt, &item.Description, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prompt_guard_rule: %w", err)
		}
		item.Enabled = enabledInt != 0
		item.Builtin = builtinInt != 0
		rules = append(rules, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate prompt_guard_rules: %w", err)
	}
	return rules, nil
}

func (r *repository) getRule(ctx context.Context, id string) (Rule, error) {
	var (
		item       Rule
		enabledInt int
		builtinInt int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, pattern, match_type, action, severity, enabled, builtin, description, created_at, updated_at
		FROM prompt_guard_rules WHERE id = ?`, id).
		Scan(
			&item.ID, &item.Name, &item.Pattern, &item.MatchType, &item.Action, &item.Severity,
			&enabledInt, &builtinInt, &item.Description, &item.CreatedAt, &item.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return Rule{}, errNotFound
	}
	if err != nil {
		return Rule{}, fmt.Errorf("get prompt_guard_rule: %w", err)
	}
	item.Enabled = enabledInt != 0
	item.Builtin = builtinInt != 0
	return item, nil
}

func (r *repository) createRule(ctx context.Context, item Rule) (Rule, error) {
	now := time.Now().UTC().Format(rfc3339NanoUTC)
	item.CreatedAt = now
	item.UpdatedAt = now
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO prompt_guard_rules
		(id, name, pattern, match_type, action, severity, enabled, builtin, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.Pattern, item.MatchType, item.Action, item.Severity,
		boolToInt(item.Enabled), boolToInt(item.Builtin), item.Description, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return Rule{}, fmt.Errorf("insert prompt_guard_rule: %w", err)
	}
	return item, nil
}

func (r *repository) updateRule(ctx context.Context, item Rule) (Rule, error) {
	item.UpdatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE prompt_guard_rules
		SET name = ?, pattern = ?, match_type = ?, action = ?, severity = ?, enabled = ?, description = ?, updated_at = ?
		WHERE id = ?`,
		item.Name, item.Pattern, item.MatchType, item.Action, item.Severity,
		boolToInt(item.Enabled), item.Description, item.UpdatedAt, item.ID,
	)
	if err != nil {
		return Rule{}, fmt.Errorf("update prompt_guard_rule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return Rule{}, errNotFound
	}
	return item, nil
}

func (r *repository) deleteRule(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM prompt_guard_rules WHERE id = ? AND builtin = 0`, id)
	if err != nil {
		return fmt.Errorf("delete prompt_guard_rule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errNotFound
	}
	return nil
}

func (r *repository) recordEvent(ctx context.Context, ev GuardEvent) error {
	if ev.CreatedAt == "" {
		ev.CreatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO prompt_guard_events
		(id, request_id, model, key_id, rule_id, rule_name, action, severity, snippet, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.ID, ev.RequestID, ev.Model, ev.KeyID, ev.RuleID, ev.RuleName, ev.Action, ev.Severity, ev.Snippet, ev.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert prompt_guard_event: %w", err)
	}
	_, _ = r.db.Writer().ExecContext(ctx, `
		DELETE FROM prompt_guard_events
		WHERE id NOT IN (
			SELECT id FROM prompt_guard_events ORDER BY created_at DESC LIMIT ?
		)`, maxStoredEvents)
	return nil
}

func (r *repository) listEvents(ctx context.Context, limit int) ([]GuardEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, request_id, model, key_id, rule_id, rule_name, action, severity, snippet, created_at
		FROM prompt_guard_events
		ORDER BY created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list prompt_guard_events: %w", err)
	}
	defer rows.Close()

	events := make([]GuardEvent, 0)
	for rows.Next() {
		var ev GuardEvent
		if err := rows.Scan(
			&ev.ID, &ev.RequestID, &ev.Model, &ev.KeyID, &ev.RuleID, &ev.RuleName,
			&ev.Action, &ev.Severity, &ev.Snippet, &ev.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan prompt_guard_event: %w", err)
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate prompt_guard_events: %w", err)
	}
	return events, nil
}

func (r *repository) clearEvents(ctx context.Context) error {
	_, err := r.db.Writer().ExecContext(ctx, `DELETE FROM prompt_guard_events`)
	if err != nil {
		return fmt.Errorf("clear prompt_guard_events: %w", err)
	}
	return nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
