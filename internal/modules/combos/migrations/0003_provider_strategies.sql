CREATE TABLE IF NOT EXISTS provider_strategies (
    provider_id TEXT PRIMARY KEY,
    strategy TEXT NOT NULL DEFAULT 'fill-first',
    updated_at TEXT NOT NULL
);
