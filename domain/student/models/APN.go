package models

import "gorm.io/gorm"

type APN struct {
	gorm.Model
	ID    uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Year  string `gorm:"type:varchar(255)" json:"year"`
	Month string `gorm:"type:varchar(255)" json:"month"`
}
