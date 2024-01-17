package models

type InterviewQuestion struct {
	Number          string `json:"number" gorm:"size:36;unique;not null"`
	InterviewNumber string `json:"interview_number" gorm:"size:36;not null"`
	Question        string `json:"question" gorm:"type:text;not null"`
	CorrectAnswer   string `json:"correct_answer" gorm:"type:text;not null"`
	Score           string `json:"score" gorm:"not null"`
}

type Answer struct {
	Number                  string `json:"number" gorm:"size:36;unique;not null"`
	InterviewQuestionNumber string `json:"interview_question_number" gorm:"size:36;not null"`
	Body                    string `json:"body" gorm:"type:text;not null"`
}
