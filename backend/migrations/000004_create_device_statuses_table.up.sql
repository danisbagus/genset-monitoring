-- ══════════════════════════════════════════════════════════════════
-- Migration: 000004_create_device_statuses_table (UP)
-- Description: Creates the device_statuses table for frequent heartbeat
--              updates and connectivity tracking.
-- ══════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS device_statuses (
    device_id           UUID            PRIMARY KEY,

    -- Connectivity
    is_online           BOOLEAN         NOT NULL DEFAULT FALSE,
    last_seen           TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Signal & Hardware health
    gsm_signal          SMALLINT        NULL, -- in dBm (e.g., -70)
    gps_connected       BOOLEAN         NOT NULL DEFAULT FALSE,
    server_connected    BOOLEAN         NOT NULL DEFAULT FALSE,
    can_connected       BOOLEAN         NOT NULL DEFAULT FALSE,
    rs485_connected     BOOLEAN         NOT NULL DEFAULT FALSE,
    sd_card_ok          BOOLEAN         NOT NULL DEFAULT FALSE,

    -- Audit
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Foreign Key
    CONSTRAINT fk_device_statuses_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- ── Indexes ────────────────────────────────────────────────────────
-- Index for online/offline filtering (useful for dashboard aggregations)
CREATE INDEX IF NOT EXISTS idx_device_statuses_is_online 
    ON device_statuses (is_online);

-- Index for last_seen to find stale/offline devices quickly
CREATE INDEX IF NOT EXISTS idx_device_statuses_last_seen 
    ON device_statuses (last_seen DESC);
