package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type InterviewQuestion struct {
	gorm.Model
	Number          string      `json:"number" gorm:"size:36;unique;not null"`
	InterviewNumber string      `json:"interview_number" gorm:"size:36;not null"`
	Question        string      `json:"question" gorm:"type:text;not null"`
	CorrectAnswer   string      `json:"correct_answer" gorm:"type:text;not null"`
	Score           int         `json:"score" gorm:"not null"`
	PossibleAnswers string      `json:"possible_answers" gorm:"type:text;not null"`
	UserNumber      string      `json:"user_number" gorm:"size:36;not null"`
	User            models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type Answer struct {
	gorm.Model
	Number         string `json:"number" gorm:"size:36;unique;not null"`
	QuestionNumber string `json:"question_number" gorm:"size:36;not null"`
	Body           string `json:"body" gorm:"type:text;not null"`
}
