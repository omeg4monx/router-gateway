CREATE TABLE IF NOT EXISTS teams (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    budget_usd REAL NOT NULL DEFAULT 0,
    spent_usd REAL NOT NULL DEFAULT 0,
    rate_limit_rpm INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS team_keys (
    team_id TEXT NOT NULL,
    key_id TEXT NOT NULL PRIMARY KEY,
    added_at TEXT NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_team_keys_team_id ON team_keys(team_id);
