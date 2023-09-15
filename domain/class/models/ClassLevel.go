package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/campus/models"
	"time"
)

type (
	ClassLevel struct {
		ID           uint           `gorm:"primaryKey"`
		Number       string         `json:"number" gorm:"size:100;unique;not null"`
		ClassNumber  string         `json:"class_id" gorm:"size:40;not null"`
		CampusNumber string         `json:"campus_id" gorm:"size:40;not null"`
		LevelNumber  string         `json:"level_number" gorm:"size:40;not null"`
		UserNumber   string         `json:"user_number" gorm:"size:40;not null"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		DeletedAt    *time.Time     `sql:"index"`
		Class        Class          `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Campus       models2.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Level        models2.Level  `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User         models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	/*ClassInfo struct {
		Number   string         `json:"number"`
		Name     string         `json:"name"`
		Subjects []SubjectsInfo `json:"subjects"`
	}

	SubjectsInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}*/
)

func (u *ClassLevel) TableName() string {
	return "classes_levels"
}
