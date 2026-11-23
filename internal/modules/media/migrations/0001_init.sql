CREATE TABLE IF NOT EXISTS media_providers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    capability TEXT NOT NULL,
    driver TEXT NOT NULL,
    base_url TEXT NOT NULL DEFAULT '',
    api_key_enc TEXT NOT NULL DEFAULT '',
    default_model TEXT NOT NULL DEFAULT '',
    extra_headers_json TEXT NOT NULL DEFAULT '{}',
    priority INTEGER NOT NULL DEFAULT 100,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_media_providers_cap_prio
    ON media_providers(capability, enabled, priority, created_at);

CREATE TABLE IF NOT EXISTS media_combos (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    capability TEXT NOT NULL,
    targets_json TEXT NOT NULL DEFAULT '[]',
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_media_combos_cap
    ON media_combos(capability, enabled, name);
