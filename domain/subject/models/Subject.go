package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/campus/models"
	"time"
)

type (
	Subject struct {
		ID          uint          `gorm:"primaryKey"`
		Number      string        `json:"number" gorm:"size:100;unique;not null"`
		Name        string        `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		LevelNumber string        `json:"level_number" gorm:"size:40;not null"`
		UserNumber  string        `json:"user_number" gorm:"size:40;not null"`
		CreatedAt   time.Time     `json:"created_at"`
		UpdatedAt   time.Time     `json:"updated_at"`
		DeletedAt   *time.Time    `sql:"index"`
		Level       models2.Level `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User        models.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Subject) TableName() string {
	return "subjects"
}
