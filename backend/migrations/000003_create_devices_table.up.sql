-- ══════════════════════════════════════════════════════════════════
-- Migration: 000003_create_devices_table (UP)
-- Description: Creates the devices table with enhanced fields, soft
--              delete, audit trail, device token management, and
--              optimised indexes for list/search/filter operations.
-- ══════════════════════════════════════════════════════════════════

-- ── Device status enum ─────────────────────────────────────────────
-- Using an enum keeps the column self-documenting and allows the
-- query planner to use a tighter btree index scan.
CREATE TYPE device_status AS ENUM ('active', 'inactive', 'maintenance');

-- ── Devices table ──────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS devices (
    -- Primary key (UUID v4 generated server-side)
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Identity
    device_code         VARCHAR(100)    NOT NULL,
    name                VARCHAR(100)    NOT NULL,
    serial_number       VARCHAR(100)    NOT NULL,
    engine_id           VARCHAR(100)    NULL,
    gsm_number          VARCHAR(50)     NULL,

    -- Firmware
    firmware_version    VARCHAR(50)     NULL,

    -- Connectivity tracking
    is_online           BOOLEAN         NOT NULL DEFAULT FALSE,
    last_seen           TIMESTAMPTZ     NULL,

    -- Lifecycle status
    status              device_status   NOT NULL DEFAULT 'active',

    -- Device authentication token (opaque, client-facing)
    -- Only the SHA-256 hash is stored; the raw token is returned
    -- exactly once at creation time and never persisted.
    hashed_device_token TEXT            NULL,

    -- Extensible metadata (firmware config, location, tags, etc.)
    metadata            JSONB           NULL,

    -- Audit fields
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ     NULL
);

-- ── Constraints ────────────────────────────────────────────────────
-- Unique device_code among non-deleted rows (partial unique index
-- handles this more efficiently — see below; the constraint below
-- is a belt-and-suspenders guard at the DB level).
ALTER TABLE devices
    ADD CONSTRAINT uq_devices_device_code   UNIQUE (device_code),
    ADD CONSTRAINT uq_devices_serial_number UNIQUE (serial_number);

-- ── Indexes ────────────────────────────────────────────────────────

-- Partial unique index: enforce uniqueness only among soft-delete-
-- active rows so deleted records do not block reuse of a code.
DROP INDEX IF EXISTS idx_devices_device_code_active;
CREATE UNIQUE INDEX idx_devices_device_code_active
    ON devices (device_code)
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_devices_serial_number_active;
CREATE UNIQUE INDEX idx_devices_serial_number_active
    ON devices (serial_number)
    WHERE deleted_at IS NULL;

-- Soft-delete lookup (nearly every query filters deleted_at IS NULL)
CREATE INDEX IF NOT EXISTS idx_devices_deleted_at
    ON devices (deleted_at);

-- Online/offline filtering
CREATE INDEX IF NOT EXISTS idx_devices_is_online
    ON devices (is_online)
    WHERE deleted_at IS NULL;

-- Status filtering
CREATE INDEX IF NOT EXISTS idx_devices_status
    ON devices (status)
    WHERE deleted_at IS NULL;

-- Identity lookups
CREATE INDEX IF NOT EXISTS idx_devices_engine_id
    ON devices (engine_id)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_devices_gsm_number
    ON devices (gsm_number)
    WHERE deleted_at IS NULL;

-- Full-text search on device_code + name via trigram
-- (requires pg_trgm; add GIN index for ILIKE/% queries)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_devices_device_code_trgm
    ON devices USING GIN (device_code gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_devices_name_trgm
    ON devices USING GIN (name gin_trgm_ops);

-- last_seen for "recently seen" queries / TTL checks
CREATE INDEX IF NOT EXISTS idx_devices_last_seen
    ON devices (last_seen DESC)
    WHERE deleted_at IS NULL;

-- JSONB metadata GIN index for key-existence and containment queries
CREATE INDEX IF NOT EXISTS idx_devices_metadata
    ON devices USING GIN (metadata)
    WHERE metadata IS NOT NULL;
