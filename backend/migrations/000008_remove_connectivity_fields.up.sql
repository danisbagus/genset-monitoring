-- ══════════════════════════════════════════════════════════════════
-- Migration: 000008_remove_connectivity_fields (UP)
-- Description: Removes redundant connectivity fields from devices and
--              device_statuses tables. These will be handled via
--              realtime presence in the future.
-- ══════════════════════════════════════════════════════════════════

-- Remove from devices table
DROP INDEX IF EXISTS idx_devices_last_seen;
DROP INDEX IF EXISTS idx_devices_is_online;

ALTER TABLE devices
    DROP COLUMN IF EXISTS is_online,
    DROP COLUMN IF EXISTS last_seen;