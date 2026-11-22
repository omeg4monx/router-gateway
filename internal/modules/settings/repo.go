package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

var (
	// ErrEmptyUpdates indicates a PATCH request supplied no setting updates.
	ErrEmptyUpdates = errors.New("settings: updates must contain at least one key")
	// ErrRevisionNotFound indicates the requested settings_history entry does not exist.
	ErrRevisionNotFound = errors.New("settings: revision not found")
)

type storedSetting struct {
	Key       string
	ValueJSON string
	Revision  int64
	UpdatedAt string
	UpdatedBy string
}

type validatedUpdate struct {
	Def           settingDefinition
	Value         any
	CanonicalJSON string
}

type applyResult struct {
	Settings []SettingItem
	Event    UpdatedEvent
}

type rollbackResult struct {
	Item  SettingItem
	Event UpdatedEvent
}

type repository struct {
	db *db.DB
}

func newRepository(database *db.DB) *repository {
	return &repository{db: database}
}

func (r *repository) List(ctx context.Context) ([]SettingItem, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("settings: database is not initialized")
	}
	return listSettingsWithQueryer(ctx, r.db.Reader())
}

type rowQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func listSettingsWithQueryer(ctx context.Context, q rowQueryer) ([]SettingItem, error) {
	rows, err := q.QueryContext(
		ctx,
		`SELECT key, value_json, revision, updated_at, updated_by FROM settings`,
	)
	if err != nil {
		return nil, fmt.Errorf("settings: query settings: %w", err)
	}
	defer rows.Close()

	persisted := make(map[string]storedSetting)
	for rows.Next() {
		var row storedSetting
		if err := rows.Scan(&row.Key, &row.ValueJSON, &row.Revision, &row.UpdatedAt, &row.UpdatedBy); err != nil {
			return nil, fmt.Errorf("settings: scan setting row: %w", err)
		}
		persisted[row.Key] = row
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("settings: iterate settings: %w", err)
	}

	items := make([]SettingItem, 0, len(builtinDefinitions))
	for _, def := range builtinDefinitions {
		if row, ok := persisted[def.Key]; ok {
			val, err := decodeStoredValue(def.Key, row.ValueJSON)
			if err != nil {
				return nil, err
			}
			items = append(items, def.buildItem(val, row.Revision, row.UpdatedAt, row.UpdatedBy))
			continue
		}
		items = append(items, def.buildItem(def.DefaultValue, 0, "", ""))
	}
	return items, nil
}

func (r *repository) ApplyUpdates(ctx context.Context, rawUpdates map[string]json.RawMessage, changedBy, reason string) (applyResult, error) {
	if r == nil || r.db == nil {
		return applyResult{}, errors.New("settings: database is not initialized")
	}
	if len(rawUpdates) == 0 {
		return applyResult{}, ErrEmptyUpdates
	}

	keys := make([]string, 0, len(rawUpdates))
	for k := range rawUpdates {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return applyResult{}, fmt.Errorf("settings: begin patch transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Validate every key upfront inside the transaction before writing any row.
	defs := make([]settingDefinition, 0, len(keys))
	for _, k := range keys {
		def, ok := lookupDefinition(k)
		if !ok {
			return applyResult{}, &ValidationError{
				Code:    codeInvalidSettingKey,
				Message: fmt.Sprintf("unknown setting key %q", k),
			}
		}
		defs = append(defs, def)
	}

	validated := make([]validatedUpdate, 0, len(keys))
	updatedValues := make(map[string]any, len(keys))
	for _, def := range defs {
		val, canonical, err := def.validateRaw(rawUpdates[def.Key])
		if err != nil {
			return applyResult{}, err
		}
		validated = append(validated, validatedUpdate{
			Def:           def,
			Value:         val,
			CanonicalJSON: canonical,
		})
		updatedValues[def.Key] = val
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	cleanReason := strings.TrimSpace(reason)

	for _, upd := range validated {
		var (
			currentValJSON sql.NullString
			nextRevision   int64
			existingVal    string
			existingRev    int64
		)
		err := tx.QueryRowContext(
			ctx,
			`SELECT value_json, revision FROM settings WHERE key = ?`,
			upd.Def.Key,
		).Scan(&existingVal, &existingRev)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			currentValJSON = sql.NullString{Valid: false}
			nextRevision = 1
		case err != nil:
			return applyResult{}, fmt.Errorf("settings: read existing %q: %w", upd.Def.Key, err)
		default:
			currentValJSON = sql.NullString{String: existingVal, Valid: true}
			nextRevision = existingRev + 1
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO settings (key, value_json, revision, updated_at, updated_by)
			 VALUES (?, ?, ?, ?, ?)
			 ON CONFLICT(key) DO UPDATE SET
			     value_json = excluded.value_json,
			     revision = excluded.revision,
			     updated_at = excluded.updated_at,
			     updated_by = excluded.updated_by`,
			upd.Def.Key,
			upd.CanonicalJSON,
			nextRevision,
			now,
			changedBy,
		); err != nil {
			return applyResult{}, fmt.Errorf("settings: upsert %q: %w", upd.Def.Key, err)
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO settings_history (setting_key, old_value_json, new_value_json, revision, changed_by, reason, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			upd.Def.Key,
			currentValJSON,
			upd.CanonicalJSON,
			nextRevision,
			changedBy,
			cleanReason,
			now,
		); err != nil {
			return applyResult{}, fmt.Errorf("settings: insert history for %q: %w", upd.Def.Key, err)
		}
	}

	snapshot, err := listSettingsWithQueryer(ctx, tx)
	if err != nil {
		return applyResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return applyResult{}, fmt.Errorf("settings: commit patch transaction: %w", err)
	}

	return applyResult{
		Settings: snapshot,
		Event: UpdatedEvent{
			Keys:      keys,
			Values:    updatedValues,
			ChangedBy: changedBy,
			Reason:    cleanReason,
			UpdatedAt: now,
		},
	}, nil
}

func (r *repository) ListHistory(ctx context.Context, key string, limit int) ([]HistoryItem, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("settings: database is not initialized")
	}

	var (
		rows *sql.Rows
		err  error
	)
	if key != "" {
		rows, err = r.db.Reader().QueryContext(
			ctx,
			`SELECT id, setting_key, old_value_json, new_value_json, revision, changed_by, reason, created_at
			 FROM settings_history
			 WHERE setting_key = ?
			 ORDER BY created_at DESC, id DESC
			 LIMIT ?`,
			key,
			limit,
		)
	} else {
		rows, err = r.db.Reader().QueryContext(
			ctx,
			`SELECT id, setting_key, old_value_json, new_value_json, revision, changed_by, reason, created_at
			 FROM settings_history
			 ORDER BY created_at DESC, id DESC
			 LIMIT ?`,
			limit,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("settings: query history: %w", err)
	}
	defer rows.Close()

	items := make([]HistoryItem, 0)
	for rows.Next() {
		var (
			item    HistoryItem
			oldJSON sql.NullString
			newJSON string
		)
		if err := rows.Scan(
			&item.ID,
			&item.SettingKey,
			&oldJSON,
			&newJSON,
			&item.Revision,
			&item.ChangedBy,
			&item.Reason,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("settings: scan history row: %w", err)
		}

		if oldJSON.Valid {
			val, err := decodeStoredValue(item.SettingKey, oldJSON.String)
			if err != nil {
				return nil, err
			}
			item.OldValue = val
		}
		newVal, err := decodeStoredValue(item.SettingKey, newJSON)
		if err != nil {
			return nil, err
		}
		item.NewValue = newVal

		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("settings: iterate history: %w", err)
	}

	return items, nil
}

func (r *repository) Rollback(ctx context.Context, historyID int64, changedBy, reason string) (rollbackResult, error) {
	if r == nil || r.db == nil {
		return rollbackResult{}, errors.New("settings: database is not initialized")
	}

	tx, err := r.db.Writer().BeginTx(ctx, nil)
	if err != nil {
		return rollbackResult{}, fmt.Errorf("settings: begin rollback transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var (
		settingKey   string
		oldValueJSON sql.NullString
	)
	err = tx.QueryRowContext(
		ctx,
		`SELECT setting_key, old_value_json FROM settings_history WHERE id = ?`,
		historyID,
	).Scan(&settingKey, &oldValueJSON)
	if errors.Is(err, sql.ErrNoRows) {
		return rollbackResult{}, ErrRevisionNotFound
	}
	if err != nil {
		return rollbackResult{}, fmt.Errorf("settings: query history #%d: %w", historyID, err)
	}

	def, ok := lookupDefinition(settingKey)
	if !ok {
		return rollbackResult{}, &ValidationError{
			Code:    codeInvalidSettingKey,
			Message: fmt.Sprintf("unknown setting key %q in history record", settingKey),
		}
	}

	restoreRaw := def.defaultJSON()
	if oldValueJSON.Valid {
		restoreRaw = oldValueJSON.String
	}

	restoredVal, restoredCanonical, err := def.validateRaw(json.RawMessage(restoreRaw))
	if err != nil {
		return rollbackResult{}, err
	}

	var (
		currentValJSON sql.NullString
		nextRevision   int64
		existingVal    string
		existingRev    int64
	)
	err = tx.QueryRowContext(
		ctx,
		`SELECT value_json, revision FROM settings WHERE key = ?`,
		def.Key,
	).Scan(&existingVal, &existingRev)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		currentValJSON = sql.NullString{Valid: false}
		nextRevision = 1
	case err != nil:
		return rollbackResult{}, fmt.Errorf("settings: read current %q for rollback: %w", def.Key, err)
	default:
		currentValJSON = sql.NullString{String: existingVal, Valid: true}
		nextRevision = existingRev + 1
	}

	now := time.Now().UTC().Format(rfc3339NanoFixed)
	effectiveReason := strings.TrimSpace(reason)
	if effectiveReason == "" {
		effectiveReason = fmt.Sprintf("rollback history #%d", historyID)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO settings (key, value_json, revision, updated_at, updated_by)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET
		     value_json = excluded.value_json,
		     revision = excluded.revision,
		     updated_at = excluded.updated_at,
		     updated_by = excluded.updated_by`,
		def.Key,
		restoredCanonical,
		nextRevision,
		now,
		changedBy,
	); err != nil {
		return rollbackResult{}, fmt.Errorf("settings: upsert rollback for %q: %w", def.Key, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO settings_history (setting_key, old_value_json, new_value_json, revision, changed_by, reason, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		def.Key,
		currentValJSON,
		restoredCanonical,
		nextRevision,
		changedBy,
		effectiveReason,
		now,
	); err != nil {
		return rollbackResult{}, fmt.Errorf("settings: insert rollback history for %q: %w", def.Key, err)
	}

	if err := tx.Commit(); err != nil {
		return rollbackResult{}, fmt.Errorf("settings: commit rollback transaction: %w", err)
	}

	item := def.buildItem(restoredVal, nextRevision, now, changedBy)
	return rollbackResult{
		Item: item,
		Event: UpdatedEvent{
			Keys:      []string{def.Key},
			Values:    map[string]any{def.Key: restoredVal},
			ChangedBy: changedBy,
			Reason:    effectiveReason,
			Rollback:  true,
			UpdatedAt: now,
		},
	}, nil
}
