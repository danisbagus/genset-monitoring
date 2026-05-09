package service

import (
	"context"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	infradb "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/database"
	inframqtt "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/mqtt"
	infraredis "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/redis"
)

// HealthService defines the contract for health checks.
type HealthService interface {
	Check(ctx context.Context) (*HealthStatus, error)
}

// HealthStatus holds the health status of each dependency.
type HealthStatus struct {
	Postgres string `json:"postgres"`
	Redis    string `json:"redis"`
	MQTT     string `json:"mqtt"`
}

// Healthy returns true if all dependencies are connected.
func (h *HealthStatus) Healthy() bool {
	return h.Postgres == "connected" &&
		h.Redis == "connected" &&
		h.MQTT == "connected"
}

// healthService implements HealthService.
type healthService struct {
	db          *gorm.DB
	redisClient *goredis.Client
	mqttClient  pahomqtt.Client
}

// NewHealthService creates a new HealthService.
func NewHealthService(db *gorm.DB, redisClient *goredis.Client, mqttClient pahomqtt.Client) HealthService {
	return &healthService{
		db:          db,
		redisClient: redisClient,
		mqttClient:  mqttClient,
	}
}

func (s *healthService) Check(_ context.Context) (*HealthStatus, error) {
	status := &HealthStatus{
		Postgres: "disconnected",
		Redis:    "disconnected",
		MQTT:     "disconnected",
	}

	if err := infradb.Ping(s.db); err == nil {
		status.Postgres = "connected"
	}

	if err := infraredis.Ping(s.redisClient); err == nil {
		status.Redis = "connected"
	}

	if inframqtt.IsConnected(s.mqttClient) {
		status.MQTT = "connected"
	}

	return status, nil
}
