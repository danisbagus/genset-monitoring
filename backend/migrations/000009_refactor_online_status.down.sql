-- ══════════════════════════════════════════════════════════════════
-- Migration: 000009_refactor_online_status (DOWN)
-- Description: Reverts removal of is_online boolean columns.
-- ══════════════════════════════════════════════════════════════════

-- 1. Revert device_latest_state
ALTER TABLE device_latest_state 
    ADD COLUMN device_online BOOLEAN NOT NULL DEFAULT FALSE,
    DROP COLUMN IF EXISTS last_online_at;

-- Re-create old index
CREATE INDEX IF NOT EXISTS idx_device_latest_state_online
    ON device_latest_state(device_online);


-- 2. Revert device_statuses
ALTER TABLE device_statuses
    ADD COLUMN is_online BOOLEAN NOT NULL DEFAULT FALSE;

-- Re-create old index
CREATE INDEX IF NOT EXISTS idx_device_statuses_is_online 
    ON device_statuses (is_online);
