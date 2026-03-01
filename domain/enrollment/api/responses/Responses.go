package responses

type EnrolledStudent struct {
	APN               string  `json:"apn"`
	ADN               string  `json:"adn"`
	Number            string  `json:"number"`
	StudentSurname    string  `json:"student_surname"`
	StudentSecondName string  `json:"student_second_name"`
	StudentOtherName  string  `json:"student_other_name"`
	Score             float64 `json:"score"`
	ActualScore       float64 `json:"actual_score"`
	Comment           string  `json:"comment"`
	Class             string  `json:"class"`
	Status            string  `json:"status"`
}

type NarrativeEnrolledStudent struct {
	APN               string `json:"apn"`
	ADN               string `json:"adn"`
	Number            string `json:"number"`
	StudentSurname    string `json:"student_surname"`
	StudentSecondName string `json:"student_second_name"`
	StudentOtherName  string `json:"student_other_name"`
	Narrative         string `json:"narrative"`
	SuggestedSolution string `json:"suggested_solution"`
	Class             string `json:"class"`
}

type Enrollment struct {
	Number                 string `json:"number"`
	StudentNumber          string `json:"student_number"`
	Adn                    string `json:"adn"`
	StudentSurname         string `json:"student_surname"`
	StudentSecondName      string `json:"student_second_name"`
	StudentOtherName       string `json:"student_other_name"`
	Class                  string `json:"class"`
	Stream                 string `json:"stream"`
	School                 string `json:"school"`
	Campus                 string `json:"campus"`
	StudyPeriod            string `json:"study_period"`
	Status                 string `json:"status"`
	EnrollmentDate         string `json:"enrollment_date"`
	ContinuingStudent      bool   `json:"continuing_student"`
	AcademicYear           string `json:"academic_year"`
	CampusCurriculumNumber string `json:"campus_curriculum_number"`
	LevelNumber            string `json:"level_number"`
	Photo                  string `json:"photo"`
	UserNumber             string `json:"user_number"`
	CampusClassNumber      string `json:"campus_class_number"`
	CampusNumber           string `json:"campus_number"`
	StudyPeriodNumber      string `json:"study_period_number"`
	ClassNumber            string `json:"class_number"`
	CurriculumNumber       string `json:"curriculum_number"`
	SchoolNumber           string `json:"school_number"`
	StreamNumber           string `json:"stream_number"`
}

// EnrollmentItemWithName is an enrollment item with the related item name for API responses.
type EnrollmentItemWithName struct {
	ID               uint    `json:"id"`
	EnrollmentNumber string  `json:"enrollment_number"`
	ItemNumber       string  `json:"item_number"`
	ItemName         string  `json:"item_name"`
	Quantity         float64 `json:"quantity"`
	Brought          bool    `json:"brought"`
}

// EnrollmentDetails is enrollment with its items (including item names) for the details endpoint.
type EnrollmentDetails struct {
	Number            string                   `json:"number"`
	StudentNumber     string                   `json:"student_number"`
	Adn               string                   `json:"adn"`
	StudentSurname    string                   `json:"student_surname"`
	StudentSecondName string                   `json:"student_second_name"`
	StudentOtherName  string                   `json:"student_other_name"`
	Class             string                   `json:"class"`
	Stream            string                   `json:"stream"`
	School            string                   `json:"school"`
	Campus            string                   `json:"campus"`
	StudyPeriod       string                   `json:"study_period"`
	Status            string                   `json:"status"`
	EnrollmentDate    string                   `json:"enrollment_date"`
	ContinuingStudent bool                     `json:"continuing_student"`
	CampusClassNumber string                   `json:"campus_class_number"`
	CampusNumber      string                   `json:"campus_number"`
	StudyPeriodNumber string                   `json:"study_period_number"`
	SchoolNumber      string                   `json:"school_number"`
	StreamNumber      string                   `json:"stream_number"`
	EnrollmentItems   []EnrollmentItemWithName `json:"enrollment_items"`
}
