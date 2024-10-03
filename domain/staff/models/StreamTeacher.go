package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	streamModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type StreamTeacher struct {
	ID                   uint                                `json:"id" gorm:"primaryKey;autoIncrement"`
	Number               string                              `json:"number" gorm:"size:36;unique;not null"`
	CampusClassNumber    string                              `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass          streamModels.CampusClass            `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber         string                              `json:"stream_number" gorm:"size:36;not null"`
	Stream               streamModels.Stream                 `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber           string                              `json:"user_number" gorm:"size:36;not null"`
	User                 models.User                         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedBy              string                              `json:"added_by" gorm:"size:36;not null"`
	AssessmentAreaNumber string                              `json:"assessment_area_number" gorm:"size:36;not null"`
	AssessmentArea       assessmentAreaModels.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber         string                              `json:"school_number" gorm:"size:36;not null"`
	School               schoolModel.School                  `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CurriculumNumber     string                              `json:"curriculum_number" gorm:"size:36;not null"`
	Curriculum           curriculumModel.Curriculum          `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber         string                              `json:"campus_number" gorm:"size:36;not null"`
	Campus               campusModel.Campus                  `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CreatedAt            time.Time                           `json:"created_at"`
	UpdatedAt            time.Time                           `json:"updated_at"`
}
