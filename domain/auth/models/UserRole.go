package models

// UserRole represents the relationship between users and roles
type UserRole struct {
	ID         uint   `gorm:"primary_key"`
	UserNumber string `json:"user_number" gorm:"size:36;not null"`
	RoleNumber string `json:"role_number" gorm:"size:36;not null"`
}
