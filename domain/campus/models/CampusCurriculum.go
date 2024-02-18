package models

import (
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type (
	CampusCurriculum struct {
		gorm.Model
		Number                 string                   `json:"number" gorm:"size:36;unique;not null"`
		CampusNumber           string                   `json:"campus_number" gorm:"size:36;not null"`
		SchoolCurriculumNumber string                   `json:"school_curriculum_number" gorm:"size:36;not null"`
		Campus                 Campus                   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		SchoolCurriculum       models2.SchoolCurriculum `gorm:"foreignKey:SchoolCurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
