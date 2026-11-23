CREATE TABLE IF NOT EXISTS usage_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT NOT NULL,
    key_id TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL,
    provider_id TEXT NOT NULL DEFAULT '',
    connection_id TEXT NOT NULL DEFAULT '',
    protocol TEXT NOT NULL DEFAULT '',
    input_tokens INTEGER NOT NULL DEFAULT 0,
    output_tokens INTEGER NOT NULL DEFAULT 0,
    cached_tokens INTEGER NOT NULL DEFAULT 0,
    cost_usd REAL NOT NULL DEFAULT 0,
    latency_ms INTEGER NOT NULL DEFAULT 0,
    ttft_ms INTEGER NOT NULL DEFAULT 0,
    status_code INTEGER NOT NULL DEFAULT 0,
    error_summary TEXT NOT NULL DEFAULT '',
    stream INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_usage_logs_created_at
    ON usage_logs(created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_usage_logs_model_created
    ON usage_logs(model, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_usage_logs_provider_created
    ON usage_logs(provider_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_usage_logs_key_created
    ON usage_logs(key_id, created_at DESC);

CREATE TABLE IF NOT EXISTS model_pricing (
    model TEXT PRIMARY KEY,
    input_per_million REAL NOT NULL DEFAULT 0,
    output_per_million REAL NOT NULL DEFAULT 0,
    cached_per_million REAL NOT NULL DEFAULT 0,
    updated_at TEXT NOT NULL
);
