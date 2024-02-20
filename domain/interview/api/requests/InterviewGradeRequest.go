package requests

type InterviewGradeRequest struct {
	InterviewNumber string   `json:"interview_number"`
	QuestionNumbers []string `json:"question_numbers"`
	Scores          []int    `json:"scores"`
}
