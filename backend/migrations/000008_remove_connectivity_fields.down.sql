-- ══════════════════════════════════════════════════════════════════
-- Migration: 000008_remove_connectivity_fields (DOWN)
-- Description: Restores connectivity fields to devices and
--              device_statuses tables.
-- ══════════════════════════════════════════════════════════════════

-- Restore to devices table
ALTER TABLE devices
    ADD COLUMN is_online BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN last_seen  TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_devices_is_online ON devices (is_online);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices (last_seen DESC);