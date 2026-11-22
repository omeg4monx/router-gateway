// Package db provides a zero-CGO SQLite WAL engine with split writer/reader
// connection pools and a transactional per-module SQL migrator.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	driverName         = "sqlite"
	writerMaxOpenConns = 1
	readerMaxOpenConns = 8
	parentDirPerm      = 0o750
)

const schemaMigrationsDDL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	module TEXT NOT NULL,
	version INTEGER NOT NULL,
	applied_at TEXT NOT NULL,
	PRIMARY KEY (module, version)
);
`

// DB wraps split SQLite connection pools: a single-connection writer pool and a
// bounded read-only reader pool sharing the same WAL-mode database file.
type DB struct {
	writer *sql.DB
	reader *sql.DB
	path   string
}

// Open initializes a SQLite WAL database at path, creating parent directories
// (mode 0750) if needed, applying performance and durability pragmas,
// bootstrapping the schema_migrations table, and opening separate writer (max 1
// conn) and read-only reader (max 8 conns) pools.
func Open(path string) (*DB, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return nil, errors.New("database path must not be empty")
	}

	absPath, err := filepath.Abs(trimmed)
	if err != nil {
		return nil, fmt.Errorf("resolve database path %q: %w", trimmed, err)
	}

	if err := os.MkdirAll(filepath.Dir(absPath), parentDirPerm); err != nil {
		return nil, fmt.Errorf("create database directory %q: %w", filepath.Dir(absPath), err)
	}

	writerDSN := buildDSN(absPath, false)
	writer, err := sql.Open(driverName, writerDSN)
	if err != nil {
		return nil, fmt.Errorf("open sqlite writer pool: %w", err)
	}
	writer.SetMaxOpenConns(writerMaxOpenConns)
	writer.SetMaxIdleConns(writerMaxOpenConns)

	ctx := context.Background()
	if err := writer.PingContext(ctx); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("ping sqlite writer pool: %w", err)
	}

	if _, err := writer.ExecContext(ctx, schemaMigrationsDDL); err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("bootstrap schema_migrations table: %w", err)
	}

	readerDSN := buildDSN(absPath, true)
	reader, err := sql.Open(driverName, readerDSN)
	if err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("open sqlite reader pool: %w", err)
	}
	reader.SetMaxOpenConns(readerMaxOpenConns)
	reader.SetMaxIdleConns(readerMaxOpenConns)

	if err := reader.PingContext(ctx); err != nil {
		_ = reader.Close()
		_ = writer.Close()
		return nil, fmt.Errorf("ping sqlite reader pool: %w", err)
	}

	return &DB{
		writer: writer,
		reader: reader,
		path:   absPath,
	}, nil
}

// Writer returns the single-connection write pool.
func (d *DB) Writer() *sql.DB {
	return d.writer
}

// Reader returns the concurrent read-only query pool.
func (d *DB) Reader() *sql.DB {
	return d.reader
}

// Ping verifies both the writer and reader connection pools are reachable.
func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.writer == nil || d.reader == nil {
		return errors.New("database is not initialized")
	}
	if err := d.writer.PingContext(ctx); err != nil {
		return fmt.Errorf("ping writer pool: %w", err)
	}
	if err := d.reader.PingContext(ctx); err != nil {
		return fmt.Errorf("ping reader pool: %w", err)
	}
	return nil
}

// Close closes both the reader and writer pools.
func (d *DB) Close() error {
	if d == nil {
		return nil
	}
	var errs []error
	if d.reader != nil {
		if err := d.reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close reader pool: %w", err))
		}
	}
	if d.writer != nil {
		if err := d.writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close writer pool: %w", err))
		}
	}
	return errors.Join(errs...)
}

func buildDSN(absPath string, readOnly bool) string {
	q := url.Values{}
	q.Add("_pragma", "busy_timeout(10000)")
	if !readOnly {
		q.Add("_pragma", "journal_mode(WAL)")
		q.Add("_pragma", "synchronous(NORMAL)")
	}
	q.Add("_pragma", "foreign_keys(ON)")
	q.Add("_pragma", "temp_store(MEMORY)")
	q.Add("_pragma", "cache_size(-2000)")
	q.Add("_pragma", "mmap_size(0)")
	if readOnly {
		q.Add("_pragma", "query_only(ON)")
	}

	u := url.URL{
		Scheme:   "file",
		Path:     filepath.ToSlash(absPath),
		RawQuery: q.Encode(),
	}
	return u.String()
}
