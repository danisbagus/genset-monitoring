-- ══════════════════════════════════════════════════════════════════
-- Migration: 000010_create_alerts_table (DOWN)
-- Description: Drops the alerts table and associated indexes.
-- ══════════════════════════════════════════════════════════════════

-- Drop indexes first (optional but good practice)
DROP INDEX IF EXISTS idx_alerts_device_active;
DROP INDEX IF EXISTS idx_alerts_created_at;
DROP INDEX IF EXISTS idx_alerts_severity;
DROP INDEX IF EXISTS idx_alerts_status;
DROP INDEX IF EXISTS idx_alerts_device_id;

-- Drop table
DROP TABLE IF EXISTS alerts;
