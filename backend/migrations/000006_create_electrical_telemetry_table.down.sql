-- ══════════════════════════════════════════════════════════════════
-- Migration: 000006_create_electrical_telemetry_table (DOWN)
-- Description: Drops the electrical_telemetry table and its indexes.
-- ══════════════════════════════════════════════════════════════════

DROP INDEX IF EXISTS idx_electrical_telemetry_device_created;
DROP INDEX IF EXISTS idx_electrical_telemetry_created_at;
DROP INDEX IF EXISTS idx_electrical_telemetry_device_id;
DROP TABLE IF EXISTS electrical_telemetry;
