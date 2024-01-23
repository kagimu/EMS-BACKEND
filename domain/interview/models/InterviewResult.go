package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

type InterviewResult struct {
	Number         string               `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber  string               `json:"student_number" gorm:"size:36;not null"`
	Student        studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	QuestionNumber string               `json:"question_number" gorm:"size:36;not null"`
	Score          int                  `json:"score"`
	Status         string               `gorm:"type:varchar(255);not null;default:'DRAFT'" json:"status"`
	UserNumber     string               `json:"userNumber" gorm:"size:36;not null"`
	User           models.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
