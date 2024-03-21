package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type ProgressReport struct {
	gorm.Model
	Number            string                       `json:"number" gorm:"size:36;unique;not null"`
	StudyPeriodNumber string                       `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod       studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassTeacher      string                       `json:"class_teacher" gorm:"type:text;"`
	SubjectTeacher    string                       `json:"subject_teacher" gorm:"type:text;"`
	Comment           string                       `json:"comment" gorm:"type:text;"`
	StudentNumber     string                       `json:"student_number" gorm:"size:36;"`
	Student           studentModel.Student         `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                       `json:"user_number" gorm:"size:36;index;not null"`
	User              authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
