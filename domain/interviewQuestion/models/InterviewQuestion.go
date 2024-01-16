package models

type InterviewQuestion struct {
	Number          string `json:"number" gorm:"size:36;unique;not null"`
	InterviewNumber string `json:"interview_number" gorm:"size:36;not null"`
	Question        string `json:"question" gorm:"type:text;not null"`
	PossibleAnswers string `json:"possible_answers" gorm:"type:text;not null"`
	CorrectAnswer   string `json:"correct_answer" gorm:"type:text;not null"`
	Score           string `json:"score" gorm:"not null"`
}
