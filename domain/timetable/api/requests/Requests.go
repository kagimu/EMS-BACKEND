package requests

type TimetableRequest struct {
	SubjectNumber string `json:"subject_number"`
	StreamNumber  string `json:"stream_number"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Day           string `json:"day"`
	CampusNumber  string `json:"campus_number"`
}

type StudyTimetableRequest struct {
	StudyPeriodNumber string `json:"study_period_number"`
	Name              string `json:"name"`
}
