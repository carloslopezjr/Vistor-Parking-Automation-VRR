-- 003_create_reminder_tokens.sql

CREATE TABLE IF NOT EXISTS reminder_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    token_hash BLOB NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at DATETIME,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_tokens_profile FOREIGN KEY (profile_id)
        REFERENCES profiles(id) ON DELETE CASCADE,
    CONSTRAINT uq_token_hash UNIQUE (token_hash)
);

CREATE INDEX IF NOT EXISTS idx_tokens_profile ON reminder_tokens (profile_id);
CREATE INDEX IF NOT EXISTS idx_tokens_expires ON reminder_tokens (expires_at);
