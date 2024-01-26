package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type InterviewGrade struct {
	gorm.Model
	Number          string         `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber   string         `json:"student_number" gorm:"size:36"`
	InterviewNumber string         `json:"interview_number" gorm:"size:36"`
	QuestionNumber  string         `json:"question_number" gorm:"size:36"`
	Score           int            `json:"score"`
	UserNumber      string         `json:"userNumber" gorm:"size:36;not null"`
	User            authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
