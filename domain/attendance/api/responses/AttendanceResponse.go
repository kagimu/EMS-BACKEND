package responses

type AttendanceResponse struct {
	Number            string `json:"number"`
	Name              string `json:"name"`
	UserNumber        string `json:"user_number"`
	StudentNumber     string `json:"student_number"`
	Date              string `json:"date"`
	CheckInTime       string `json:"check_in_time"`
	CheckOutTime      string `json:"check_out_time"`
	Present           string `json:"present"`
	Class             string `json:"class"`
	School            string `json:"school"`
	Campus            string `json:"campus"`
	AbsenceReason     string `json:"absence_reason"`
	CampusNumber      string `json:"campus_number"`
	CampusClassNumber string `json:"campus_class_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	SchoolNumber      string `json:"school_number"`
	StreamNumber      string `json:"stream_number"`
	CheckoutTime      string `json:"checkout_time"`
	CheckoutReason    string `json:"checkout_reason"`
}
