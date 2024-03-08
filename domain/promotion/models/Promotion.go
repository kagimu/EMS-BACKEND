package models

import (
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	Number                string                     `json:"number" gorm:"size:36;unique;not null"`
	Condition             string                     `json:"condition" gorm:"type:text"`
	AssessmentNumber      string                     `json:"assessment_number" gorm:"size:36;"`
	ReasonForAltering     string                     `json:"reason_for_altering" gorm:"type:text"`
	GradeNumber           string                     `json:"grade_number" gorm:"size:36;"`
	StudentNumber         string                     `json:"student_number" gorm:"size:36;"`
	AssessmentGradeNumber string                     `json:"assessment_grade_number"  gorm:"size:36"`
	AssessmentGrade       string                     `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Student               studentModel.Student       `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Assessment            assessmentModel.Assessment `gorm:"foreignKey:AssessmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AuthorizedByNumber    string                     `json:"authorized_by_number" gorm:"size:36;"`
	AuthorizedBy          authModel.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber            string                     `json:"user_number" gorm:"size:36;"`
	User                  authModel.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
