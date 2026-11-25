package keys

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

// APIKeySummary represents the public metadata for a gateway API key.
// It intentionally excludes key_hash and raw secret material.
type APIKeySummary struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	KeyPrefix     string   `json:"key_prefix"`
	AllowedModels []string `json:"allowed_models"`
	Scopes        []string `json:"scopes"`
	Tags          []string `json:"tags"`
	RateLimitRPM  int      `json:"rate_limit_rpm"`
	QuotaUSD      float64  `json:"quota_usd"`
	Disabled      bool     `json:"disabled"`
	ExpiresAt     *string  `json:"expires_at"`
	LastUsedAt    *string  `json:"last_used_at"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type createParams struct {
	ID            string
	Name          string
	KeyHash       string
	KeyPrefix     string
	AllowedModels []string
	Scopes        []string
	Tags          []string
	RateLimitRPM  int
	QuotaUSD      float64
	ExpiresAt     *string
}

type updateParams struct {
	Name          *string
	AllowedModels *[]string
	Scopes        *[]string
	Tags          *[]string
	RateLimitRPM  *int
	QuotaUSD      *float64
	Disabled      *bool
	ExpiresAtSet  bool
	ExpiresAt     *string
}

type updateResult struct {
	Summary APIKeySummary
	KeyHash string
}

type rotateResult struct {
	Summary APIKeySummary
	OldHash string
	NewHash string
}

type deleteResult struct {
	OldHash   string
	DeletedAt string
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

type rowScanner interface {
	Scan(dest ...any) error
}

const selectSummaryColumns = `id, name, key_prefix, allowed_models_json, scopes_json, tags_json,
	rate_limit_rpm, quota_usd, disabled, expires_at, last_used_at, created_at, updated_at`

func scanSummary(sc rowScanner) (APIKeySummary, error) {
	var (
		item              APIKeySummary
		allowedModelsJSON string
		scopesJSON        string
		tagsJSON          string
		disabledInt       int
		expiresAt         sql.NullString
		lastUsedAt        sql.NullString
	)
	if err := sc.Scan(
		&item.ID,
		&item.Name,
		&item.KeyPrefix,
		&allowedModelsJSON,
		&scopesJSON,
		&tagsJSON,
		&item.RateLimitRPM,
		&item.QuotaUSD,
		&disabledInt,
		&expiresAt,
		&lastUsedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return APIKeySummary{}, err
	}

	var err error
	if item.AllowedModels, err = decodeStringSlice(allowedModelsJSON); err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: decode allowed_models for %q: %w", item.ID, err)
	}
	if item.Scopes, err = decodeStringSlice(scopesJSON); err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: decode scopes for %q: %w", item.ID, err)
	}
	if item.Tags, err = decodeStringSlice(tagsJSON); err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: decode tags for %q: %w", item.ID, err)
	}
	item.Disabled = disabledInt != 0
	if expiresAt.Valid {
		s := expiresAt.String
		item.ExpiresAt = &s
	}
	if lastUsedAt.Valid {
		s := lastUsedAt.String
		item.LastUsedAt = &s
	}
	return item, nil
}

func (r *repository) List(ctx context.Context) ([]APIKeySummary, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("keys: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT `+selectSummaryColumns+` FROM api_keys ORDER BY created_at DESC, id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("keys: list api_keys: %w", err)
	}
	defer rows.Close()

	items := make([]APIKeySummary, 0)
	for rows.Next() {
		item, err := scanSummary(rows)
		if err != nil {
			return nil, fmt.Errorf("keys: scan api_key row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("keys: iterate api_keys: %w", err)
	}
	return items, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (APIKeySummary, error) {
	if r == nil || r.db == nil {
		return APIKeySummary{}, errors.New("keys: database is not initialized")
	}

	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT `+selectSummaryColumns+` FROM api_keys WHERE id = ?`,
		id,
	)
	item, err := scanSummary(row)
	if errors.Is(err, sql.ErrNoRows) {
		return APIKeySummary{}, ErrKeyNotFound
	}
	if err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: get api_key %q: %w", id, err)
	}
	return item, nil
}

func (r *repository) LookupByHash(ctx context.Context, keyHash string) (cachedEntry, error) {
	if r == nil || r.db == nil {
		return cachedEntry{}, errors.New("keys: database is not initialized")
	}

	var (
		info              contracts.KeyInfo
		allowedModelsJSON string
		scopesJSON        string
		tagsJSON          string
		disabledInt       int
		expiresAt         sql.NullString
	)
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT id, name, allowed_models_json, scopes_json, tags_json, rate_limit_rpm, quota_usd, disabled, expires_at
		 FROM api_keys
		 WHERE key_hash = ?`,
		keyHash,
	).Scan(
		&info.ID,
		&info.Name,
		&allowedModelsJSON,
		&scopesJSON,
		&tagsJSON,
		&info.RateLimitRPM,
		&info.QuotaUSD,
		&disabledInt,
		&expiresAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return cachedEntry{}, ErrKeyInvalid
	}
	if err != nil {
		return cachedEntry{}, fmt.Errorf("keys: lookup key by hash: %w", err)
	}

	if info.AllowedModels, err = decodeStringSlice(allowedModelsJSON); err != nil {
		return cachedEntry{}, fmt.Errorf("keys: decode allowed_models: %w", err)
	}
	if info.Scopes, err = decodeStringSlice(scopesJSON); err != nil {
		return cachedEntry{}, fmt.Errorf("keys: decode scopes: %w", err)
	}
	if info.Tags, err = decodeStringSlice(tagsJSON); err != nil {
		return cachedEntry{}, fmt.Errorf("keys: decode tags: %w", err)
	}
	info.Disabled = disabledInt != 0

	entry := cachedEntry{info: info}
	if expiresAt.Valid && expiresAt.String != "" {
		parsed, err := time.Parse(time.RFC3339Nano, expiresAt.String)
		if err != nil {
			return cachedEntry{}, fmt.Errorf("keys: parse stored expires_at: %w", err)
		}
		utc := parsed.UTC()
		entry.expiresAt = &utc
	}
	return entry, nil
}

func (r *repository) Create(ctx context.Context, p createParams) (APIKeySummary, error) {
	if r == nil || r.db == nil {
		return APIKeySummary{}, errors.New("keys: database is not initialized")
	}

	allowedJSON, err := encodeStringSlice(p.AllowedModels)
	if err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: encode allowed_models: %w", err)
	}
	scopesJSON, err := encodeStringSlice(p.Scopes)
	if err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: encode scopes: %w", err)
	}
	tagsJSON, err := encodeStringSlice(p.Tags)
	if err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: encode tags: %w", err)
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	var expVal any
	if p.ExpiresAt != nil {
		expVal = *p.ExpiresAt
	}

	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO api_keys (
			id, name, key_hash, key_prefix,
			allowed_models_json, scopes_json, tags_json,
			rate_limit_rpm, quota_usd, disabled,
			expires_at, last_used_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, NULL, ?, ?)`,
		p.ID,
		p.Name,
		p.KeyHash,
		p.KeyPrefix,
		allowedJSON,
		scopesJSON,
		tagsJSON,
		p.RateLimitRPM,
		p.QuotaUSD,
		expVal,
		now,
		now,
	); err != nil {
		return APIKeySummary{}, fmt.Errorf("keys: insert api_key: %w", err)
	}

	var expiresCopy *string
	if p.ExpiresAt != nil {
		s := *p.ExpiresAt
		expiresCopy = &s
	}

	return APIKeySummary{
		ID:            p.ID,
		Name:          p.Name,
		KeyPrefix:     p.KeyPrefix,
		AllowedModels: cloneStringSlice(p.AllowedModels),
		Scopes:        cloneStringSlice(p.Scopes),
		Tags:          cloneStringSlice(p.Tags),
		RateLimitRPM:  p.RateLimitRPM,
		QuotaUSD:      p.QuotaUSD,
		Disabled:      false,
		ExpiresAt:     expiresCopy,
		LastUsedAt:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (r *repository) Update(ctx context.Context, id string, p updateParams) (updateResult, error) {
	if r == nil || r.db == nil {
		return updateResult{}, errors.New("keys: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return updateResult{}, fmt.Errorf("keys: begin update transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var (
		keyHash           string
		current           APIKeySummary
		allowedModelsJSON string
		scopesJSON        string
		tagsJSON          string
		disabledInt       int
		expiresAt         sql.NullString
		lastUsedAt        sql.NullString
	)
	err = tx.QueryRowContext(
		ctx,
		`SELECT key_hash, `+selectSummaryColumns+` FROM api_keys WHERE id = ?`,
		id,
	).Scan(
		&keyHash,
		&current.ID,
		&current.Name,
		&current.KeyPrefix,
		&allowedModelsJSON,
		&scopesJSON,
		&tagsJSON,
		&current.RateLimitRPM,
		&current.QuotaUSD,
		&disabledInt,
		&expiresAt,
		&lastUsedAt,
		&current.CreatedAt,
		&current.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return updateResult{}, ErrKeyNotFound
	}
	if err != nil {
		return updateResult{}, fmt.Errorf("keys: read api_key %q for update: %w", id, err)
	}

	if current.AllowedModels, err = decodeStringSlice(allowedModelsJSON); err != nil {
		return updateResult{}, fmt.Errorf("keys: decode allowed_models: %w", err)
	}
	if current.Scopes, err = decodeStringSlice(scopesJSON); err != nil {
		return updateResult{}, fmt.Errorf("keys: decode scopes: %w", err)
	}
	if current.Tags, err = decodeStringSlice(tagsJSON); err != nil {
		return updateResult{}, fmt.Errorf("keys: decode tags: %w", err)
	}
	current.Disabled = disabledInt != 0
	if expiresAt.Valid {
		s := expiresAt.String
		current.ExpiresAt = &s
	}
	if lastUsedAt.Valid {
		s := lastUsedAt.String
		current.LastUsedAt = &s
	}

	if p.Name != nil {
		current.Name = *p.Name
	}
	if p.AllowedModels != nil {
		current.AllowedModels = cloneStringSlice(*p.AllowedModels)
	}
	if p.Scopes != nil {
		current.Scopes = cloneStringSlice(*p.Scopes)
	}
	if p.Tags != nil {
		current.Tags = cloneStringSlice(*p.Tags)
	}
	if p.RateLimitRPM != nil {
		current.RateLimitRPM = *p.RateLimitRPM
	}
	if p.QuotaUSD != nil {
		current.QuotaUSD = *p.QuotaUSD
	}
	if p.Disabled != nil {
		current.Disabled = *p.Disabled
	}
	if p.ExpiresAtSet {
		if p.ExpiresAt != nil {
			s := *p.ExpiresAt
			current.ExpiresAt = &s
		} else {
			current.ExpiresAt = nil
		}
	}

	newAllowedJSON, err := encodeStringSlice(current.AllowedModels)
	if err != nil {
		return updateResult{}, fmt.Errorf("keys: encode allowed_models: %w", err)
	}
	newScopesJSON, err := encodeStringSlice(current.Scopes)
	if err != nil {
		return updateResult{}, fmt.Errorf("keys: encode scopes: %w", err)
	}
	newTagsJSON, err := encodeStringSlice(current.Tags)
	if err != nil {
		return updateResult{}, fmt.Errorf("keys: encode tags: %w", err)
	}

	newDisabledInt := 0
	if current.Disabled {
		newDisabledInt = 1
	}
	var newExpiresVal any
	if current.ExpiresAt != nil {
		newExpiresVal = *current.ExpiresAt
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	current.UpdatedAt = now

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE api_keys
		 SET name = ?,
		     allowed_models_json = ?,
		     scopes_json = ?,
		     tags_json = ?,
		     rate_limit_rpm = ?,
		     quota_usd = ?,
		     disabled = ?,
		     expires_at = ?,
		     updated_at = ?
		 WHERE id = ?`,
		current.Name,
		newAllowedJSON,
		newScopesJSON,
		newTagsJSON,
		current.RateLimitRPM,
		current.QuotaUSD,
		newDisabledInt,
		newExpiresVal,
		now,
		id,
	); err != nil {
		return updateResult{}, fmt.Errorf("keys: update api_key %q: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return updateResult{}, fmt.Errorf("keys: commit update transaction: %w", err)
	}

	return updateResult{
		Summary: current,
		KeyHash: keyHash,
	}, nil
}

func (r *repository) Rotate(ctx context.Context, id, newHash, newPrefix string) (rotateResult, error) {
	if r == nil || r.db == nil {
		return rotateResult{}, errors.New("keys: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return rotateResult{}, fmt.Errorf("keys: begin rotate transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var oldHash string
	err = tx.QueryRowContext(ctx, `SELECT key_hash FROM api_keys WHERE id = ?`, id).Scan(&oldHash)
	if errors.Is(err, sql.ErrNoRows) {
		return rotateResult{}, ErrKeyNotFound
	}
	if err != nil {
		return rotateResult{}, fmt.Errorf("keys: read api_key %q for rotate: %w", id, err)
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE api_keys SET key_hash = ?, key_prefix = ?, updated_at = ? WHERE id = ?`,
		newHash,
		newPrefix,
		now,
		id,
	); err != nil {
		return rotateResult{}, fmt.Errorf("keys: rotate api_key %q: %w", id, err)
	}

	summary, err := scanSummary(tx.QueryRowContext(
		ctx,
		`SELECT `+selectSummaryColumns+` FROM api_keys WHERE id = ?`,
		id,
	))
	if err != nil {
		return rotateResult{}, fmt.Errorf("keys: reload rotated api_key %q: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return rotateResult{}, fmt.Errorf("keys: commit rotate transaction: %w", err)
	}

	return rotateResult{
		Summary: summary,
		OldHash: oldHash,
		NewHash: newHash,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id string) (deleteResult, error) {
	if r == nil || r.db == nil {
		return deleteResult{}, errors.New("keys: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return deleteResult{}, fmt.Errorf("keys: begin delete transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var oldHash string
	err = tx.QueryRowContext(ctx, `SELECT key_hash FROM api_keys WHERE id = ?`, id).Scan(&oldHash)
	if errors.Is(err, sql.ErrNoRows) {
		return deleteResult{}, ErrKeyNotFound
	}
	if err != nil {
		return deleteResult{}, fmt.Errorf("keys: read api_key %q for delete: %w", id, err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM api_keys WHERE id = ?`, id); err != nil {
		return deleteResult{}, fmt.Errorf("keys: delete api_key %q: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return deleteResult{}, fmt.Errorf("keys: commit delete transaction: %w", err)
	}

	return deleteResult{
		OldHash:   oldHash,
		DeletedAt: time.Now().UTC().Format(rfc3339NanoFixed),
	}, nil
}

func (r *repository) TouchLastUsed(ctx context.Context, id, usedAt string) error {
	if r == nil || r.db == nil {
		return errors.New("keys: database is not initialized")
	}
	_, err := r.db.Writer().ExecContext(ctx, `UPDATE api_keys SET last_used_at = ? WHERE id = ?`, usedAt, id)
	if err != nil {
		return fmt.Errorf("keys: touch last_used_at for %q: %w", id, err)
	}
	return nil
}

func (r *repository) SumUsageCost(ctx context.Context, keyID string) (float64, bool) {
	if r == nil || r.db == nil || keyID == "" {
		return 0, false
	}
	var total float64
	if err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(cost_usd), 0.0) FROM usage_logs WHERE key_id = ?`,
		keyID,
	).Scan(&total); err != nil {
		return 0, false
	}
	return total, true
}

func encodeStringSlice(items []string) (string, error) {
	normalized := cloneStringSlice(items)
	data, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeStringSlice(raw string) ([]string, error) {
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil, err
	}
	if items == nil {
		return []string{}, nil
	}
	return items, nil
}
