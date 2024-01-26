package requests

type InterviewGradeRequest struct {
	InterviewNumber string   `json:"interview_number"`
	QuestionNumber  []string `json:"question_number"`
	Score           []int    `json:"score"`
}
