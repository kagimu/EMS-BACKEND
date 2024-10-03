package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"time"
)

type Billable struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Number            string                       `json:"number" gorm:"size:36;unique;not null"`
	Description       string                       `json:"description" gorm:"size:2000"`
	BillType          string                       `json:"bill_type" gorm:"size:36;not null;default:CONTINUING"`
	Status            string                       `json:"status" gorm:"size:36;not null;default:ACTIVE"`
	UserNumber        string                       `json:"user_number" gorm:"size:36;not null"`
	User              authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber string                       `json:"study_period_number" gorm:"size:36;not null"`
	StudyPeriod       studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber       string                       `json:"class_number" gorm:"size:36;not null"`
	Class             classModels.Class            `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CampusNumber      string                       `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusModel.Campus           `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ModeOfStudyNumber string                       `json:"mode_of_study_number" gorm:"size:36;not null"`
	ModeOfStudy       models.ModeOfStudy           `gorm:"foreignKey:ModeOfStudyNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	SchoolNumber      string                       `json:"school_number" gorm:"size:36;not null"`
	School            models.School                `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CurriculumNumber  string                       `json:"curriculum_number" gorm:"size:36;not null"`
	Curriculum        models2.Curriculum           `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (billable *Billable) TableName() string {
	return "billables"
}
