package requests

type InterviewResultRequest struct {
	StudentNumber  string   `json:"student_number"`
	QuestionNumber []string `json:"question_number"`
	Score          []int    `json:"score"`
	Status         string   `json:"status"`
}
