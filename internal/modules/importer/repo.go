package importer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

const rfc3339NanoFixed = "2006-01-02T15:04:05.000000000Z"

// ImportRun records a dry-run or live legacy data.sqlite migration run.
type ImportRun struct {
	ID               string `json:"id"`
	SourcePath       string `json:"source_path"`
	DryRun           bool   `json:"dry_run"`
	Status           string `json:"status"`
	SettingsImported int    `json:"settings_imported"`
	CombosImported   int    `json:"combos_imported"`
	KeysImported     int    `json:"keys_imported"`
	ProvidersFound   int    `json:"providers_found"`
	Summary          string `json:"summary"`
	Error            string `json:"error"`
	CreatedAt        string `json:"created_at"`
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

func (r *repository) InsertRun(ctx context.Context, run ImportRun) (ImportRun, error) {
	raw := make([]byte, 12)
	if _, err := rand.Read(raw); err != nil {
		return ImportRun{}, fmt.Errorf("importer: generate run id: %w", err)
	}
	run.ID = "imp_" + hex.EncodeToString(raw)
	run.CreatedAt = r.now().UTC().Format(rfc3339NanoFixed)
	_, err := r.db.Writer().ExecContext(
		ctx,
		`INSERT INTO importer_runs
		 (id, source_path, dry_run, status, settings_imported, combos_imported, keys_imported, providers_found, summary, error, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		run.ID, run.SourcePath, boolToInt(run.DryRun), run.Status,
		run.SettingsImported, run.CombosImported, run.KeysImported, run.ProvidersFound,
		run.Summary, run.Error, run.CreatedAt,
	)
	if err != nil {
		return ImportRun{}, fmt.Errorf("importer: insert run: %w", err)
	}
	return run, nil
}

func (r *repository) ListRuns(ctx context.Context, limit int) ([]ImportRun, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Reader().QueryContext(
		ctx,
		`SELECT id, source_path, dry_run, status, settings_imported, combos_imported, keys_imported, providers_found, summary, error, created_at
		 FROM importer_runs
		 ORDER BY created_at DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("importer: query runs: %w", err)
	}
	defer rows.Close()

	runs := make([]ImportRun, 0)
	for rows.Next() {
		var (
			run    ImportRun
			dryInt int
		)
		if err := rows.Scan(
			&run.ID, &run.SourcePath, &dryInt, &run.Status,
			&run.SettingsImported, &run.CombosImported, &run.KeysImported, &run.ProvidersFound,
			&run.Summary, &run.Error, &run.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("importer: scan run: %w", err)
		}
		run.DryRun = dryInt != 0
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("importer: iterate runs: %w", err)
	}
	return runs, nil
}

func (r *repository) ClearRuns(ctx context.Context) (int64, error) {
	res, err := r.db.Writer().ExecContext(ctx, `DELETE FROM importer_runs`)
	if err != nil {
		return 0, fmt.Errorf("importer: clear runs: %w", err)
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
