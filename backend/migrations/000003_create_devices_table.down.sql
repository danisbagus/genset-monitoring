-- ══════════════════════════════════════════════════════════════════
-- Migration: 000003_create_devices_table (DOWN)
-- Description: Rolls back the devices table and associated objects
-- ══════════════════════════════════════════════════════════════════

-- Drop indexes (most are dropped implicitly with the table,
-- but drop explicitly for clarity and to avoid dependency issues)
DROP INDEX IF EXISTS idx_devices_metadata;
DROP INDEX IF EXISTS idx_devices_last_seen;
DROP INDEX IF EXISTS idx_devices_name_trgm;
DROP INDEX IF EXISTS idx_devices_device_code_trgm;
DROP INDEX IF EXISTS idx_devices_status;
DROP INDEX IF EXISTS idx_devices_is_online;
DROP INDEX IF EXISTS idx_devices_deleted_at;
DROP INDEX IF EXISTS idx_devices_serial_number_active;
DROP INDEX IF EXISTS idx_devices_device_code_active;

-- Drop table (CASCADE removes all dependent constraints/indexes)
DROP TABLE IF EXISTS devices CASCADE;

-- Drop custom enum type
DROP TYPE IF EXISTS device_status;
