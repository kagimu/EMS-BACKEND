package models

import (
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	Number                    string                        `json:"number" gorm:"size:36;unique;not null"`
	Condition                 string                        `json:"condition" gorm:"type:text"`
	AssessmentNumber          string                        `json:"assessment_number" gorm:"size:36;"`
	ReasonForAltering         string                        `json:"reason_for_altering" gorm:"type:text"`
	PreviousCampusClassNumber string                        `json:"previous_campus_class_number" gorm:"size:36;"`
	PreviousCampusClass       campusClassModels.CampusClass `gorm:"foreignKey:PreviousCampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	NextCampusClassNumber     string                        `json:"next_campus_class_number" gorm:"size:36;"`
	NextCampusClass           campusClassModels.CampusClass `gorm:"foreignKey:NextCampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status                    string                        `json:"status" gorm:"size:36;"`
	Type                      string                        `json:"type" gorm:"size:20;"`
	AssessmentGradeNumber     string                        `json:"assessment_grade_number"  gorm:"size:36"`
	AssessmentGrade           string                        `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber             string                        `json:"student_number" gorm:"size:36;"`
	Student                   studentModel.Student          `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Assessment                assessmentModel.Assessment    `gorm:"foreignKey:AssessmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Authorizer                string                        `json:"authorizer" gorm:"type:text"`
	UserNumber                string                        `json:"user_number" gorm:"size:36;"`
	User                      authModel.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
