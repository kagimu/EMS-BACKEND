package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type Scope struct {
	gorm.Model
	Number           string                       `json:"number" gorm:"size:36;unique;not null"`
	UserNumber       string                       `json:"user_number" gorm:"size:36;"`
	User             authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber     string                       `json:"campus_number" gorm:"size:36;"`
	Campus           campusModel.Campus           `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber     string                       `json:"school_number" gorm:"size:36;"`
	School           schoolModel.School           `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CurriculumNumber string                       `json:"curriculum_number" gorm:"size:36;"`
	Curriculum       campusModel.CampusCurriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber      string                       `json:"class_number" gorm:"size:36;"`
	Class            classModel.CampusClass       `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber     string                       `json:"stream_number" gorm:"size:36;"`
	Stream           classModel.Stream            `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedBy          string                       `json:"added_by" gorm:"size:36;not null"`
	AddedByUser      authModel.User               `gorm:"foreignKey:AddedBy;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
