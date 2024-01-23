package requests

type InterviewQuestionRequest struct {
	InterviewNumber string   `json:"interview_number"`
	Question        string   `json:"question"`
	PossibleAnswers []string `json:"possible_answers"`
	CorrectAnswer   string   `json:"correct_answer"`
	Score           int      `json:"score"`
}
