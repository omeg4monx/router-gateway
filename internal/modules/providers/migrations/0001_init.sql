CREATE TABLE IF NOT EXISTS provider_nodes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    protocol TEXT NOT NULL,
    base_url TEXT NOT NULL,
    models_json TEXT NOT NULL DEFAULT '[]',
    headers_json TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS provider_connections (
    id TEXT PRIMARY KEY,
    provider_id TEXT NOT NULL,
    name TEXT NOT NULL,
    protocol TEXT NOT NULL,
    base_url TEXT NOT NULL,
    auth_type TEXT NOT NULL,
    credentials_enc TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 100,
    enabled INTEGER NOT NULL DEFAULT 1,
    models_json TEXT NOT NULL DEFAULT '[]',
    headers_json TEXT NOT NULL DEFAULT '{}',
    token_expires_at TEXT,
    last_tested_at TEXT,
    last_test_status TEXT NOT NULL DEFAULT '',
    last_test_error TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_provider_connections_provider_id
    ON provider_connections(provider_id);

CREATE INDEX IF NOT EXISTS idx_provider_connections_enabled_priority
    ON provider_connections(enabled, priority, id);

CREATE TABLE IF NOT EXISTS oauth_sessions (
    state TEXT PRIMARY KEY,
    provider_id TEXT NOT NULL,
    connection_name TEXT NOT NULL,
    code_verifier_enc TEXT NOT NULL,
    redirect_uri TEXT NOT NULL,
    client_id TEXT NOT NULL,
    client_secret_enc TEXT NOT NULL DEFAULT '',
    base_url TEXT NOT NULL DEFAULT '',
    models_json TEXT NOT NULL DEFAULT '[]',
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_oauth_sessions_expires_at
    ON oauth_sessions(expires_at);
