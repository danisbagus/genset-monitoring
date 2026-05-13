package model

import (
	"time"

	"github.com/google/uuid"
)

// EngineTelemetry represents engine performance data and sensor readings.
type EngineTelemetry struct {
	ID                       uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID                 uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	Device                   *Device   `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	Speed                    *int32    `json:"speed"`
	FuelRate                 *float32  `json:"fuel_rate"`
	RatedPower               *float32  `json:"rated_power"`
	RatedSpeed               *int32    `json:"rated_speed"`
	OilFilterOutPressure     *float32  `json:"oil_filter_out_pressure"`
	DesiredOperatingSpeed    *int32    `json:"desired_operating_speed"`
	OilPressure              *float32  `json:"oil_pressure"`
	CoolantTemperature       *float32  `json:"coolant_temperature"`
	TripFuel                 *float32  `json:"trip_fuel"`
	TotalFuel                *float32  `json:"total_fuel"`
	AvgFuelRate              *float32  `json:"avg_fuel_rate"`
	IntakeManifoldPressure   *float32  `json:"intake_manifold_pressure"`
	IntakeManifoldTemperature *float32  `json:"intake_manifold_temperature"`
	KeyswitchBattPotential   *float32  `json:"keyswitch_batt_potential"`
	BattVolt                 *float32  `json:"batt_volt"`
	EcuTemperature           *float32  `json:"ecu_temperature"`
	RunTime                  *int64    `json:"run_time"`
	FuelLevelTop             *float32  `json:"fuel_level_top"`
	FuelLevelBottom          *float32  `json:"fuel_level_bottom"`
	FuelLevelPressure1       *float32  `gorm:"column:fuel_level_pressure_1" json:"fuel_level_pressure_1"`
	FuelLevelPressure2       *float32  `gorm:"column:fuel_level_pressure_2" json:"fuel_level_pressure_2"`
	TurboPressure            *float32  `json:"turbo_pressure"`
	CreatedAt                time.Time `gorm:"not null;default:now();index" json:"created_at"`
}

// TableName overrides the table name used by EngineTelemetry to `engine_telemetry`.
func (EngineTelemetry) TableName() string {
	return "engine_telemetry"
}

// ElectricalTelemetry represents AC power, voltage, current, and power factor readings.
type ElectricalTelemetry struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	Device        *Device   `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	ChargeAltVolt *float32  `json:"charge_alt_volt"`
	Frequency     *float32  `json:"frequency"`
	L1NVolt       *float32  `gorm:"column:l1_n_volt" json:"l1_n_volt"`
	L2NVolt       *float32  `gorm:"column:l2_n_volt" json:"l2_n_volt"`
	L3NVolt       *float32  `gorm:"column:l3_n_volt" json:"l3_n_volt"`
	L1L2Volt      *float32  `json:"l1_l2_volt"`
	L2L3Volt      *float32  `json:"l2_l3_volt"`
	L3L1Volt      *float32  `json:"l3_l1_volt"`
	L1Curr        *float32  `json:"l1_curr"`
	L2Curr        *float32  `json:"l2_curr"`
	L3Curr        *float32  `json:"l3_curr"`
	EarthCurr     *float32  `json:"earth_curr"`
	L1Va          *float32  `json:"l1_va"`
	L2Va          *float32  `json:"l2_va"`
	L3Va          *float32  `json:"l3_va"`
	TotalVa       *float32  `json:"total_va"`
	L1Var         *float32  `json:"l1_var"`
	L2Var         *float32  `json:"l2_var"`
	L3Var         *float32  `json:"l3_var"`
	TotalVar      *float32  `json:"total_var"`
	PfL1          *float32  `json:"pf_l1"`
	PfL2          *float32  `json:"pf_l2"`
	PfL3          *float32  `json:"pf_l3"`
	PfAvg         *float32  `json:"pf_avg"`
	PercentFv     *float32  `json:"percent_fv"`
	PercentFp     *float32  `json:"percent_fp"`
	CreatedAt     time.Time `gorm:"not null;default:now();index" json:"created_at"`
}

// TableName overrides the table name used by ElectricalTelemetry to `electrical_telemetry`.
func (ElectricalTelemetry) TableName() string {
	return "electrical_telemetry"
}
