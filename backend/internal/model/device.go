package model

// Device represents a physical genset device.
type Device struct {
	Base
	Name     string `gorm:"not null"           json:"name"`
	Serial   string `gorm:"uniqueIndex;not null" json:"serial"`
	Location string `json:"location"`
	IsActive bool   `gorm:"default:true"       json:"is_active"`
}
