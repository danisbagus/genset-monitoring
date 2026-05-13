-- ══════════════════════════════════════════════════════════════════
-- Migration: 000005_create_engine_telemetry_table (UP)
-- Description: Creates the engine_telemetry table for storing detailed
--              engine performance data and sensor readings.
-- ══════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS engine_telemetry (
    id                          BIGSERIAL       PRIMARY KEY,

    -- Identity & Association
    device_id                   UUID            NOT NULL,

    -- Engine Performance
    speed                       INTEGER         NULL, -- Engine Speed (RPM)
    fuel_rate                   REAL            NULL, -- Engine Fuel Rate
    rated_power                 REAL            NULL, -- Engine Rated Power
    rated_speed                 INTEGER         NULL, -- Engine Rated Speed

    -- Oil & Pressure
    oil_filter_out_pressure     REAL            NULL, -- Engine Oil Filter Outlet Pressure
    desired_operating_speed     INTEGER         NULL, -- Engine Desired Operating Speed
    oil_pressure                REAL            NULL, -- Engine Oil Pressure

    -- Fuel & Fluids
    coolant_temperature         REAL            NULL, -- Engine Coolant Temperature
    trip_fuel                   REAL            NULL, -- Engine Trip Fuel
    total_fuel                  REAL            NULL, -- Engine Total Fuel Used
    avg_fuel_rate               REAL            NULL, -- Engine Average Fuel Economy

    -- Air & Intake
    intake_manifold_pressure    REAL            NULL, -- Engine Intake Manifold 1 Pressure
    intake_manifold_temperature REAL            NULL, -- Engine Intake Manifold 1 Temperature

    -- Electrical
    keyswitch_batt_potential    REAL            NULL, -- Engine Keyswitch Battery Potential
    batt_volt                   REAL            NULL, -- Electrical Potential (Voltage)
    ecu_temperature             REAL            NULL, -- Engine ECU Temperature

    -- Operational
    run_time                    BIGINT          NULL, -- Engine Total Hours of Operation

    -- Fuel Levels (Secondary/Tank)
    fuel_level_top              REAL            NULL,
    fuel_level_bottom           REAL            NULL,
    fuel_level_pressure_1       REAL            NULL,
    fuel_level_pressure_2       REAL            NULL,

    -- Forced Induction
    turbo_pressure              REAL            NULL, -- Engine Turbocharger 1 Boost Pressure

    -- Audit
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    -- Foreign Key
    CONSTRAINT fk_engine_telemetry_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
        ON DELETE CASCADE
);

-- ── Indexes ────────────────────────────────────────────────────────

-- Index for device association lookups
CREATE INDEX IF NOT EXISTS idx_engine_telemetry_device_id 
    ON engine_telemetry (device_id);

-- Index for time-series queries and ordering
CREATE INDEX IF NOT EXISTS idx_engine_telemetry_created_at 
    ON engine_telemetry (created_at DESC);

-- Composite index for efficient historical data retrieval per device
CREATE INDEX IF NOT EXISTS idx_engine_telemetry_device_created
    ON engine_telemetry (device_id, created_at DESC);