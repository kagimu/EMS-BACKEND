package models

import (
	"time"
)

type (
	Campus struct {
		ID         uint       `gorm:"primaryKey"`
		Number     string     `json:"number" gorm:"size:100;unique;not null"`
		Name       string     `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		Phone      string     `json:"phone" form:"phone" binding:"required" gorm:"size:40;unique;not null"`
		Email      string     `json:"email" form:"email" binding:"required" gorm:"size:40;not null"`
		Location   string     `json:"location" form:"location" binding:"required" gorm:"size:200;not null"`
		Level      string     `json:"level" form:"level" binding:"required" gorm:"size:200;not null"`
		UserNumber string     `json:"user_number" gorm:"size:40;not null"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		DeletedAt  *time.Time `sql:"index"`
	}
)

func (u *Campus) TableName() string {
	return "campuses"
}
