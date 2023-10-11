package requests

type TimetableRequest struct {
	SubjectNumber string `json:"subject_number"`
	StreamNumber  string `json:"stream_number"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Day           string `json:"day"`
	CampusNumber  string `json:"campus_number"`
}

type CalendarEventRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type GeneralCalendarRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Name      string `json:"name"`
}
