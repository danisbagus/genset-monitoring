package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/pkg/hashutil"
)

// ── DTOs ─────────────────────────────────────────────────────────

// DeviceListItem is a compact device representation used in list responses.
type DeviceListItem struct {
	ID              uuid.UUID `json:"id"`
	DeviceCode      string    `json:"device_code"`
	Name            string    `json:"name"`
	SerialNumber    string    `json:"serial_number"`
	EngineID        string    `json:"engine_id"`
	GSMNumber       string    `json:"gsm_number"`
	FirmwareVersion string    `json:"firmware_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DeviceDetail is the full device representation for detail responses.
type DeviceDetail struct {
	ID              uuid.UUID   `json:"id"`
	DeviceCode      string      `json:"device_code"`
	Name            string      `json:"name"`
	SerialNumber    string      `json:"serial_number"`
	EngineID        string      `json:"engine_id"`
	GSMNumber       string      `json:"gsm_number"`
	FirmwareVersion string      `json:"firmware_version"`
	Status          string      `json:"status"`
	Metadata        interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// DeviceCreatedOutput is returned after a successful device creation.
// The DeviceToken is returned ONLY once and must not be logged or stored.
type DeviceCreatedOutput struct {
	ID          uuid.UUID `json:"id"`
	DeviceToken string    `json:"device_token"`
}

// DeviceListOutput wraps the list items with pagination metadata.
type DeviceListOutput struct {
	Devices    []*DeviceListItem `json:"devices"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta holds standard pagination metadata.
type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// ── Input structs ─────────────────────────────────────────────────

// DeviceListInput carries the validated query parameters for list.
type DeviceListInput struct {
	Search  string
	Status  string
	Page    int
	Limit   int
	SortBy  string
	SortDir string
}

// CreateDeviceInput is the payload required to create a device.
type CreateDeviceInput struct {
	DeviceCode   string
	Name         string
	SerialNumber string
	EngineID     string
	GSMNumber    string
}

// UpdateDeviceInput carries partial update fields.
// Only non-empty / non-nil fields are applied.
type UpdateDeviceInput struct {
	Name            *string
	EngineID        *string
	GSMNumber       *string
	FirmwareVersion *string
	Status          *model.DeviceLifecycle
}

// ── Interface ─────────────────────────────────────────────────────

// DeviceService defines the contract for device management operations.
type DeviceService interface {
	List(ctx context.Context, input DeviceListInput) (*DeviceListOutput, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DeviceDetail, error)
	Create(ctx context.Context, input CreateDeviceInput) (*DeviceCreatedOutput, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateDeviceInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ── Implementation ────────────────────────────────────────────────

type deviceService struct {
	deviceRepo repository.DeviceRepository
	log        *zap.Logger
}

// NewDeviceService constructs a DeviceService.
func NewDeviceService(deviceRepo repository.DeviceRepository, log *zap.Logger) DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
		log:        log,
	}
}

// List returns a paginated, filtered list of devices.
func (s *deviceService) List(ctx context.Context, input DeviceListInput) (*DeviceListOutput, error) {
	// Clamp / default pagination
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 || input.Limit > 100 {
		input.Limit = 10
	}

	// Allow-list sortable columns to prevent SQL injection via ORDER BY
	allowed := map[string]bool{
		"created_at":       true,
		"updated_at":       true,
		"name":             true,
		"device_code":      true,
		"firmware_version": true,
	}
	sortBy := "created_at"
	if allowed[input.SortBy] {
		sortBy = input.SortBy
	}

	result, err := s.deviceRepo.FindAll(ctx, repository.DeviceFilter{
		Search:  input.Search,
		Status:  input.Status,
		Page:    input.Page,
		Limit:   input.Limit,
		SortBy:  sortBy,
		SortDir: input.SortDir,
	})
	if err != nil {
		return nil, fmt.Errorf("deviceService.List: %w", err)
	}

	items := make([]*DeviceListItem, 0, len(result.Devices))
	for _, d := range result.Devices {
		items = append(items, &DeviceListItem{
			ID:              d.ID,
			DeviceCode:      d.DeviceCode,
			Name:            d.Name,
			SerialNumber:    d.SerialNumber,
			EngineID:        d.EngineID,
			GSMNumber:       d.GSMNumber,
			FirmwareVersion: d.FirmwareVersion,
			CreatedAt:       d.CreatedAt,
			UpdatedAt:       d.UpdatedAt,
		})
	}

	return &DeviceListOutput{
		Devices: items,
		Pagination: PaginationMeta{
			Page:  input.Page,
			Limit: input.Limit,
			Total: result.Total,
		},
	}, nil
}

// GetByID returns the full device detail.
func (s *deviceService) GetByID(ctx context.Context, id uuid.UUID) (*DeviceDetail, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return nil, ErrDeviceNotFound
		}
		return nil, fmt.Errorf("deviceService.GetByID: %w", err)
	}

	return deviceToDetail(device), nil
}

// Create registers a new device.
//   - Validates uniqueness of device_code and serial_number.
//   - Generates a cryptographically secure device token.
//   - Stores only the SHA-256 hash of the token.
//   - Returns the raw token ONCE; it is not retrievable afterwards.
func (s *deviceService) Create(ctx context.Context, input CreateDeviceInput) (*DeviceCreatedOutput, error) {
	// Uniqueness: device_code
	exists, err := s.deviceRepo.ExistsByDeviceCode(ctx, input.DeviceCode, nil)
	if err != nil {
		return nil, fmt.Errorf("deviceService.Create: %w", err)
	}
	if exists {
		return nil, ErrDeviceCodeExists
	}

	// Uniqueness: serial_number
	exists, err = s.deviceRepo.ExistsBySerialNumber(ctx, input.SerialNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("deviceService.Create: %w", err)
	}
	if exists {
		return nil, ErrSerialNumberExists
	}

	// Generate opaque device token (same mechanism as refresh tokens)
	rawToken, tokenHash, err := hashutil.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("deviceService.Create: generate token: %w", err)
	}

	device := &model.Device{
		DeviceCode:        input.DeviceCode,
		Name:              input.Name,
		SerialNumber:      input.SerialNumber,
		EngineID:          input.EngineID,
		GSMNumber:         input.GSMNumber,
		Status:            model.DeviceLifecycleActive,
		HashedDeviceToken: tokenHash,
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("deviceService.Create: %w", err)
	}

	s.log.Info("device created",
		zap.String("device_id", device.ID.String()),
		zap.String("device_code", device.DeviceCode),
	)

	return &DeviceCreatedOutput{
		ID:          device.ID,
		DeviceToken: rawToken,
	}, nil
}

// Update applies a partial update to a device.
func (s *deviceService) Update(ctx context.Context, id uuid.UUID, input UpdateDeviceInput) error {
	// Ensure device exists before building the update map
	if _, err := s.deviceRepo.FindByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return fmt.Errorf("deviceService.Update: %w", err)
	}

	fields := make(map[string]interface{})
	if input.Name != nil && *input.Name != "" {
		fields["name"] = *input.Name
	}
	if input.EngineID != nil {
		fields["engine_id"] = *input.EngineID
	}
	if input.GSMNumber != nil {
		fields["gsm_number"] = *input.GSMNumber
	}
	if input.FirmwareVersion != nil {
		fields["firmware_version"] = *input.FirmwareVersion
	}
	if input.Status != nil {
		fields["status"] = string(*input.Status)
	}

	if len(fields) == 0 {
		return nil // nothing to update – idempotent
	}

	// GORM auto-sets updated_at when using Updates(map)
	if err := s.deviceRepo.Update(ctx, id, fields); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return fmt.Errorf("deviceService.Update: %w", err)
	}

	s.log.Info("device updated",
		zap.String("device_id", id.String()),
		zap.Any("fields", fields),
	)
	return nil
}

// Delete soft-deletes a device, preserving the audit trail.
func (s *deviceService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.deviceRepo.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			return ErrDeviceNotFound
		}
		return fmt.Errorf("deviceService.Delete: %w", err)
	}

	s.log.Info("device soft-deleted", zap.String("device_id", id.String()))
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────

func deviceToDetail(d *model.Device) *DeviceDetail {
	detail := &DeviceDetail{
		ID:              d.ID,
		DeviceCode:      d.DeviceCode,
		Name:            d.Name,
		SerialNumber:    d.SerialNumber,
		EngineID:        d.EngineID,
		GSMNumber:       d.GSMNumber,
		FirmwareVersion: d.FirmwareVersion,
		Status:          string(d.Status),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	// Deserialise JSONB metadata; ignore if empty / invalid JSON
	if len(d.Metadata) > 0 {
		var m interface{}
		if err := json.Unmarshal(d.Metadata, &m); err == nil {
			detail.Metadata = m
		}
	}

	return detail
}

// ── Sentinel errors ───────────────────────────────────────────────

var (
	ErrDeviceNotFound     = errors.New("device not found")
	ErrDeviceCodeExists   = errors.New("device_code already in use")
	ErrSerialNumberExists = errors.New("serial_number already in use")
)
