package models

import "gorm.io/gorm"

// Tracks and generates the SerialCode in the Product model
// DO NOT DIRECTLY INSERT INTO THIS TABLE
type ProductSerialCode struct {
	gorm.Model
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
}
