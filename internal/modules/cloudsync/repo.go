package cloudsync

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

	policyMerge      = "merge"
	policyRemoteWins = "remote_wins"
	policyLocalWins  = "local_wins"
)

// Config represents the persisted cloud sync configuration without exposing plaintext tokens.
type Config struct {
	Enabled         bool   `json:"enabled"`
	EndpointURL     string `json:"endpoint_url"`
	HasAuthToken    bool   `json:"has_auth_token"`
	AuthTokenMask   string `json:"auth_token_mask"`
	SyncIntervalMin int    `json:"sync_interval_min"`
	ConflictPolicy  string `json:"conflict_policy"`
	IncludeSettings bool   `json:"include_settings"`
	IncludeCombos   bool   `json:"include_combos"`
	LastSyncAt      string `json:"last_sync_at"`
	LastStatus      string `json:"last_status"`
	LastError       string `json:"last_error"`
	UpdatedAt       string `json:"updated_at"`
}

type storedConfig struct {
	Config
	TokenEnc string
}

// HistoryEvent records a single push, pull, export, or restore operation.
type HistoryEvent struct {
	ID        string `json:"id"`
	Direction string `json:"direction"`
	Status    string `json:"status"`
	Summary   string `json:"summary"`
	Checksum  string `json:"checksum"`
	ItemCount int    `json:"item_count"`
	Error     string `json:"error"`
	CreatedAt string `json:"created_at"`
}

type repository struct {
	db  *db.DB
	now func() time.Time
}

func newRepository(database *db.DB, clock func() time.Time) *repository {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	return &repository{db: database, now: clock}
}

func (r *repository) nowUTC() string {
	return r.now().UTC().Format(rfc3339NanoFixed)
}

func (r *repository) GetConfig(ctx context.Context) (storedConfig, error) {
	row := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT enabled, endpoint_url, token_enc, token_mask, sync_interval_min, conflict_policy,
		        include_settings, include_combos, last_sync_at, last_status, last_error, updated_at
		 FROM cloudsync_config WHERE id = 1`,
	)
	var (
		sc          storedConfig
		enabledInt  int
		incSettings int
		incCombos   int
	)
	err := row.Scan(
		&enabledInt,
		&sc.EndpointURL,
		&sc.TokenEnc,
		&sc.AuthTokenMask,
		&sc.SyncIntervalMin,
		&sc.ConflictPolicy,
		&incSettings,
		&incCombos,
		&sc.LastSyncAt,
		&sc.LastStatus,
		&sc.LastError,
		&sc.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return storedConfig{
			Config: Config{
				SyncIntervalMin: 30,
				ConflictPolicy:  policyMerge,
				IncludeSettings: true,
				IncludeCombos:   true,
			},
		}, nil
	}
	if err != nil {
		return storedConfig{}, fmt.Errorf("cloudsync: query config: %w", err)
	}
	sc.Enabled = enabledInt != 0
	sc.HasAuthToken = strings.TrimSpace(sc.TokenEnc) != ""
	sc.IncludeSettings = incSettings != 0
	sc.IncludeCombos = incCombos != 0
	return sc, nil
}

func (r *repository) SaveConfig(ctx context.Context, sc storedConfig) (storedConfig, error) {
	sc.UpdatedAt = r.nowUTC()
	sc.HasAuthToken = strings.TrimSpace(sc.TokenEnc) != ""
	_, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO cloudsync_config
		 (id, enabled, endpoint_url, token_enc, token_mask, sync_interval_min, conflict_policy,
		  include_settings, include_combos, last_sync_at, last_status, last_error, updated_at)
		 VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   enabled = excluded.enabled,
		   endpoint_url = excluded.endpoint_url,
		   token_enc = excluded.token_enc,
		   token_mask = excluded.token_mask,
		   sync_interval_min = excluded.sync_interval_min,
		   conflict_policy = excluded.conflict_policy,
		   include_settings = excluded.include_settings,
		   include_combos = excluded.include_combos,
		   last_sync_at = excluded.last_sync_at,
		   last_status = excluded.last_status,
		   last_error = excluded.last_error,
		   updated_at = excluded.updated_at`,
		boolToInt(sc.Enabled),
		sc.EndpointURL,
		sc.TokenEnc,
		sc.AuthTokenMask,
		sc.SyncIntervalMin,
		sc.ConflictPolicy,
		boolToInt(sc.IncludeSettings),
		boolToInt(sc.IncludeCombos),
		sc.LastSyncAt,
		sc.LastStatus,
		sc.LastError,
		sc.UpdatedAt,
	)
	if err != nil {
		return storedConfig{}, fmt.Errorf("cloudsync: save config: %w", err)
	}
	return sc, nil
}

func (r *repository) RecordSyncStatus(ctx context.Context, status, errMsg string) error {
	now := r.nowUTC()
	sc, err := r.GetConfig(ctx)
	if err != nil {
		return err
	}
	sc.LastSyncAt = now
	sc.LastStatus = status
	sc.LastError = errMsg
	_, err = r.SaveConfig(ctx, sc)
	return err
}

func (r *repository) InsertHistory(ctx context.Context, direction, status, summary, checksum string, itemCount int, errMsg string) (HistoryEvent, error) {
	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return HistoryEvent{}, fmt.Errorf("cloudsync: generate history id: %w", err)
	}
	ev := HistoryEvent{
		ID:        "sync_" + hex.EncodeToString(raw),
		Direction: direction,
		Status:    status,
		Summary:   summary,
		Checksum:  checksum,
		ItemCount: itemCount,
		Error:     errMsg,
		CreatedAt: r.nowUTC(),
	}
	_, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO cloudsync_history (id, direction, status, summary, checksum, item_count, error, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		ev.ID, ev.Direction, ev.Status, ev.Summary, ev.Checksum, ev.ItemCount, ev.Error, ev.CreatedAt,
	)
	if err != nil {
		return HistoryEvent{}, fmt.Errorf("cloudsync: insert history: %w", err)
	}
	return ev, nil
}

func (r *repository) ListHistory(ctx context.Context, limit int) ([]HistoryEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT id, direction, status, summary, checksum, item_count, error, created_at
		 FROM cloudsync_history
		 ORDER BY created_at DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("cloudsync: query history: %w", err)
	}
	defer rows.Close()

	events := make([]HistoryEvent, 0)
	for rows.Next() {
		var ev HistoryEvent
		if err := rows.Scan(&ev.ID, &ev.Direction, &ev.Status, &ev.Summary, &ev.Checksum, &ev.ItemCount, &ev.Error, &ev.CreatedAt); err != nil {
			return nil, fmt.Errorf("cloudsync: scan history: %w", err)
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cloudsync: iterate history: %w", err)
	}
	return events, nil
}

func (r *repository) ClearHistory(ctx context.Context) (int64, error) {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM cloudsync_history`)
	if err != nil {
		return 0, fmt.Errorf("cloudsync: clear history: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func maskSecret(secret string) string {
	s := strings.TrimSpace(secret)
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "…" + s[len(s)-4:]
}
