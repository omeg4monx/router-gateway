CREATE TABLE IF NOT EXISTS proxy_pools (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    scheme TEXT NOT NULL,
    url_encrypted TEXT NOT NULL,
    display_url TEXT NOT NULL,
    has_credentials INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'unknown',
    last_latency_ms INTEGER NOT NULL DEFAULT 0,
    last_error TEXT NOT NULL DEFAULT '',
    last_tested_at TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS proxy_bindings (
    provider_id TEXT PRIMARY KEY,
    pool_id TEXT NOT NULL REFERENCES proxy_pools(id) ON DELETE CASCADE,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_proxy_bindings_pool_id ON proxy_bindings(pool_id);
