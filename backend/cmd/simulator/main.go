package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// EngineTelemetryPayload matches the API request structure
type EngineTelemetryPayload struct {
	Speed              int32   `json:"speed"`
	OilPressure        float32 `json:"oil_pressure"`
	CoolantTemperature float32 `json:"coolant_temperature"`
	FuelRate           float32 `json:"fuel_rate"`
	BattVolt           float32 `json:"batt_volt"`
}

// DeviceState holds the current values for a device to ensure smooth transitions
type DeviceState struct {
	ID                 string
	Speed              float64
	OilPressure        float64
	CoolantTemperature float64
	FuelRate           float64
	BattVolt           float64
}

// EngineTelemetrySimulator handles the simulation of multiple devices
type EngineTelemetrySimulator struct {
	apiUrl    string
	interval  time.Duration
	deviceIDs []string
	token     string
	states    []*DeviceState
	logger    *zap.Logger
	client    *http.Client
}

func NewSimulator(apiUrl string, interval time.Duration, ids []string, token string, logger *zap.Logger) *EngineTelemetrySimulator {
	states := make([]*DeviceState, len(ids))
	for i, id := range ids {
		states[i] = &DeviceState{
			ID:                 id,
			Speed:              1500,
			OilPressure:        4.5,
			CoolantTemperature: 80,
			FuelRate:           25,
			BattVolt:           26,
		}
	}

	return &EngineTelemetrySimulator{
		apiUrl:    apiUrl,
		interval:  interval,
		deviceIDs: ids,
		token:     token,
		states:    states,
		logger:    logger,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (sim *EngineTelemetrySimulator) Run() {
	ticker := time.NewTicker(sim.interval)
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sim.logger.Info("Starting telemetry simulator",
		zap.Int("device_count", len(sim.deviceIDs)),
		zap.Duration("interval", sim.interval),
		zap.String("api_url", sim.apiUrl),
	)

	for {
		select {
		case <-ticker.C:
			for _, state := range sim.states {
				sim.updateState(state)
				sim.sendTelemetry(state)
			}
		case <-stop:
			sim.logger.Info("Simulator shutting down gracefully...")
			return
		}
	}
}

func (sim *EngineTelemetrySimulator) updateState(s *DeviceState) {
	// Smooth randomization: current + small delta
	s.Speed = sim.clamp(s.Speed+(rand.Float64()*40-20), 1400, 1800)
	s.OilPressure = sim.clamp(s.OilPressure+(rand.Float64()*0.4-0.2), 3.0, 6.0)
	s.CoolantTemperature = sim.clamp(s.CoolantTemperature+(rand.Float64()*1.0-0.5), 70, 95)
	s.FuelRate = sim.clamp(s.FuelRate+(rand.Float64()*2.0-1.0), 10, 40)
	s.BattVolt = sim.clamp(s.BattVolt+(rand.Float64()*0.2-0.1), 24, 28)
}

func (sim *EngineTelemetrySimulator) clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (sim *EngineTelemetrySimulator) sendTelemetry(state *DeviceState) {
	url := fmt.Sprintf("%s/devices/%s/engine", sim.apiUrl, state.ID)

	payload := EngineTelemetryPayload{
		Speed:              int32(state.Speed),
		OilPressure:        float32(state.OilPressure),
		CoolantTemperature: float32(state.CoolantTemperature),
		FuelRate:           float32(state.FuelRate),
		BattVolt:           float32(state.BattVolt),
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		sim.logger.Error("Failed to create request", zap.Error(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if sim.token != "" {
		req.Header.Set("Authorization", "Bearer "+sim.token)
	}

	resp, err := sim.client.Do(req)
	if err != nil {
		sim.logger.Warn("Failed to send telemetry", zap.String("device_id", state.ID), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		sim.logger.Warn("API returned non-success status",
			zap.String("device_id", state.ID),
			zap.Int("status", resp.StatusCode))
		return
	}

	sim.logger.Debug("Telemetry sent successfully", zap.String("device_id", state.ID))
}

func main() {
	apiUrl := flag.String("url", "http://localhost:8080/api/v1", "Backend API Base URL")
	interval := flag.Duration("interval", 1*time.Second, "Simulation interval")
	deviceIDs := flag.String("devices", "", "Comma-separated list of Device UUIDs")
	token := flag.String("token", "", "JWT Bearer Token for Authentication")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	if *deviceIDs == "" {
		logger.Fatal("No device IDs provided. Use -devices flag with comma-separated UUIDs")
	}

	ids := strings.Split(*deviceIDs, ",")
	sim := NewSimulator(*apiUrl, *interval, ids, *token, logger)
	sim.Run()
}
