package proxypools

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoUTC = "2006-01-02T15:04:05.000000000Z"

var errNotFound = errors.New("proxypools: not found")

// PoolSummary is the public, credential-redacted view of a configured proxy pool.
type PoolSummary struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Scheme         string `json:"scheme"`
	DisplayURL     string `json:"display_url"`
	HasCredentials bool   `json:"has_credentials"`
	Enabled        bool   `json:"enabled"`
	Status         string `json:"status"`
	LastLatencyMS  int64  `json:"last_latency_ms"`
	LastError      string `json:"last_error"`
	LastTestedAt   string `json:"last_tested_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// BindingSummary maps a provider identifier to a proxy pool ID.
type BindingSummary struct {
	ProviderID string `json:"provider_id"`
	PoolID     string `json:"pool_id"`
	UpdatedAt  string `json:"updated_at"`
}

type storedPool struct {
	PoolSummary
	URLEncrypted string
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) listPools(ctx context.Context) ([]storedPool, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, name, scheme, url_encrypted, display_url, has_credentials, enabled,
		       status, last_latency_ms, last_error, last_tested_at, created_at, updated_at
		FROM proxy_pools
		ORDER BY created_at ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("proxypools: list pools: %w", err)
	}
	defer rows.Close()

	out := make([]storedPool, 0)
	for rows.Next() {
		var item storedPool
		var hasCreds, enabled int
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Scheme, &item.URLEncrypted, &item.DisplayURL,
			&hasCreds, &enabled, &item.Status, &item.LastLatencyMS, &item.LastError,
			&item.LastTestedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("proxypools: scan pool: %w", err)
		}
		item.HasCredentials = hasCreds != 0
		item.Enabled = enabled != 0
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("proxypools: iterate pools: %w", err)
	}
	return out, nil
}

func (r *repository) getPool(ctx context.Context, id string) (storedPool, error) {
	var item storedPool
	var hasCreds, enabled int
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT id, name, scheme, url_encrypted, display_url, has_credentials, enabled,
		       status, last_latency_ms, last_error, last_tested_at, created_at, updated_at
		FROM proxy_pools
		WHERE id = ?`, id).Scan(
		&item.ID, &item.Name, &item.Scheme, &item.URLEncrypted, &item.DisplayURL,
		&hasCreds, &enabled, &item.Status, &item.LastLatencyMS, &item.LastError,
		&item.LastTestedAt, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return storedPool{}, errNotFound
	}
	if err != nil {
		return storedPool{}, fmt.Errorf("proxypools: get pool: %w", err)
	}
	item.HasCredentials = hasCreds != 0
	item.Enabled = enabled != 0
	return item, nil
}

func (r *repository) createPool(ctx context.Context, p storedPool) error {
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO proxy_pools (
			id, name, scheme, url_encrypted, display_url, has_credentials, enabled,
			status, last_latency_ms, last_error, last_tested_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Scheme, p.URLEncrypted, p.DisplayURL,
		boolToInt(p.HasCredentials), boolToInt(p.Enabled),
		p.Status, p.LastLatencyMS, p.LastError, p.LastTestedAt, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("proxypools: insert pool: %w", err)
	}
	return nil
}

func (r *repository) updatePool(ctx context.Context, p storedPool) error {
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE proxy_pools
		SET name = ?, scheme = ?, url_encrypted = ?, display_url = ?, has_credentials = ?,
		    enabled = ?, updated_at = ?
		WHERE id = ?`,
		p.Name, p.Scheme, p.URLEncrypted, p.DisplayURL,
		boolToInt(p.HasCredentials), boolToInt(p.Enabled), p.UpdatedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("proxypools: update pool: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}

func (r *repository) updatePoolHealth(ctx context.Context, id, status string, latencyMS int64, lastErr, testedAt string) error {
	res, err := r.db.Writer().ExecContext(ctx, `
		UPDATE proxy_pools
		SET status = ?, last_latency_ms = ?, last_error = ?, last_tested_at = ?, updated_at = ?
		WHERE id = ?`,
		status, latencyMS, lastErr, testedAt, testedAt, id,
	)
	if err != nil {
		return fmt.Errorf("proxypools: update pool health: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}

func (r *repository) deletePool(ctx context.Context, id string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM proxy_pools WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("proxypools: delete pool: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}

func (r *repository) listBindings(ctx context.Context) ([]BindingSummary, error) {
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT provider_id, pool_id, updated_at
		FROM proxy_bindings
		ORDER BY provider_id ASC`)
	if err != nil {
		return nil, fmt.Errorf("proxypools: list bindings: %w", err)
	}
	defer rows.Close()

	out := make([]BindingSummary, 0)
	for rows.Next() {
		var b BindingSummary
		if err := rows.Scan(&b.ProviderID, &b.PoolID, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("proxypools: scan binding: %w", err)
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("proxypools: iterate bindings: %w", err)
	}
	return out, nil
}

func (r *repository) upsertBinding(ctx context.Context, providerID, poolID, now string) error {
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO proxy_bindings (provider_id, pool_id, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(provider_id) DO UPDATE SET pool_id = excluded.pool_id, updated_at = excluded.updated_at`,
		providerID, poolID, now,
	)
	if err != nil {
		return fmt.Errorf("proxypools: upsert binding: %w", err)
	}
	return nil
}

func (r *repository) deleteBinding(ctx context.Context, providerID string) error {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM proxy_bindings WHERE provider_id = ?`, providerID)
	if err != nil {
		return fmt.Errorf("proxypools: delete binding: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errNotFound
	}
	return nil
}

type parsedProxyURL struct {
	Raw            string
	Scheme         string
	DisplayURL     string
	HasCredentials bool
	Parsed         *url.URL
}

func parseAndValidateProxyURL(raw string) (parsedProxyURL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return parsedProxyURL{}, errors.New("proxy url must not be empty")
	}
	u, err := url.Parse(trimmed)
	if err != nil {
		return parsedProxyURL{}, errors.New("malformed proxy url")
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "http", "https", "socks5", "socks5h":
	default:
		return parsedProxyURL{}, fmt.Errorf("unsupported proxy scheme %q (allowed: http, https, socks5, socks5h)", u.Scheme)
	}
	if u.Hostname() == "" {
		return parsedProxyURL{}, errors.New("proxy host must not be empty")
	}

	hasCreds := u.User != nil
	display := scheme + "://"
	if hasCreds {
		display += "***@"
	}
	display += u.Host
	if u.Path != "" && u.Path != "/" {
		display += u.Path
	}

	return parsedProxyURL{
		Raw:            trimmed,
		Scheme:         scheme,
		DisplayURL:     display,
		HasCredentials: hasCreds,
		Parsed:         u,
	}, nil
}

func formatUTC(t time.Time) string {
	return t.UTC().Format(rfc3339NanoUTC)
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
