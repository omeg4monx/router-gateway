package db_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/omeg4monx/router-gateway/internal/platform/db"
)

func TestOpenAppliesWALPragmasAndSplitPools(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "nested", "data.sqlite")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) failed: %v", dbPath, err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	info, err := os.Stat(filepath.Dir(dbPath))
	if err != nil {
		t.Fatalf("expected parent dir to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", filepath.Dir(dbPath))
	}

	if err := database.Ping(ctx); err != nil {
		t.Fatalf("Ping() failed: %v", err)
	}

	if got := database.Writer().Stats().MaxOpenConnections; got != 1 {
		t.Fatalf("Writer().Stats().MaxOpenConnections = %d, want 1", got)
	}
	if got := database.Reader().Stats().MaxOpenConnections; got != 8 {
		t.Fatalf("Reader().Stats().MaxOpenConnections = %d, want 8", got)
	}

	var journalMode string
	if err := database.Writer().QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("query writer PRAGMA journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Fatalf("writer PRAGMA journal_mode = %q, want %q", journalMode, "wal")
	}

	var readerJournalMode string
	if err := database.Reader().QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&readerJournalMode); err != nil {
		t.Fatalf("query reader PRAGMA journal_mode: %v", err)
	}
	if readerJournalMode != "wal" {
		t.Fatalf("reader PRAGMA journal_mode = %q, want %q", readerJournalMode, "wal")
	}

	checkIntPragma := func(poolName string, query string, want int) {
		t.Helper()
		var got int
		pool := database.Writer()
		if poolName == "reader" {
			pool = database.Reader()
		}
		if err := pool.QueryRowContext(ctx, query).Scan(&got); err != nil {
			t.Fatalf("query %s %s: %v", poolName, query, err)
		}
		if got != want {
			t.Fatalf("%s %s = %d, want %d", poolName, query, got, want)
		}
	}

	checkIntPragma("writer", "PRAGMA foreign_keys", 1)
	checkIntPragma("writer", "PRAGMA busy_timeout", 10000)
	checkIntPragma("writer", "PRAGMA synchronous", 1) // NORMAL = 1
	checkIntPragma("writer", "PRAGMA temp_store", 2)  // MEMORY = 2
	checkIntPragma("writer", "PRAGMA cache_size", -2000)
	checkIntPragma("writer", "PRAGMA mmap_size", 0)

	checkIntPragma("reader", "PRAGMA foreign_keys", 1)
	checkIntPragma("reader", "PRAGMA busy_timeout", 10000)
	checkIntPragma("reader", "PRAGMA query_only", 1)

	if _, err := database.Reader().ExecContext(ctx, "CREATE TABLE should_fail (id INTEGER PRIMARY KEY)"); err == nil {
		t.Fatalf("expected Reader() write statement to fail under read-only query path, got nil")
	}
}

func TestMigrateModuleIdempotentAndAtomicRollback(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "migrate.sqlite")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	// Nil filesystem is a no-op.
	if err := database.MigrateModule(ctx, "alpha", nil); err != nil {
		t.Fatalf("MigrateModule(nil) failed: %v", err)
	}

	// Empty module names must be rejected.
	if err := database.MigrateModule(ctx, "   ", fstest.MapFS{}); err == nil {
		t.Fatalf("expected error for blank module name in MigrateModule")
	}
	if _, err := database.AppliedVersions(ctx, ""); err == nil {
		t.Fatalf("expected error for blank module name in AppliedVersions")
	}

	// Initial two migrations applied in ascending version order.
	initialFS := fstest.MapFS{
		"0002_add_col.sql": &fstest.MapFile{
			Data: []byte(`
				ALTER TABLE items ADD COLUMN label TEXT NOT NULL DEFAULT 'default';
				INSERT INTO items (name, label) VALUES ('seed', 'initial');
			`),
		},
		"0001_init.sql": &fstest.MapFile{
			Data: []byte(`
				CREATE TABLE items (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT NOT NULL
				);
			`),
		},
	}

	if err := database.MigrateModule(ctx, "alpha", initialFS); err != nil {
		t.Fatalf("first MigrateModule() failed: %v", err)
	}

	versions, err := database.AppliedVersions(ctx, "alpha")
	if err != nil {
		t.Fatalf("AppliedVersions() failed: %v", err)
	}
	if !reflect.DeepEqual(versions, []int{1, 2}) {
		t.Fatalf("AppliedVersions(alpha) = %v, want [1 2]", versions)
	}

	// Re-running the exact same migrations must be a no-op (seed row count remains 1).
	if err := database.MigrateModule(ctx, "alpha", initialFS); err != nil {
		t.Fatalf("idempotent re-run MigrateModule() failed: %v", err)
	}

	var count int
	if err := database.Reader().QueryRowContext(ctx, "SELECT COUNT(*) FROM items").Scan(&count); err != nil {
		t.Fatalf("count items: %v", err)
	}
	if count != 1 {
		t.Fatalf("items count after idempotent re-run = %d, want 1", count)
	}

	// Failing 0003 migration with valid DDL followed by syntax error must roll back atomically.
	brokenFS := fstest.MapFS{
		"0001_init.sql":    initialFS["0001_init.sql"],
		"0002_add_col.sql": initialFS["0002_add_col.sql"],
		"0003_bad.sql": &fstest.MapFile{
			Data: []byte(`
				CREATE TABLE partial_should_rollback (id INTEGER PRIMARY KEY);
				INSERT INTO partial_should_rollback (id) VALUES (1);
				THIS IS NOT VALID SQL SYNTAX;
			`),
		},
	}

	if err := database.MigrateModule(ctx, "alpha", brokenFS); err == nil {
		t.Fatalf("expected error from broken 0003_bad.sql migration, got nil")
	}

	versions, err = database.AppliedVersions(ctx, "alpha")
	if err != nil {
		t.Fatalf("AppliedVersions() after rollback failed: %v", err)
	}
	if !reflect.DeepEqual(versions, []int{1, 2}) {
		t.Fatalf("AppliedVersions(alpha) after rollback = %v, want [1 2]", versions)
	}

	var partialTableCount int
	if err := database.Reader().QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'partial_should_rollback'",
	).Scan(&partialTableCount); err != nil {
		t.Fatalf("query sqlite_master: %v", err)
	}
	if partialTableCount != 0 {
		t.Fatalf("partial DDL table exists after failed migration rollback")
	}

	// Nested *.sql files and independent module tracking.
	nestedFS := fstest.MapFS{
		"migrations/0001_nested.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE beta_records (id INTEGER PRIMARY KEY);`),
		},
	}
	if err := database.MigrateModule(ctx, "beta", nestedFS); err != nil {
		t.Fatalf("nested MigrateModule(beta) failed: %v", err)
	}
	betaVersions, err := database.AppliedVersions(ctx, "beta")
	if err != nil {
		t.Fatalf("AppliedVersions(beta) failed: %v", err)
	}
	if !reflect.DeepEqual(betaVersions, []int{1}) {
		t.Fatalf("AppliedVersions(beta) = %v, want [1]", betaVersions)
	}

	// Duplicate migration versions must be rejected.
	dupFS := fstest.MapFS{
		"0001_first.sql":  &fstest.MapFile{Data: []byte(`SELECT 1;`)},
		"1_duplicate.sql": &fstest.MapFile{Data: []byte(`SELECT 1;`)},
	}
	if err := database.MigrateModule(ctx, "gamma", dupFS); err == nil {
		t.Fatalf("expected error for duplicate migration versions, got nil")
	}

	// Non-positive version must be rejected.
	zeroFS := fstest.MapFS{
		"0000_invalid.sql": &fstest.MapFile{Data: []byte(`SELECT 1;`)},
	}
	if err := database.MigrateModule(ctx, "delta", zeroFS); err == nil {
		t.Fatalf("expected error for non-positive migration version 0, got nil")
	}
}

func TestConcurrentReadersAndWritersUnderWAL(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "concurrent.sqlite")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	setupFS := fstest.MapFS{
		"0001_events.sql": &fstest.MapFile{
			Data: []byte(`
				CREATE TABLE events (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					worker INTEGER NOT NULL,
					seq INTEGER NOT NULL,
					payload TEXT NOT NULL
				);
			`),
		},
	}
	if err := database.MigrateModule(ctx, "events", setupFS); err != nil {
		t.Fatalf("MigrateModule() failed: %v", err)
	}

	const (
		numReaders = 8
		numWriters = 4
		opsPerTask = 50
	)

	var (
		wg   sync.WaitGroup
		errs = make(chan error, (numReaders+numWriters)*opsPerTask)
	)

	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < opsPerTask; i++ {
				_, err := database.Writer().ExecContext(
					ctx,
					"INSERT INTO events (worker, seq, payload) VALUES (?, ?, ?)",
					workerID,
					i,
					fmt.Sprintf("w%d-op%d", workerID, i),
				)
				if err != nil {
					errs <- fmt.Errorf("writer %d op %d: %w", workerID, i, err)
					return
				}
			}
		}(w)
	}

	for r := 0; r < numReaders; r++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for i := 0; i < opsPerTask; i++ {
				var count int
				if err := database.Reader().QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&count); err != nil {
					errs <- fmt.Errorf("reader %d op %d: %w", readerID, i, err)
					return
				}
			}
		}(r)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent WAL operation failed: %v", err)
	}

	var total int
	if err := database.Reader().QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&total); err != nil {
		t.Fatalf("final count failed: %v", err)
	}
	if want := numWriters * opsPerTask; total != want {
		t.Fatalf("final row count = %d, want %d", total, want)
	}
}
