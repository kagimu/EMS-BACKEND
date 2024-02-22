package models

import "gorm.io/gorm"

// Role represents the database model of roles
type Role struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	Name         string `json:"name" gorm:"size:20;not null"`
	UserNumber   string `json:"user_number" gorm:"size:36;index;not null"`
	User         User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber string `gorm:"size:36;not null" json:"school_number"`
}
