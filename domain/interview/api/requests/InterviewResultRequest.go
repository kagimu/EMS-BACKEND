package requests

type InterviewResultRequest struct {
	StudentNumber   string   `json:"student_number"`
	QuestionNumbers []string `json:"question_numbers"`
	Scores          []int    `json:"scores"`
	Status          string   `json:"status"`
}
