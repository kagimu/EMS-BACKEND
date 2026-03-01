package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Metric struct {
	gorm.Model
	Number           string `json:"number" gorm:"size:36;unique;not null"`
	Name             string `json:"name"`
	Code             int32  `json:"code"`
	CurriculumNumber string `json:"curriculum_number" gorm:"size:36;not null"`
	//ExamType         string         `json:"exam_type"`
	UserNumber string         `json:"user_number" gorm:"size:36;index;not null"`
	User       authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
