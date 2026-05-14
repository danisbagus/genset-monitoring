-- ══════════════════════════════════════════════════════════════════
-- Migration: 000007_create_device_latest_state_table (DOWN)
-- Description: Drops the device_latest_state table and its indexes.
-- ══════════════════════════════════════════════════════════════════

DROP INDEX IF EXISTS idx_device_latest_state_last_seen;
DROP INDEX IF EXISTS idx_device_latest_state_running;
DROP INDEX IF EXISTS idx_device_latest_state_online;

DROP TABLE IF EXISTS device_latest_state;
