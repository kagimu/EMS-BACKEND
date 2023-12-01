package models

import (
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers/models"
	"time"
)

type (
	SchoolCurriculum struct {
		ID               uint                       `gorm:"primaryKey"`
		Number           string                     `json:"number" gorm:"size:36;unique;not null"`
		CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;index;not null"`
		SchoolNumber     string                     `json:"school_number" gorm:"size:36;index;not null"`
		CreatedAt        time.Time                  `json:"created_at"`
		UpdatedAt        time.Time                  `json:"updated_at"`
		DeletedAt        *time.Time                 `sql:"index"`
		Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		School           School                     `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SchoolCurriculum) TableName() string {
	return "school_curriculums"
}
