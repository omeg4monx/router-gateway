package telegram

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
	timestampFormat = "2006-01-02T15:04:05.000000000Z"
	maxStoredEvents = 200
)

type storedConfig struct {
	Enabled            bool
	BotTokenEnc        string
	BotTokenMask       string
	ChatID             string
	UserID             string
	WebhookURL         string
	WebhookSecretEnc   string
	NotifyOnJobFailure bool
	NotifyOnAuthChange bool
	NotifyOnEmit       bool
	UpdatedAt          string
}

// Event represents a recorded Telegram webhook, command, test, or alert event.
type Event struct {
	ID        string `json:"id"`
	Direction string `json:"direction"`
	Kind      string `json:"kind"`
	Command   string `json:"command"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Error     string `json:"error"`
	CreatedAt string `json:"created_at"`
}

type repo struct {
	db *db.DB
}

func newRepo(database *db.DB) *repo {
	return &repo{db: database}
}

func nowTimestamp() string {
	return time.Now().UTC().Format(timestampFormat)
}

func randomHexID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

func defaultStoredConfig() storedConfig {
	return storedConfig{
		NotifyOnJobFailure: true,
		NotifyOnAuthChange: true,
		UpdatedAt:          nowTimestamp(),
	}
}

func (r *repo) getConfig(ctx context.Context) (storedConfig, error) {
	if r == nil || r.db == nil {
		return defaultStoredConfig(), nil
	}
	var (
		cfg           storedConfig
		enabled       int
		notifyJobFail int
		notifyAuthChg int
		notifyEmit    int
	)
	err := r.db.Reader().QueryRowContext(ctx, `
		SELECT enabled, bot_token_enc, bot_token_mask, chat_id, user_id,
		       webhook_url, webhook_secret_enc, notify_on_job_failure,
		       notify_on_auth_change, notify_on_emit, updated_at
		FROM telegram_config
		WHERE id = 1
	`).Scan(
		&enabled, &cfg.BotTokenEnc, &cfg.BotTokenMask, &cfg.ChatID, &cfg.UserID,
		&cfg.WebhookURL, &cfg.WebhookSecretEnc, &notifyJobFail,
		&notifyAuthChg, &notifyEmit, &cfg.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return defaultStoredConfig(), nil
	}
	if err != nil {
		return storedConfig{}, fmt.Errorf("query telegram config: %w", err)
	}
	cfg.Enabled = enabled != 0
	cfg.NotifyOnJobFailure = notifyJobFail != 0
	cfg.NotifyOnAuthChange = notifyAuthChg != 0
	cfg.NotifyOnEmit = notifyEmit != 0
	return cfg, nil
}

func (r *repo) saveConfig(ctx context.Context, cfg storedConfig) error {
	if r == nil || r.db == nil {
		return nil
	}
	cfg.UpdatedAt = nowTimestamp()
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO telegram_config (
			id, enabled, bot_token_enc, bot_token_mask, chat_id, user_id,
			webhook_url, webhook_secret_enc, notify_on_job_failure,
			notify_on_auth_change, notify_on_emit, updated_at
		) VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			enabled = excluded.enabled,
			bot_token_enc = excluded.bot_token_enc,
			bot_token_mask = excluded.bot_token_mask,
			chat_id = excluded.chat_id,
			user_id = excluded.user_id,
			webhook_url = excluded.webhook_url,
			webhook_secret_enc = excluded.webhook_secret_enc,
			notify_on_job_failure = excluded.notify_on_job_failure,
			notify_on_auth_change = excluded.notify_on_auth_change,
			notify_on_emit = excluded.notify_on_emit,
			updated_at = excluded.updated_at
	`,
		boolToInt(cfg.Enabled),
		cfg.BotTokenEnc,
		cfg.BotTokenMask,
		cfg.ChatID,
		cfg.UserID,
		cfg.WebhookURL,
		cfg.WebhookSecretEnc,
		boolToInt(cfg.NotifyOnJobFailure),
		boolToInt(cfg.NotifyOnAuthChange),
		boolToInt(cfg.NotifyOnEmit),
		cfg.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save telegram config: %w", err)
	}
	return nil
}

func (r *repo) recordEvent(ctx context.Context, ev Event) error {
	if r == nil || r.db == nil {
		return nil
	}
	if ev.ID == "" {
		ev.ID = randomHexID("tgev")
	}
	if ev.CreatedAt == "" {
		ev.CreatedAt = nowTimestamp()
	}
	_, err := r.db.Writer().ExecContext(ctx, `
		INSERT INTO telegram_events (id, direction, kind, command, status, summary, error, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, ev.ID, ev.Direction, ev.Kind, ev.Command, ev.Status, ev.Summary, ev.Error, ev.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert telegram event: %w", err)
	}
	_, _ = r.db.Writer().ExecContext(ctx, `
		DELETE FROM telegram_events
		WHERE id NOT IN (
			SELECT id FROM telegram_events ORDER BY created_at DESC LIMIT ?
		)
	`, maxStoredEvents)
	return nil
}

func (r *repo) listEvents(ctx context.Context, limit int) ([]Event, error) {
	if r == nil || r.db == nil {
		return []Event{}, nil
	}
	if limit <= 0 || limit > maxStoredEvents {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(ctx, `
		SELECT id, direction, kind, command, status, summary, error, created_at
		FROM telegram_events
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list telegram events: %w", err)
	}
	defer rows.Close()

	out := make([]Event, 0, limit)
	for rows.Next() {
		var ev Event
		if err := rows.Scan(
			&ev.ID, &ev.Direction, &ev.Kind, &ev.Command,
			&ev.Status, &ev.Summary, &ev.Error, &ev.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan telegram event: %w", err)
		}
		out = append(out, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate telegram events: %w", err)
	}
	return out, nil
}

func (r *repo) clearEvents(ctx context.Context) (int64, error) {
	if r == nil || r.db == nil {
		return 0, nil
	}
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM telegram_events`)
	if err != nil {
		return 0, fmt.Errorf("clear telegram events: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
