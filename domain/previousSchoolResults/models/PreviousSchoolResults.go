package models

import (
	applicationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type PreviousSchoolResults struct {
	gorm.Model
	Number               string                              `json:"number" gorm:"size:36;unique;not null"`
	AssessmentAreaNumber string                              `json:"assessment_area_number" gorm:"size:36; not null"`
	AssessmentArea       assessmentAreaModels.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Score                float64                             `json:"score" gorm:"not null"`
	StudentNumber        string                              `json:"student_number" gorm:"size:36;not null"`
	Student              models.Student                      `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Attachment           string                              `json:"attachment" form:"attachment"`
	ApplicationNumber    string                              `json:"application_number" gorm:"size:36;not null"`
	Application          applicationModels.Application       `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber           string                              `json:"user_number" gorm:"size:36;not null"`
	User                 authModels.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
