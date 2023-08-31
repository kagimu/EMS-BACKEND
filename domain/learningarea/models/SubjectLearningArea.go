package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/campus/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/subject/models"
	"time"
)

type (
	SubjectLearningArea struct {
		ID                 uint            `gorm:"primaryKey"`
		Number             string          `json:"number" gorm:"size:100;unique;not null"`
		LearningAreaNumber string          `json:"learning_area_number" gorm:"size:40;not null"`
		SubjectNumber      string          `json:"subject_number" gorm:"size:40;not null"`
		CampusNumber       string          `json:"campus_number" gorm:"size:40;not null"`
		UserNumber         string          `json:"user_number" gorm:"size:40;not null"`
		CreatedAt          time.Time       `json:"created_at"`
		UpdatedAt          time.Time       `json:"updated_at"`
		DeletedAt          *time.Time      `sql:"index"`
		Subject            models3.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		LearningArea       LearningArea    `gorm:"foreignKey:LearningAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Campus             models.Campus   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubjectLearningArea) TableName() string {
	return "subject_learning_areas"
}
