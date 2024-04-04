package models

import (
	assesmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type PromotionSetting struct {
	gorm.Model
	Number                string                               `json:"number" gorm:"size:36;unique;not null"`
	Name                  string                               `json:"name" gorm:"size:100;not null"`
	Status                string                               `json:"status" gorm:"size:20;"`
	Type                  string                               `json:"type" gorm:"size:50"`
	AssessmentGradeNumber string                               `json:"assessment_grade_number" gorm:"size:36;"`
	AssessmentGrade       assesmentGradeModels.AssessmentGrade `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	NextClassNumber       string                               `json:"next_class_number" gorm:"size:36;not null"`
	NextClass             classModels.CampusClass              `gorm:"foreignKey:NextClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	PreviousClassNumber   string                               `json:"previous_class_number" gorm:"size:36;not null"`
	PreviousClass         classModels.CampusClass              `gorm:"foreignKey:PreviousClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber            string                               `json:"user_number" gorm:"size:36;"`
	User                  authModel.User                       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
