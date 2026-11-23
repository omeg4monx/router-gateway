package tunnel

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	tsLayout         = "2006-01-02T15:04:05.000000000Z"
	maxStoredEvents  = 200
	defaultEventList = 50
)

// ConfigRecord represents persisted tunnel settings.
type ConfigRecord struct {
	Provider      string `json:"provider"`
	LocalURL      string `json:"local_url"`
	TailscaleMode string `json:"tailscale_mode"`
	AutoStart     bool   `json:"auto_start"`
	UpdatedAt     string `json:"updated_at"`
}

// Event represents a persisted tunnel lifecycle audit event.
type Event struct {
	ID        string `json:"id"`
	Provider  string `json:"provider"`
	EventType string `json:"event_type"`
	Message   string `json:"message"`
	PublicURL string `json:"public_url"`
	CreatedAt string `json:"created_at"`
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func nowUTC() string {
	return time.Now().UTC().Format(tsLayout)
}

func newID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

func (r *repo) loadOrInitConfig(ctx context.Context, defaultLocalURL string) (ConfigRecord, error) {
	row := r.db.Reader().QueryRowContext(ctx, `
		SELECT provider, local_url, tailscale_mode, auto_start, updated_at
		FROM tunnel_config
		WHERE id = 1
	`)
	var (
		cfg       ConfigRecord
		autoStart int
	)
	err := row.Scan(&cfg.Provider, &cfg.LocalURL, &cfg.TailscaleMode, &autoStart, &cfg.UpdatedAt)
	if err == nil {
		cfg.AutoStart = autoStart != 0
		if cfg.LocalURL == "" {
			cfg.LocalURL = defaultLocalURL
		}
		return cfg, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return ConfigRecord{}, fmt.Errorf("load tunnel config: %w", err)
	}

	cfg = ConfigRecord{
		Provider:      ProviderCloudflare,
		LocalURL:      defaultLocalURL,
		TailscaleMode: TailscaleModeFunnel,
		AutoStart:     false,
		UpdatedAt:     nowUTC(),
	}
	if err := r.saveConfig(ctx, cfg); err != nil {
		return ConfigRecord{}, err
	}
	return cfg, nil
}

func (r *repo) saveConfig(ctx context.Context, cfg ConfigRecord) error {
	if cfg.UpdatedAt == "" {
		cfg.UpdatedAt = nowUTC()
	}
	autoInt := 0
	if cfg.AutoStart {
		autoInt = 1
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO tunnel_config (id, provider, local_url, tailscale_mode, auto_start, updated_at)
		VALUES (1, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			provider = excluded.provider,
			local_url = excluded.local_url,
			tailscale_mode = excluded.tailscale_mode,
			auto_start = excluded.auto_start,
			updated_at = excluded.updated_at
	`, cfg.Provider, cfg.LocalURL, cfg.TailscaleMode, autoInt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save tunnel config: %w", err)
	}
	return nil
}

func (r *repo) appendEvent(ctx context.Context, provider, eventType, message, publicURL string) error {
	ev := Event{
		ID:        newID("tev"),
		Provider:  provider,
		EventType: eventType,
		Message:   message,
		PublicURL: publicURL,
		CreatedAt: nowUTC(),
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO tunnel_events (id, provider, event_type, message, public_url, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, ev.ID, ev.Provider, ev.EventType, ev.Message, ev.PublicURL, ev.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert tunnel event: %w", err)
	}
	_, _ = r.db.Writer().ExecContext(ctx, `
		DELETE FROM tunnel_events
		WHERE id NOT IN (
			SELECT id FROM tunnel_events ORDER BY created_at DESC LIMIT ?
		)
	`, maxStoredEvents)
	return nil
}

func (r *repo) listEvents(ctx context.Context, limit int) ([]Event, error) {
	if limit <= 0 {
		limit = defaultEventList
	}
	if limit > maxStoredEvents {
		limit = maxStoredEvents
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, provider, event_type, message, public_url, created_at
		FROM tunnel_events
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list tunnel events: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]Event, 0)
	for rows.Next() {
		var ev Event
		if err := rows.Scan(&ev.ID, &ev.Provider, &ev.EventType, &ev.Message, &ev.PublicURL, &ev.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return out, rows.Err()
}
