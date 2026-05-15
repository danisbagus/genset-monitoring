-- ══════════════════════════════════════════════════════════════════
-- Migration: 000010_create_alerts_table (UP)
-- Description: Creates the alerts table to track system and telemetry 
--              threshold violations with acknowledgment workflow.
-- ══════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS alerts (
    -- Primary key (UUID v4)
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relations
    device_id           UUID NOT NULL,

    -- Alert classification
    type                VARCHAR(50) NOT NULL,    -- e.g., 'engine', 'electrical', 'connectivity'
    severity            VARCHAR(20) NOT NULL,    -- e.g., 'critical', 'warning', 'info'

    -- Content
    title               VARCHAR(255) NOT NULL,
    message             TEXT NOT NULL,

    -- Contextual telemetry (optional, for threshold-based alerts)
    metric_name         VARCHAR(100),
    metric_value        DOUBLE PRECISION,
    threshold_value     DOUBLE PRECISION,

    -- Lifecycle status
    status              VARCHAR(20) NOT NULL DEFAULT 'active',

    -- Acknowledgment tracking
    acknowledged_at     TIMESTAMPTZ,
    acknowledged_by     UUID,

    -- Resolution tracking
    resolved_at         TIMESTAMPTZ,

    -- Audit
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_alerts_device 
        FOREIGN KEY (device_id) 
        REFERENCES devices(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_alerts_acknowledged_by 
        FOREIGN KEY (acknowledged_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL,

    CONSTRAINT chk_alerts_status 
        CHECK (status IN ('active', 'acknowledged', 'resolved')),
        
    CONSTRAINT chk_alerts_severity 
        CHECK (severity IN ('critical', 'warning', 'info'))
);

-- ── Indexes ────────────────────────────────────────────────────────

-- Efficient filtering by device
CREATE INDEX IF NOT EXISTS idx_alerts_device_id 
    ON alerts(device_id);

-- Filter/sort by status and severity (common dashboard queries)
CREATE INDEX IF NOT EXISTS idx_alerts_status 
    ON alerts(status);

CREATE INDEX IF NOT EXISTS idx_alerts_severity 
    ON alerts(severity);

-- Ordering by creation time (latest first)
CREATE INDEX IF NOT EXISTS idx_alerts_created_at 
    ON alerts(created_at DESC);

-- Composite index for active alerts per device
CREATE INDEX IF NOT EXISTS idx_alerts_device_active 
    ON alerts(device_id) 
    WHERE status = 'active';
