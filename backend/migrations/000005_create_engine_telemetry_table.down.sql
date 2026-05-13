-- ══════════════════════════════════════════════════════════════════
-- Migration: 000005_create_engine_telemetry_table (DOWN)
-- Description: Drops the engine_telemetry table and its indexes.
-- ══════════════════════════════════════════════════════════════════

DROP INDEX IF EXISTS idx_engine_telemetry_device_created;
DROP INDEX IF EXISTS idx_engine_telemetry_created_at;
DROP INDEX IF EXISTS idx_engine_telemetry_device_id;
DROP TABLE IF EXISTS engine_telemetry;
