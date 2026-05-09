package model

// UserRole defines valid user roles.
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleOperator UserRole = "operator"
	RoleViewer   UserRole = "viewer"
)

// User represents a system user.
type User struct {
	Base
	Username     string   `gorm:"uniqueIndex;not null;size:100"  json:"username"`
	Email        string   `gorm:"uniqueIndex;not null;size:255"  json:"email"`
	PasswordHash string   `gorm:"not null"                       json:"-"`
	Role         UserRole `gorm:"not null;default:'viewer'"      json:"role"`
	IsActive     bool     `gorm:"not null;default:true"          json:"is_active"`
}

// TableName overrides the table name used by GORM.
func (User) TableName() string { return "users" }
