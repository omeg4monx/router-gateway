CREATE TABLE IF NOT EXISTS cloudsync_config (
    id               INTEGER PRIMARY KEY CHECK (id = 1),
    enabled          INTEGER NOT NULL DEFAULT 0,
    endpoint_url     TEXT NOT NULL DEFAULT '',
    token_enc        TEXT NOT NULL DEFAULT '',
    token_mask       TEXT NOT NULL DEFAULT '',
    sync_interval_min INTEGER NOT NULL DEFAULT 30,
    conflict_policy  TEXT NOT NULL DEFAULT 'merge',
    include_settings INTEGER NOT NULL DEFAULT 1,
    include_combos   INTEGER NOT NULL DEFAULT 1,
    last_sync_at     TEXT NOT NULL DEFAULT '',
    last_status      TEXT NOT NULL DEFAULT '',
    last_error       TEXT NOT NULL DEFAULT '',
    updated_at       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cloudsync_history (
    id         TEXT PRIMARY KEY,
    direction  TEXT NOT NULL,
    status     TEXT NOT NULL,
    summary    TEXT NOT NULL DEFAULT '',
    checksum   TEXT NOT NULL DEFAULT '',
    item_count INTEGER NOT NULL DEFAULT 0,
    error      TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_cloudsync_history_created ON cloudsync_history(created_at DESC);
