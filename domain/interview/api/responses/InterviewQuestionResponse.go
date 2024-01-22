package responses

type InterviewQuestionResponse struct {
	InterviewNumber string
	Question        string
	PossibleAnswers string
	CorrectAnswer   string
	Score           int
}
