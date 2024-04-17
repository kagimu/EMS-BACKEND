package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	"gorm.io/gorm"
)

type ScopeUserCurriculum struct {
	gorm.Model
	Number           string                     `json:"number" gorm:"size:36;unique;not null"`
	UserNumber       string                     `json:"user_number" gorm:"size:36;not null"`
	User             authModel.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;not null"`
	Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber    string                     `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy          authModel.User             `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
