-- ══════════════════════════════════════════════════════════════════
-- Migration: 000004_create_device_statuses_table (DOWN)
-- Description: Drops the device_statuses table
-- ══════════════════════════════════════════════════════════════════

DROP INDEX IF EXISTS idx_device_statuses_last_seen;
DROP INDEX IF EXISTS idx_device_statuses_is_online;
DROP TABLE IF EXISTS device_statuses;
