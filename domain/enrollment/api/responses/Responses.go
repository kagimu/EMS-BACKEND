package responses

type EnrolledStudent struct {
	Number            string  `json:"number"`
	StudentSurname    string  `json:"student_surname"`
	StudentSecondName string  `json:"student_second_name"`
	StudentOtherName  string  `json:"student_other_name"`
	Score             float64 `json:"score"`
	Comment           string  `json:"comment"`
	Class             string  `json:"class"`
	Status            string  `json:"status"`
}

type Enrollment struct {
	Number            string `json:"number"`
	StudentNumber     string `json:"student_number"`
	Adn               string `json:"adn"`
	StudentSurname    string `json:"student_surname"`
	StudentSecondName string `json:"student_second_name"`
	StudentOtherName  string `json:"student_other_name"`
	Class             string `json:"class"`
	Stream            string `json:"stream"`
	School            string `json:"school"`
	Campus            string `json:"campus"`
	StudyPeriod       string `json:"study_period"`
	Status            string `json:"status"`
	EnrollmentDate    string `json:"enrollment_date"`
	ContinuingStudent bool   `json:"continuing_student"`
	AcademicYear      string `json:"academic_year"`
}
