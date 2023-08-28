package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/school/models"
	"time"
)

type (
	Campus struct {
		ID           uint           `gorm:"primaryKey"`
		Number       string         `json:"number" gorm:"size:100;unique;not null"`
		Name         string         `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		Phone        string         `json:"phone" form:"phone" binding:"required" gorm:"size:40;unique;not null"`
		Email        string         `json:"email" form:"email" binding:"required" gorm:"size:40;not null"`
		Location     string         `json:"location" form:"location" binding:"required" gorm:"size:200;not null"`
		UserNumber   string         `json:"user_number" gorm:"size:40;not null"`
		SchoolNumber string         `json:"school_number" gorm:"size:40;not null"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		DeletedAt    *time.Time     `sql:"index"`
		Level        []Level        `json:"contacts" gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User         models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		School       models2.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Campus) TableName() string {
	return "campuses"
}
