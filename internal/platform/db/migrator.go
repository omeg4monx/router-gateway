package db

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type migrationFile struct {
	version int
	name    string
	path    string
}

// MigrateModule discovers versioned SQL migration files (NNNN_*.sql) in fsys,
// sorts them by version ascending, and applies each unapplied migration for
// module inside an atomic writer transaction. If fsys is nil, MigrateModule is
// a no-op.
func (d *DB) MigrateModule(ctx context.Context, module string, fsys fs.FS) error {
	if d == nil || d.writer == nil {
		return errors.New("database is not initialized")
	}

	mod := strings.TrimSpace(module)
	if mod == "" {
		return errors.New("module name must not be empty")
	}

	if fsys == nil {
		return nil
	}

	migrations, err := collectMigrations(fsys)
	if err != nil {
		return fmt.Errorf("collect migrations for module %q: %w", mod, err)
	}
	if len(migrations) == 0 {
		return nil
	}

	applied, err := d.appliedVersionSet(ctx, mod)
	if err != nil {
		return fmt.Errorf("load applied versions for module %q: %w", mod, err)
	}

	for _, m := range migrations {
		if applied[m.version] {
			continue
		}

		sqlBytes, err := fs.ReadFile(fsys, m.path)
		if err != nil {
			return fmt.Errorf("read migration %q for module %q: %w", m.path, mod, err)
		}

		tx, err := d.writer.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration tx for module %q version %d (%s): %w", mod, m.version, m.path, err)
		}

		if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("execute migration for module %q version %d (%s): %w", mod, m.version, m.path, err)
		}

		appliedAt := time.Now().UTC().Format(time.RFC3339Nano)
		if _, err := tx.ExecContext(
			ctx,
			"INSERT INTO schema_migrations (module, version, applied_at) VALUES (?, ?, ?)",
			mod,
			m.version,
			appliedAt,
		); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration for module %q version %d (%s): %w", mod, m.version, m.path, err)
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("commit migration for module %q version %d (%s): %w", mod, m.version, m.path, err)
		}

		applied[m.version] = true
	}

	return nil
}

// AppliedVersions returns the sorted list of migration versions recorded for
// module in schema_migrations.
func (d *DB) AppliedVersions(ctx context.Context, module string) ([]int, error) {
	if d == nil || d.reader == nil {
		return nil, errors.New("database is not initialized")
	}

	mod := strings.TrimSpace(module)
	if mod == "" {
		return nil, errors.New("module name must not be empty")
	}

	rows, err := d.reader.QueryContext(
		ctx,
		"SELECT version FROM schema_migrations WHERE module = ? ORDER BY version ASC",
		mod,
	)
	if err != nil {
		return nil, fmt.Errorf("query applied versions for module %q: %w", mod, err)
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied version for module %q: %w", mod, err)
		}
		versions = append(versions, version)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied versions for module %q: %w", mod, err)
	}

	return versions, nil
}

func (d *DB) appliedVersionSet(ctx context.Context, module string) (map[int]bool, error) {
	rows, err := d.writer.QueryContext(
		ctx,
		"SELECT version FROM schema_migrations WHERE module = ?",
		module,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applied, nil
}

func collectMigrations(fsys fs.FS) ([]migrationFile, error) {
	seen := make(map[int]string)
	var migrations []migrationFile

	if err := collectMigrationsFromDir(fsys, ".", seen, &migrations); err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func collectMigrationsFromDir(fsys fs.FS, dir string, seen map[int]string, out *[]migrationFile) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("read migration directory %q: %w", dir, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		entryPath := name
		if dir != "." {
			entryPath = path.Join(dir, name)
		}

		if entry.IsDir() {
			if err := collectMigrationsFromDir(fsys, entryPath, seen, out); err != nil {
				return err
			}
			continue
		}

		if !strings.HasSuffix(strings.ToLower(name), ".sql") {
			continue
		}

		version, err := parseMigrationVersion(name)
		if err != nil {
			return fmt.Errorf("invalid migration file %q: %w", entryPath, err)
		}

		if prev, exists := seen[version]; exists {
			return fmt.Errorf("duplicate migration version %d in %q and %q", version, prev, entryPath)
		}
		seen[version] = entryPath

		*out = append(*out, migrationFile{
			version: version,
			name:    name,
			path:    entryPath,
		})
	}

	return nil
}

func parseMigrationVersion(filename string) (int, error) {
	stem := filename[:len(filename)-len(".sql")]
	sep := strings.IndexByte(stem, '_')
	if sep <= 0 || sep == len(stem)-1 {
		return 0, fmt.Errorf("filename %q must match NNNN_description.sql", filename)
	}

	prefix := stem[:sep]
	for i := 0; i < len(prefix); i++ {
		if prefix[i] < '0' || prefix[i] > '9' {
			return 0, fmt.Errorf("version prefix %q in %q must be digits", prefix, filename)
		}
	}

	version, err := strconv.Atoi(prefix)
	if err != nil || version <= 0 {
		return 0, fmt.Errorf("version %q in %q must be a positive integer", prefix, filename)
	}

	return version, nil
}
