package mqtt

// MQTTService defines the contract for MQTT operations.
type MQTTService interface {
	// PublishTelemetry publishes telemetry data to the MQTT broker.
	// Topic format: genset/<deviceID>/telemetry
	PublishTelemetry(deviceID string, payload []byte) error

	// PublishAlert publishes an alert to the MQTT broker.
	// Topic format: genset/<deviceID>/alert
	PublishAlert(deviceID string, payload []byte) error
}
