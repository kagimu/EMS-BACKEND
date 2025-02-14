package responses

type AttendanceResponse struct {
	Number        string `json:"number"`
	Name          string `json:"user_number"`
	Date          string `json:"date"`
	CheckInTime   string `json:"check_in_time"`
	CheckOutTime  string `json:"check_out_time"`
	Present       string `json:"present"`
	Class         string `json:"class"`
	School        string `json:"school"`
	Campus        string `json:"campus"`
	AbsenceReason string `json:"absence_reason"`
}
