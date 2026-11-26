CREATE TABLE IF NOT EXISTS proxy_bindings_v2 (
    provider_id TEXT PRIMARY KEY,
    pool_id TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

INSERT OR IGNORE INTO proxy_bindings_v2 (provider_id, pool_id, updated_at)
SELECT provider_id, pool_id, updated_at FROM proxy_bindings;

DROP TABLE proxy_bindings;

ALTER TABLE proxy_bindings_v2 RENAME TO proxy_bindings;

CREATE INDEX IF NOT EXISTS idx_proxy_bindings_pool_id ON proxy_bindings(pool_id);
