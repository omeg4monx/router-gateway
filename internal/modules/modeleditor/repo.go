package modeleditor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoUTC = "2006-01-02T15:04:05.000000000Z"

var errNotFound = errors.New("modeleditor: not found")

// ModelOverride defines custom model parameters, prompt wrappers, or target rewrites.
type ModelOverride struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	ModelPattern       string         `json:"model_pattern"`
	TargetModel        string         `json:"target_model"`
	SystemPromptPrefix string         `json:"system_prompt_prefix"`
	SystemPromptSuffix string         `json:"system_prompt_suffix"`
	Temperature        *float64       `json:"temperature,omitempty"`
	MaxTokens          *int           `json:"max_tokens,omitempty"`
	TopP               *float64       `json:"top_p,omitempty"`
	ExtraParams        map[string]any `json:"extra_params"`
	StripParams        []string       `json:"strip_params"`
	Priority           int            `json:"priority"`
	Enabled            bool           `json:"enabled"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
	compiledPattern    *regexp.Regexp `json:"-"`
}

// BypassRule conditionally rewrites models or bypasses downstream interceptors.
type BypassRule struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	ModelPattern    string         `json:"model_pattern"`
	KeyID           string         `json:"key_id"`
	Protocol        string         `json:"protocol"`
	HeaderName      string         `json:"header_name"`
	HeaderValue     string         `json:"header_value"`
	Action          string         `json:"action"`
	TargetModel     string         `json:"target_model"`
	Priority        int            `json:"priority"`
	Enabled         bool           `json:"enabled"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
	compiledPattern *regexp.Regexp `json:"-"`
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) listOverrides(ctx context.Context) ([]ModelOverride, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, model_pattern, target_model, system_prompt_prefix, system_prompt_suffix,
		       temperature, max_tokens, top_p, extra_params_json, strip_params_json,
		       priority, enabled, created_at, updated_at
		FROM model_overrides
		ORDER BY priority ASC, created_at ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list model_overrides: %w", err)
	}
	defer rows.Close()

	items := make([]ModelOverride, 0)
	for rows.Next() {
		item, scanErr := scanOverride(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate model_overrides: %w", err)
	}
	return items, nil
}

func (r *repository) getOverride(ctx context.Context, id string) (ModelOverride, error) {
	row := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, model_pattern, target_model, system_prompt_prefix, system_prompt_suffix,
		       temperature, max_tokens, top_p, extra_params_json, strip_params_json,
		       priority, enabled, created_at, updated_at
		FROM model_overrides WHERE id = ?`, id)
	item, err := scanOverride(row)
	if errors.Is(err, sql.ErrNoRows) {
		return ModelOverride{}, errNotFound
	}
	return item, err
}

func (r *repository) createOverride(ctx context.Context, item ModelOverride) (ModelOverride, error) {
	now := time.Now().UTC().Format(rfc3339NanoUTC)
	item.CreatedAt = now
	item.UpdatedAt = now
	extraJSON, stripJSON, err := marshalOverrideJSON(item)
	if err != nil {
		return ModelOverride{}, err
	}
	_, err = r.db.Writer().ExecContext(ctx, `
		INSERT INTO model_overrides
		(id, name, model_pattern, target_model, system_prompt_prefix, system_prompt_suffix,
		 temperature, max_tokens, top_p, extra_params_json, strip_params_json,
		 priority, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.ModelPattern, item.TargetModel, item.SystemPromptPrefix, item.SystemPromptSuffix,
		nullFloat(item.Temperature), nullInt(item.MaxTokens), nullFloat(item.TopP),
		extraJSON, stripJSON, item.Priority, boolToInt(item.Enabled), item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return ModelOverride{}, fmt.Errorf("insert model_override: %w", err)
	}
	return item, nil
}

func (r *repository) updateOverride(ctx context.Context, item ModelOverride) (ModelOverride, error) {
	item.UpdatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	extraJSON, stripJSON, err := marshalOverrideJSON(item)
	if err != nil {
		return ModelOverride{}, err
	}
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE model_overrides
		SET name = ?, model_pattern = ?, target_model = ?, system_prompt_prefix = ?, system_prompt_suffix = ?,
		    temperature = ?, max_tokens = ?, top_p = ?, extra_params_json = ?, strip_params_json = ?,
		    priority = ?, enabled = ?, updated_at = ?
		WHERE id = ?`,
		item.Name, item.ModelPattern, item.TargetModel, item.SystemPromptPrefix, item.SystemPromptSuffix,
		nullFloat(item.Temperature), nullInt(item.MaxTokens), nullFloat(item.TopP),
		extraJSON, stripJSON, item.Priority, boolToInt(item.Enabled), item.UpdatedAt, item.ID,
	)
	if err != nil {
		return ModelOverride{}, fmt.Errorf("update model_override: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ModelOverride{}, errNotFound
	}
	return item, nil
}

func (r *repository) deleteOverride(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM model_overrides WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete model_override: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errNotFound
	}
	return nil
}

func (r *repository) listBypassRules(ctx context.Context) ([]BypassRule, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, model_pattern, key_id, protocol, header_name, header_value,
		       action, target_model, priority, enabled, created_at, updated_at
		FROM model_bypass_rules
		ORDER BY priority ASC, created_at ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list model_bypass_rules: %w", err)
	}
	defer rows.Close()

	items := make([]BypassRule, 0)
	for rows.Next() {
		item, scanErr := scanBypass(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate model_bypass_rules: %w", err)
	}
	return items, nil
}

func (r *repository) getBypassRule(ctx context.Context, id string) (BypassRule, error) {
	row := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, model_pattern, key_id, protocol, header_name, header_value,
		       action, target_model, priority, enabled, created_at, updated_at
		FROM model_bypass_rules WHERE id = ?`, id)
	item, err := scanBypass(row)
	if errors.Is(err, sql.ErrNoRows) {
		return BypassRule{}, errNotFound
	}
	return item, err
}

func (r *repository) createBypassRule(ctx context.Context, item BypassRule) (BypassRule, error) {
	now := time.Now().UTC().Format(rfc3339NanoUTC)
	item.CreatedAt = now
	item.UpdatedAt = now
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO model_bypass_rules
		(id, name, model_pattern, key_id, protocol, header_name, header_value,
		 action, target_model, priority, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.ModelPattern, item.KeyID, item.Protocol, item.HeaderName, item.HeaderValue,
		item.Action, item.TargetModel, item.Priority, boolToInt(item.Enabled), item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return BypassRule{}, fmt.Errorf("insert model_bypass_rule: %w", err)
	}
	return item, nil
}

func (r *repository) updateBypassRule(ctx context.Context, item BypassRule) (BypassRule, error) {
	item.UpdatedAt = time.Now().UTC().Format(rfc3339NanoUTC)
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE model_bypass_rules
		SET name = ?, model_pattern = ?, key_id = ?, protocol = ?, header_name = ?, header_value = ?,
		    action = ?, target_model = ?, priority = ?, enabled = ?, updated_at = ?
		WHERE id = ?`,
		item.Name, item.ModelPattern, item.KeyID, item.Protocol, item.HeaderName, item.HeaderValue,
		item.Action, item.TargetModel, item.Priority, boolToInt(item.Enabled), item.UpdatedAt, item.ID,
	)
	if err != nil {
		return BypassRule{}, fmt.Errorf("update model_bypass_rule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return BypassRule{}, errNotFound
	}
	return item, nil
}

func (r *repository) deleteBypassRule(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM model_bypass_rules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete model_bypass_rule: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanOverride(s scanner) (ModelOverride, error) {
	var (
		item       ModelOverride
		temp       sql.NullFloat64
		maxTok     sql.NullInt64
		topP       sql.NullFloat64
		extraRaw   string
		stripRaw   string
		enabledInt int
	)
	if err := s.Scan(
		&item.ID, &item.Name, &item.ModelPattern, &item.TargetModel,
		&item.SystemPromptPrefix, &item.SystemPromptSuffix,
		&temp, &maxTok, &topP, &extraRaw, &stripRaw,
		&item.Priority, &enabledInt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return ModelOverride{}, err
	}
	if temp.Valid {
		v := temp.Float64
		item.Temperature = &v
	}
	if maxTok.Valid {
		v := int(maxTok.Int64)
		item.MaxTokens = &v
	}
	if topP.Valid {
		v := topP.Float64
		item.TopP = &v
	}
	item.ExtraParams = make(map[string]any)
	if extraRaw != "" {
		_ = json.Unmarshal([]byte(extraRaw), &item.ExtraParams)
	}
	item.StripParams = make([]string, 0)
	if stripRaw != "" {
		_ = json.Unmarshal([]byte(stripRaw), &item.StripParams)
	}
	item.Enabled = enabledInt != 0
	return item, nil
}

func scanBypass(s scanner) (BypassRule, error) {
	var (
		item       BypassRule
		enabledInt int
	)
	if err := s.Scan(
		&item.ID, &item.Name, &item.ModelPattern, &item.KeyID, &item.Protocol,
		&item.HeaderName, &item.HeaderValue, &item.Action, &item.TargetModel,
		&item.Priority, &enabledInt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return BypassRule{}, err
	}
	item.Enabled = enabledInt != 0
	return item, nil
}

func marshalOverrideJSON(item ModelOverride) (extraJSON, stripJSON string, err error) {
	extra := item.ExtraParams
	if extra == nil {
		extra = map[string]any{}
	}
	strip := item.StripParams
	if strip == nil {
		strip = []string{}
	}
	eb, err := json.Marshal(extra)
	if err != nil {
		return "", "", fmt.Errorf("marshal extra_params: %w", err)
	}
	sb, err := json.Marshal(strip)
	if err != nil {
		return "", "", fmt.Errorf("marshal strip_params: %w", err)
	}
	return string(eb), string(sb), nil
}

func nullFloat(v *float64) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *v, Valid: true}
}

func nullInt(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
