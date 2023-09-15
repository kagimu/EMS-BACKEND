package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"time"
)

type (
	LearningArea struct {
		ID         uint        `gorm:"primaryKey"`
		Number     string      `json:"number" gorm:"size:100;unique;not null"`
		Name       string      `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		Code       string      `json:"code" form:"code" binding:"required" gorm:"size:60;unique;not null"`
		UserNumber string      `json:"user_number" gorm:"size:40;not null"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		DeletedAt  *time.Time  `sql:"index"`
		User       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	LearningAreaInfo struct {
		Number   string         `json:"number"`
		Name     string         `json:"name"`
		Subjects []SubjectsInfo `json:"subjects"`
	}

	SubjectsInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)

func (u *LearningArea) TableName() string {
	return "learning_areas"
}
