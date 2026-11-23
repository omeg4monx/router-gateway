CREATE TABLE IF NOT EXISTS pxpipe_config (
    id                    INTEGER PRIMARY KEY CHECK (id = 1),
    enabled               INTEGER NOT NULL DEFAULT 1,
    auto_start            INTEGER NOT NULL DEFAULT 1,
    strip_null_fields     INTEGER NOT NULL DEFAULT 1,
    collapse_whitespace   INTEGER NOT NULL DEFAULT 1,
    dedupe_system_prompts INTEGER NOT NULL DEFAULT 1,
    max_history_messages  INTEGER NOT NULL DEFAULT 0,
    updated_at            TEXT    NOT NULL
);
