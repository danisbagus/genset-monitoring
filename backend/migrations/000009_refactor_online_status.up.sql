-- ══════════════════════════════════════════════════════════════════
-- Migration: 000009_refactor_online_status (UP)
-- Description: Removes is_online boolean columns and adds last_online_at
--              to device_latest_state for computed connectivity.
-- ══════════════════════════════════════════════════════════════════

-- 1. Update device_latest_state
ALTER TABLE device_latest_state 
    DROP COLUMN IF EXISTS device_online,
    ADD COLUMN last_online_at TIMESTAMPTZ NULL;

-- Remove old index
DROP INDEX IF EXISTS idx_device_latest_state_online;

-- Add new index for last_online_at
CREATE INDEX IF NOT EXISTS idx_device_latest_state_last_online
    ON device_latest_state(last_online_at DESC);


-- 2. Update device_statuses
ALTER TABLE device_statuses
    DROP COLUMN IF EXISTS is_online;

-- Remove old index
DROP INDEX IF EXISTS idx_device_statuses_is_online;
