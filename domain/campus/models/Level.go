package models

import (
	"time"
)

type (
	Level struct {
		ID           uint       `gorm:"primaryKey"`
		Number       string     `json:"number" gorm:"size:36;unique;not null"`
		Name         string     `json:"name" form:"name" binding:"required" gorm:"size:2000;"`
		CampusNumber string     `json:"campus_number" gorm:"size:36;not null"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `sql:"index"`
	}
)

func (u *Level) TableName() string {
	return "levels"
}
