CREATE TABLE IF NOT EXISTS importer_runs (
    id                TEXT PRIMARY KEY,
    source_path       TEXT NOT NULL,
    dry_run           INTEGER NOT NULL DEFAULT 1,
    status            TEXT NOT NULL,
    settings_imported INTEGER NOT NULL DEFAULT 0,
    combos_imported   INTEGER NOT NULL DEFAULT 0,
    keys_imported     INTEGER NOT NULL DEFAULT 0,
    providers_found   INTEGER NOT NULL DEFAULT 0,
    summary           TEXT NOT NULL DEFAULT '',
    error             TEXT NOT NULL DEFAULT '',
    created_at        TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_importer_runs_created ON importer_runs(created_at DESC);
