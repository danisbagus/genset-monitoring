-- ══════════════════════════════════════════════════════════════════
-- Migration: 000006_create_electrical_telemetry_table (UP)
-- Description: Creates the electrical_telemetry table for storing AC
--              power, voltage, current, and power factor readings.
-- ══════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS electrical_telemetry (
    id              BIGSERIAL       PRIMARY KEY,

    -- Identity & Association
    device_id       UUID            NOT NULL,

    -- Alternator
    charge_alt_volt REAL            NULL, -- Charge Alternator Voltage

    -- General AC
    frequency       REAL            NULL, -- Generator Frequency (Hz)

    -- Phase to Neutral Voltage
    l1_n_volt       REAL            NULL, -- Phase L1-N Voltage
    l2_n_volt       REAL            NULL, -- Phase L2-N Voltage
    l3_n_volt       REAL            NULL, -- Phase L3-N Voltage

    -- Phase to Phase Voltage
    l1_l2_volt      REAL            NULL, -- Phase L1-L2 Voltage
    l2_l3_volt      REAL            NULL, -- Phase L2-L3 Voltage
    l3_l1_volt      REAL            NULL, -- Phase L3-L1 Voltage

    -- Current
    l1_curr         REAL            NULL, -- Phase L1 Current
    l2_curr         REAL            NULL, -- Phase L2 Current
    l3_curr         REAL            NULL, -- Phase L3 Current
    earth_curr      REAL            NULL, -- Earth Current

    -- Apparent Power (VA)
    l1_va           REAL            NULL, -- Phase L1 Apparent Power
    l2_va           REAL            NULL, -- Phase L2 Apparent Power
    l3_va           REAL            NULL, -- Phase L3 Apparent Power
    total_va        REAL            NULL, -- Total Apparent Power

    -- Reactive Power (VAr)
    l1_var          REAL            NULL, -- Phase L1 Reactive Power
    l2_var          REAL            NULL, -- Phase L2 Reactive Power
    l3_var          REAL            NULL, -- Phase L3 Reactive Power
    total_var       REAL            NULL, -- Total Reactive Power

    -- Power Factor
    pf_l1           REAL            NULL, -- Phase L1 Power Factor
    pf_l2           REAL            NULL, -- Phase L2 Power Factor
    pf_l3           REAL            NULL, -- Phase L3 Power Factor
    pf_avg          REAL            NULL, -- Average Power Factor

    -- Percentages
    percent_fv      REAL            NULL, -- Percentage of Full Voltage
    percent_fp      REAL            NULL, -- Percentage of Full Power

    -- Audit
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Foreign Key
    CONSTRAINT fk_electrical_telemetry_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- ── Indexes ────────────────────────────────────────────────────────

-- Index for device association lookups
CREATE INDEX IF NOT EXISTS idx_electrical_telemetry_device_id 
    ON electrical_telemetry (device_id);

-- Index for time-series queries and ordering
CREATE INDEX IF NOT EXISTS idx_electrical_telemetry_created_at 
    ON electrical_telemetry (created_at DESC);

-- Composite index for efficient historical data retrieval per device
CREATE INDEX IF NOT EXISTS idx_electrical_telemetry_device_created
    ON electrical_telemetry (device_id, created_at DESC);
