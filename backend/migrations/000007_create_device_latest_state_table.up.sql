-- ══════════════════════════════════════════════════════════════════
-- Migration: 000007_create_device_latest_state_table (UP)
-- Description: Creates the device_latest_state table for realtime
--              monitoring snapshots/read models.
-- ══════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS device_latest_state (
    -- Identity
    device_id               UUID            PRIMARY KEY,

    -- Device Status
    device_online           BOOLEAN         NOT NULL DEFAULT FALSE,
    engine_running          BOOLEAN         NOT NULL DEFAULT FALSE,

    -- Engine Metrics
    speed                   INTEGER         NULL,
    coolant_temperature     REAL            NULL,
    oil_pressure            REAL            NULL,

    -- Fuel
    fuel_level              REAL            NULL,

    -- Electrical
    batt_volt               REAL            NULL,
    frequency               REAL            NULL,
    total_va                REAL            NULL,
    pf_avg                  REAL            NULL,

    -- Telemetry Time
    telemetry_recorded_at   TIMESTAMPTZ     NULL,
    last_seen_at            TIMESTAMPTZ     NULL,

    -- Snapshot Metadata
    updated_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Foreign Key
    CONSTRAINT fk_device_latest_state_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- ── Indexes ────────────────────────────────────────────────────────

-- Index for filtering online/offline devices in dashboard
CREATE INDEX IF NOT EXISTS idx_device_latest_state_online
    ON device_latest_state(device_online);

-- Index for filtering running/stopped engines in dashboard
CREATE INDEX IF NOT EXISTS idx_device_latest_state_running
    ON device_latest_state(engine_running);

-- Index for sorting by most recent activity
CREATE INDEX IF NOT EXISTS idx_device_latest_state_last_seen
    ON device_latest_state(last_seen_at DESC);
