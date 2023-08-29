package models

// UserRole represents the relationship between users and roles
type UserRole struct {
	ID     uint
	UserID string
	RoleID uint
}
