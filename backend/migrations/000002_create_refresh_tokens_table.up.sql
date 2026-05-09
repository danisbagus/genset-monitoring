-- ══════════════════════════════════════════════════════════════════
-- Migration: 000002_create_refresh_tokens_table (UP)
-- Description: Creates refresh tokens table for secure token rotation
-- ══════════════════════════════════════════════════════════════════

-- ── Refresh tokens table ───────────────────────────────────────────
-- Security design:
--   • token_hash stores a SHA-256 / bcrypt hash of the opaque token,
--     so even if the DB is compromised raw tokens are not leaked.
--   • Each refresh generates a NEW token (rotation) and revokes the old one.
--   • Revoked tokens are kept for audit; a background job can prune expired rows.
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL,
    token_hash  TEXT        NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN     NOT NULL DEFAULT FALSE,
    revoked_at  TIMESTAMPTZ NULL,
    user_agent  TEXT        NULL,
    ip_address  VARCHAR(45) NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);

-- ── Indexes ────────────────────────────────────────────────────────
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_hash
    ON refresh_tokens (token_hash);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id
    ON refresh_tokens (user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at
    ON refresh_tokens (expires_at);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked
    ON refresh_tokens (revoked)
    WHERE revoked = FALSE;
