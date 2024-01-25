package requests

import "errors"

type InterviewQuestionRequest struct {
	InterviewNumber string   `json:"interview_number"`
	Question        string   `json:"question"`
	PossibleAnswers []string `json:"possible_answers"`
	CorrectAnswer   string   `json:"correct_answer"`
	Score           int64    `json:"score"`
}

func (i *InterviewQuestionRequest) Validate() error {
	if i.InterviewNumber == "" {
		return errors.New("venue is required")
	}
	if i.Question == "" {
		return errors.New("question is required")
	}
	if len(i.PossibleAnswers) == 0 {
		return errors.New("possible answers are required")
	}

	return nil
}
