package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	gradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

type PreAdmissionInterviewTest struct {
	gorm.Model
	Number        string               `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber string               `json:"student_number" gorm:"size:36;not null"`
	SubjectNumber string               `json:"subject_number" gorm:"size:36;not null"`
	Subject       subjectModel.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Marks         string               `json:"marks" gorm:"size:36;not null"`
	GradeNumber   string               `json:"grade_number" gorm:"size:36;not null"`
	Grade         gradeModel.Grade     `gorm:"foreignKey:GradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber   string               `json:"preferred_class" gorm:"size:36;not null"`
	Class         classModel.Class     `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Student       models.Student       `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber    string               `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (p *PreAdmissionInterviewTest) TableName() string {
	return "pre_interview_test"
}
