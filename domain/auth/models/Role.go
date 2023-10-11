package models

import "gorm.io/gorm"

// Role represents the database model of roles
type Role struct {
	gorm.Model
	ID           uint   `gorm:"primary_key"`
	Name         string `gorm:"unique;not null" json:"name"`
	SchoolNumber string `gorm:"size:36;not null" json:"school_number"`
}
