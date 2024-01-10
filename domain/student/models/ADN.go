package models

import "gorm.io/gorm"

// Tracks and generates the ADN string in the Student model
// DO NOT INSERT TO THIS TABLE
type ADN struct {
	gorm.Model
	ID    uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Year  string `gorm:"type:varchar(255)" json:"year"`
	Month string `gorm:"type:varchar(255)" json:"month"`
}
