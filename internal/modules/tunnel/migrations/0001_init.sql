CREATE TABLE IF NOT EXISTS tunnel_config (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    provider TEXT NOT NULL DEFAULT 'cloudflare',
    local_url TEXT NOT NULL DEFAULT '',
    tailscale_mode TEXT NOT NULL DEFAULT 'funnel',
    auto_start INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tunnel_events (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    event_type TEXT NOT NULL,
    message TEXT NOT NULL,
    public_url TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tunnel_events_created_at
    ON tunnel_events(created_at DESC);
