package gateway

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	defaultChainTTL  = time.Hour
	rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"
)

type chainStore struct {
	db  *db.DB
	ttl time.Duration
}

func newChainStore(database *db.DB, ttl time.Duration) *chainStore {
	if ttl <= 0 {
		ttl = defaultChainTTL
	}
	return &chainStore{
		db:  database,
		ttl: ttl,
	}
}

func (s *chainStore) loadChain(ctx context.Context, responseID, keyID string, now time.Time) ([]CanonicalMessage, bool, error) {
	if s == nil || s.db == nil || responseID == "" || keyID == "" {
		return nil, false, nil
	}

	nowStr := now.UTC().Format(rfc3339NanoFixed)
	var raw string
	err := s.db.Reader().QueryRowContext(
		ctx,
		`SELECT messages_json FROM response_chains WHERE response_id = ? AND key_id = ? AND expires_at > ?`,
		responseID, keyID, nowStr,
	).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("gateway: load response chain: %w", err)
	}

	var msgs []CanonicalMessage
	if err := json.Unmarshal([]byte(raw), &msgs); err != nil {
		return nil, false, fmt.Errorf("gateway: decode response chain: %w", err)
	}
	return msgs, true, nil
}

func (s *chainStore) saveChain(ctx context.Context, responseID, keyID, model string, messages []CanonicalMessage, now time.Time) error {
	if s == nil || s.db == nil || responseID == "" || keyID == "" {
		return nil
	}

	encoded, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("gateway: encode response chain: %w", err)
	}

	utcNow := now.UTC()
	createdAt := utcNow.Format(rfc3339NanoFixed)
	expiresAt := utcNow.Add(s.ttl).Format(rfc3339NanoFixed)

	if _, err := s.db.Writer().ExecContext(ctx, `DELETE FROM response_chains WHERE expires_at <= ?`, createdAt); err != nil {
		return fmt.Errorf("gateway: prune expired response chains: %w", err)
	}

	if _, err := s.db.Writer().ExecContext(
		ctx,
		`INSERT OR REPLACE INTO response_chains (response_id, key_id, model, messages_json, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		responseID, keyID, model, string(encoded), expiresAt, createdAt,
	); err != nil {
		return fmt.Errorf("gateway: save response chain: %w", err)
	}
	return nil
}
