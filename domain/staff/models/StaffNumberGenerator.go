package models

import "gorm.io/gorm"

// Tracks and generates the Staff Number string in the Staff model
// DO NOT INSERT TO THIS TABLE
type StaffNumberGenerator struct {
	gorm.Model
	ID   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Year string `gorm:"type:varchar(255)" json:"year"`
}
