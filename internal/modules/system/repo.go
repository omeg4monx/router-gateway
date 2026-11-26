package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/contracts"
	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const (
	keyAdminPasswordHash = "admin_password_hash"
	keyRequireLogin      = "require_login"
	defaultRequireLogin  = "true"
)

// ErrAlreadyInitialized indicates that the system admin password has already
// been configured.
var ErrAlreadyInitialized = errors.New("system: already initialized")

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) InitState(ctx context.Context) (initialized, requireLogin bool, err error) {
	if r == nil || r.db == nil {
		return false, true, errors.New("system: database is not initialized")
	}

	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT key, value FROM system_state WHERE key IN (?, ?)`,
		keyAdminPasswordHash,
		keyRequireLogin,
	)
	if err != nil {
		return false, true, fmt.Errorf("system: query init state: %w", err)
	}
	defer rows.Close()

	requireLogin = true
	for rows.Next() {
		var key, val string
		if err := rows.Scan(&key, &val); err != nil {
			return false, true, fmt.Errorf("system: scan init state: %w", err)
		}
		switch key {
		case keyAdminPasswordHash:
			if strings.TrimSpace(val) != "" {
				initialized = true
			}
		case keyRequireLogin:
			requireLogin = !strings.EqualFold(strings.TrimSpace(val), "false")
		}
	}
	if err := rows.Err(); err != nil {
		return false, true, fmt.Errorf("system: iterate init state: %w", err)
	}

	return initialized, requireLogin, nil
}

func (r *repository) IsInitialized(ctx context.Context) (bool, error) {
	if r == nil || r.db == nil {
		return false, errors.New("system: database is not initialized")
	}

	var val string
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT value FROM system_state WHERE key = ?`,
		keyAdminPasswordHash,
	).Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("system: query admin password state: %w", err)
	}
	return strings.TrimSpace(val) != "", nil
}

func (r *repository) EnsureDefaultRequireLogin(ctx context.Context) error {
	if r == nil || r.db == nil {
		return errors.New("system: database is not initialized")
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT OR IGNORE INTO system_state (key, value, updated_at) VALUES (?, ?, ?)`,
		keyRequireLogin,
		defaultRequireLogin,
		now,
	); err != nil {
		return fmt.Errorf("system: ensure default require_login: %w", err)
	}
	return nil
}

// InitAdminPassword atomically inserts the initial admin password hash. If the
// key already exists, it rolls back and returns ErrAlreadyInitialized.
func (r *repository) InitAdminPassword(ctx context.Context, passwordHash string) error {
	return r.insertAdminPassword(ctx, passwordHash, false)
}

// BootstrapAdminPassword atomically inserts the initial admin password hash if
// absent and succeeds idempotently when the system is already initialized.
func (r *repository) BootstrapAdminPassword(ctx context.Context, passwordHash string) error {
	return r.insertAdminPassword(ctx, passwordHash, true)
}

func (r *repository) insertAdminPassword(ctx context.Context, passwordHash string, idempotent bool) error {
	if r == nil || r.db == nil {
		return errors.New("system: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("system: begin init transaction: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := tx.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO system_state (key, value, updated_at) VALUES (?, ?, ?)`,
		keyAdminPasswordHash,
		passwordHash,
		now,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("system: insert admin password hash: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("system: check rows affected: %w", err)
	}
	if affected == 0 && !idempotent {
		_ = tx.Rollback()
		return ErrAlreadyInitialized
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT OR IGNORE INTO system_state (key, value, updated_at) VALUES (?, ?, ?)`,
		keyRequireLogin,
		defaultRequireLogin,
		now,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("system: insert default require_login: %w", err)
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("system: commit init transaction: %w", err)
	}
	return nil
}

func (r *repository) AdminPasswordHash(ctx context.Context) (string, error) {
	if r == nil || r.db == nil {
		return "", errors.New("system: database is not initialized")
	}

	var val string
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT value FROM system_state WHERE key = ?`,
		keyAdminPasswordHash,
	).Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return "", contracts.ErrNotInitialized
	}
	if err != nil {
		return "", fmt.Errorf("system: query admin password hash: %w", err)
	}
	if strings.TrimSpace(val) == "" {
		return "", contracts.ErrNotInitialized
	}
	return val, nil
}

func (r *repository) UpdateAdminPasswordHash(ctx context.Context, hash string) error {
	if r == nil || r.db == nil {
		return errors.New("system: database is not initialized")
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := r.db.Writer().ExecContext(
		ctx,
		`UPDATE system_state SET value = ?, updated_at = ? WHERE key = ?`,
		hash,
		now,
		keyAdminPasswordHash,
	)
	if err != nil {
		return fmt.Errorf("system: update admin password hash: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("system: check updated admin password rows: %w", err)
	}
	if affected == 0 {
		return contracts.ErrNotInitialized
	}
	return nil
}

func (r *repository) RequireLogin(ctx context.Context) (bool, error) {
	if r == nil || r.db == nil {
		return true, errors.New("system: database is not initialized")
	}

	var val string
	err := r.db.Reader().QueryRowContext(
		ctx,
		`SELECT value FROM system_state WHERE key = ?`,
		keyRequireLogin,
	).Scan(&val)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return true, fmt.Errorf("system: query require_login: %w", err)
	}
	return !strings.EqualFold(strings.TrimSpace(val), "false"), nil
}

func (r *repository) SetRequireLogin(ctx context.Context, required bool) error {
	if r == nil || r.db == nil {
		return errors.New("system: database is not initialized")
	}

	val := "false"
	if required {
		val = "true"
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO system_state (key, value, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		keyRequireLogin,
		val,
		now,
	); err != nil {
		return fmt.Errorf("system: set require_login: %w", err)
	}
	return nil
}
