package requests

type ReportsFilterRequest struct {
	SchoolNumber string `form:"school_number"`
	CampusNumber string `form:"campus_number"`
	AcademicYear string `form:"academic_year"`
	DateFrom     string `form:"date_from"`
	DateTo       string `form:"date_to"`
	Level        string `form:"level"`
	ClassNo      string `form:"class_no"`
}
