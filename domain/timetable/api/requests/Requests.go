package requests

type TimetableRequest struct {
	SubjectNumber string `json:"subject_number"`
	StreamNumber  string `json:"stream_number"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}
