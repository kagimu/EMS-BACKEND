package models

import (
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers/models"
	"gorm.io/gorm"
)

type (
	SchoolCurriculum struct {
		gorm.Model
		Number           string                     `json:"number" gorm:"size:36;unique;not null"`
		CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;index;not null"`
		SchoolNumber     string                     `json:"school_number" gorm:"size:36;index;not null"`
		Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		School           School                     `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
