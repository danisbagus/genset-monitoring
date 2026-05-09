-- ══════════════════════════════════════════════════════════════════
-- Migration: 000001_create_users_table (UP)
-- Description: Creates the users table with security and audit fields
-- ══════════════════════════════════════════════════════════════════

-- Enable pgcrypto for gen_random_uuid() (idempotent)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── Users table ────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'viewer',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ  NULL
);

-- ── Constraints ────────────────────────────────────────────────────
ALTER TABLE users
    ADD CONSTRAINT uq_users_username UNIQUE (username),
    ADD CONSTRAINT uq_users_email    UNIQUE (email),
    ADD CONSTRAINT chk_users_role    CHECK  (role IN ('admin', 'operator', 'viewer'));

-- ── Indexes ────────────────────────────────────────────────────────
-- Partial unique indexes that exclude soft-deleted rows
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_active
    ON users (username)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
    ON users (email)
    WHERE deleted_at IS NULL;

-- General query indexes
CREATE INDEX IF NOT EXISTS idx_users_deleted_at  ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_role        ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_is_active   ON users (is_active);
