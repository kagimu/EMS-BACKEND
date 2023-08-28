package models

import (
	"time"
)

type (
	Level struct {
		ID           uint       `gorm:"primaryKey"`
		Number       string     `json:"number" gorm:"size:100;unique;not null"`
		Name         string     `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		CampusNumber string     `json:"campus-number_number" gorm:"size:40;not null"`
		CreatedAt    time.Time  `json:"created_at"`
		UpdatedAt    time.Time  `json:"updated_at"`
		DeletedAt    *time.Time `sql:"index"`
		Campus       Campus     `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Level) TableName() string {
	return "levels"
}
