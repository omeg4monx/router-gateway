package media

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const tsLayout = "2006-01-02T15:04:05.000000000Z"

var (
	errNotFound      = errors.New("media resource not found")
	errAlreadyExists = errors.New("media resource name already exists")
)

// ProviderSummary is the API-safe representation of a media provider (never
// exposing plaintext or encrypted API keys).
type ProviderSummary struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Capability   string            `json:"capability"`
	Driver       string            `json:"driver"`
	BaseURL      string            `json:"base_url"`
	HasAPIKey    bool              `json:"has_api_key"`
	DefaultModel string            `json:"default_model"`
	ExtraHeaders map[string]string `json:"extra_headers"`
	Priority     int               `json:"priority"`
	Enabled      bool              `json:"enabled"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

type providerRecord struct {
	ProviderSummary
	APIKeyEnc string
}

// Combo represents an ordered media provider fallback chain for a capability.
type Combo struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Capability string   `json:"capability"`
	Targets    []string `json:"targets"`
	Enabled    bool     `json:"enabled"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func newID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

func nowUTC() string {
	return time.Now().UTC().Format(tsLayout)
}

func isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") || strings.Contains(msg, "constraint failed: unique")
}

func (r *repo) listProviders(ctx context.Context) ([]providerRecord, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, capability, driver, base_url, api_key_enc, default_model,
		       extra_headers_json, priority, enabled, created_at, updated_at
		FROM media_providers
		ORDER BY priority ASC, created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list media providers: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]providerRecord, 0)
	for rows.Next() {
		rec, err := scanProvider(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

func (r *repo) getProvider(ctx context.Context, id string) (providerRecord, error) {
	row := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, capability, driver, base_url, api_key_enc, default_model,
		       extra_headers_json, priority, enabled, created_at, updated_at
		FROM media_providers
		WHERE id = ?
	`, id)
	rec, err := scanProvider(row)
	if errors.Is(err, sql.ErrNoRows) {
		return providerRecord{}, errNotFound
	}
	return rec, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProvider(s rowScanner) (providerRecord, error) {
	var (
		rec        providerRecord
		headersRaw string
		enabledInt int
	)
	if err := s.Scan(
		&rec.ID,
		&rec.Name,
		&rec.Capability,
		&rec.Driver,
		&rec.BaseURL,
		&rec.APIKeyEnc,
		&rec.DefaultModel,
		&headersRaw,
		&rec.Priority,
		&enabledInt,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	); err != nil {
		return providerRecord{}, err
	}
	rec.Enabled = enabledInt != 0
	rec.HasAPIKey = strings.TrimSpace(rec.APIKeyEnc) != ""
	rec.ExtraHeaders = make(map[string]string)
	if strings.TrimSpace(headersRaw) != "" {
		_ = json.Unmarshal([]byte(headersRaw), &rec.ExtraHeaders)
	}
	return rec, nil
}

func (r *repo) createProvider(ctx context.Context, rec providerRecord) (providerRecord, error) {
	if rec.ID == "" {
		rec.ID = newID("mp")
	}
	now := nowUTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	if rec.ExtraHeaders == nil {
		rec.ExtraHeaders = make(map[string]string)
	}
	headersJSON, err := json.Marshal(rec.ExtraHeaders)
	if err != nil {
		return providerRecord{}, fmt.Errorf("marshal extra headers: %w", err)
	}
	enabledInt := 0
	if rec.Enabled {
		enabledInt = 1
	}
	_, err = r.db.Writer().ExecContext(ctx, `
		INSERT INTO media_providers (
			id, name, capability, driver, base_url, api_key_enc, default_model,
			extra_headers_json, priority, enabled, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, rec.ID, rec.Name, rec.Capability, rec.Driver, rec.BaseURL, rec.APIKeyEnc,
		rec.DefaultModel, string(headersJSON), rec.Priority, enabledInt, rec.CreatedAt, rec.UpdatedAt)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return providerRecord{}, errAlreadyExists
		}
		return providerRecord{}, fmt.Errorf("insert media provider: %w", err)
	}
	rec.HasAPIKey = strings.TrimSpace(rec.APIKeyEnc) != ""
	return rec, nil
}

func (r *repo) updateProvider(ctx context.Context, rec providerRecord) (providerRecord, error) {
	rec.UpdatedAt = nowUTC()
	if rec.ExtraHeaders == nil {
		rec.ExtraHeaders = make(map[string]string)
	}
	headersJSON, err := json.Marshal(rec.ExtraHeaders)
	if err != nil {
		return providerRecord{}, fmt.Errorf("marshal extra headers: %w", err)
	}
	enabledInt := 0
	if rec.Enabled {
		enabledInt = 1
	}
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE media_providers
		SET name = ?, capability = ?, driver = ?, base_url = ?, api_key_enc = ?,
		    default_model = ?, extra_headers_json = ?, priority = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`, rec.Name, rec.Capability, rec.Driver, rec.BaseURL, rec.APIKeyEnc,
		rec.DefaultModel, string(headersJSON), rec.Priority, enabledInt, rec.UpdatedAt, rec.ID)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return providerRecord{}, errAlreadyExists
		}
		return providerRecord{}, fmt.Errorf("update media provider: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return providerRecord{}, errNotFound
	}
	rec.HasAPIKey = strings.TrimSpace(rec.APIKeyEnc) != ""
	return rec, nil
}

func (r *repo) deleteProvider(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM media_providers WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete media provider: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}

func (r *repo) listCombos(ctx context.Context) ([]Combo, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, capability, targets_json, enabled, created_at, updated_at
		FROM media_combos
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list media combos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]Combo, 0)
	for rows.Next() {
		c, err := scanCombo(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *repo) getCombo(ctx context.Context, id string) (Combo, error) {
	row := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, capability, targets_json, enabled, created_at, updated_at
		FROM media_combos
		WHERE id = ?
	`, id)
	c, err := scanCombo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Combo{}, errNotFound
	}
	return c, err
}

func scanCombo(s rowScanner) (Combo, error) {
	var (
		c          Combo
		targetsRaw string
		enabledInt int
	)
	if err := s.Scan(&c.ID, &c.Name, &c.Capability, &targetsRaw, &enabledInt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return Combo{}, err
	}
	c.Enabled = enabledInt != 0
	c.Targets = make([]string, 0)
	if strings.TrimSpace(targetsRaw) != "" {
		_ = json.Unmarshal([]byte(targetsRaw), &c.Targets)
	}
	return c, nil
}

func (r *repo) createCombo(ctx context.Context, c Combo) (Combo, error) {
	if c.ID == "" {
		c.ID = newID("mc")
	}
	now := nowUTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	targetsJSON, err := json.Marshal(c.Targets)
	if err != nil {
		return Combo{}, fmt.Errorf("marshal targets: %w", err)
	}
	enabledInt := 0
	if c.Enabled {
		enabledInt = 1
	}
	_, err = r.db.Writer().ExecContext(ctx, `
		INSERT INTO media_combos (id, name, capability, targets_json, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, c.ID, c.Name, c.Capability, string(targetsJSON), enabledInt, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return Combo{}, errAlreadyExists
		}
		return Combo{}, fmt.Errorf("insert media combo: %w", err)
	}
	return c, nil
}

func (r *repo) updateCombo(ctx context.Context, c Combo) (Combo, error) {
	c.UpdatedAt = nowUTC()
	targetsJSON, err := json.Marshal(c.Targets)
	if err != nil {
		return Combo{}, fmt.Errorf("marshal targets: %w", err)
	}
	enabledInt := 0
	if c.Enabled {
		enabledInt = 1
	}
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE media_combos
		SET name = ?, capability = ?, targets_json = ?, enabled = ?, updated_at = ?
		WHERE id = ?
	`, c.Name, c.Capability, string(targetsJSON), enabledInt, c.UpdatedAt, c.ID)
	if err != nil {
		if isUniqueConstraintErr(err) {
			return Combo{}, errAlreadyExists
		}
		return Combo{}, fmt.Errorf("update media combo: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Combo{}, errNotFound
	}
	return c, nil
}

func (r *repo) deleteCombo(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM media_combos WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete media combo: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}
