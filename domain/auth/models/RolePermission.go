package models

import "gorm.io/gorm"

// RolePermission stores the relationship between roles and permissions
type RolePermission struct {
	gorm.Model
	ID         uint
	RoleNumber string `gorm:"size:36;not null" json:"role_number"`
	Permission string `gorm:"size:36;not null" json:"permission"`
	Role       Role   `gorm:"foreignKey:RoleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
