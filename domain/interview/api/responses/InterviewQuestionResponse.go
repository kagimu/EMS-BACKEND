package responses

type InterviewQuestionResponse struct {
	InterviewNumber string
	Question        string
	PossibleAnswers string
	CorrectAnswer   string
	Score           int
}

type StudentInterviewResults struct {
	Name    string
	Score   float64
	Grade   string
	Comment string
}
