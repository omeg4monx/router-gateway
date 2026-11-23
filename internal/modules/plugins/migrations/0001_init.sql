CREATE TABLE IF NOT EXISTS plugin_configs (
    key         TEXT PRIMARY KEY,
    enabled     INTEGER NOT NULL DEFAULT 0,
    models_json TEXT NOT NULL DEFAULT '[]',
    config_json TEXT NOT NULL DEFAULT '{}',
    updated_at  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS plugin_skills (
    id          TEXT PRIMARY KEY,
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category    TEXT NOT NULL DEFAULT 'general',
    content     TEXT NOT NULL,
    enabled     INTEGER NOT NULL DEFAULT 1,
    builtin     INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_plugin_skills_slug ON plugin_skills(slug);
CREATE INDEX IF NOT EXISTS idx_plugin_skills_category ON plugin_skills(category);
