CREATE TABLE IF NOT EXISTS combos (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    targets_json TEXT NOT NULL,
    strategy TEXT NOT NULL DEFAULT 'ordered',
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_combos_enabled_name
    ON combos(enabled, name);

CREATE TABLE IF NOT EXISTS model_aliases (
    alias TEXT PRIMARY KEY,
    target_model TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS disabled_models (
    model TEXT PRIMARY KEY,
    reason TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);
