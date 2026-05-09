package model

// User represents a system user.
type User struct {
	Base
	Name     string `gorm:"not null"          json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null"          json:"-"`
	Role     string `gorm:"default:'viewer'"  json:"role"`
	IsActive bool   `gorm:"default:true"      json:"is_active"`
}
