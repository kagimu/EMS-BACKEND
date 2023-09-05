package models

import "github.com/jinzhu/gorm"

// RolePermission stores the relationship between roles and permissions
type RolePermission struct {
	gorm.Model
	ID               uint
	RoleNumber       string `gorm:"size:36;not null" json:"role_number"`
	PermissionNumber string `gorm:"size:36;not null" json:"permission_number"`
}
