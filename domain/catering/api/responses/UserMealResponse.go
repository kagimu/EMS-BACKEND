package responses

type UserMealResponse struct {
	UserNumber           string               `json:"user_number"`
	Name                 string               `json:"name"`
	FeedingGroupName     string               `json:"feeding_group_name"`
	FeedingGroupNumber   string               `json:"feeding_group_number"`
	Meals                []SingleUserMealItem `json:"meals"`
	FeedingSessionNumber string               `json:"feeding_session_number"`
	Served               bool                 `json:"served"`
}

// LearnerMealResponse extends UserMealResponse with full enrollment details.
// It is used by the GET catering/meals/learners endpoint and only includes
// users that are currently enrolled (active enrollment).
type LearnerMealResponse struct {
	// Meal / session fields (mirrors UserMealResponse)
	UserNumber           string               `json:"user_number"`
	Name                 string               `json:"name"`
	FeedingGroupName     string               `json:"feeding_group_name"`
	FeedingGroupNumber   string               `json:"feeding_group_number"`
	Meals                []SingleUserMealItem `json:"meals"`
	FeedingSessionNumber string               `json:"feeding_session_number"`
	Served               bool                 `json:"served"`

	// Learner identity
	StudentSurname    string `json:"student_surname"`
	StudentSecondName string `json:"student_second_name"`
	StudentOtherName  string `json:"student_other_name"`
	Gender            string `json:"gender"`
	ADN               string `json:"adn"`
	Photo             string `json:"photo"`

	// Enrollment details
	EnrollmentNumber     string `json:"enrollment_number"`
	EnrollmentStatus     string `json:"enrollment_status"`
	EnrollmentDate       string `json:"enrollment_date"`
	ContinuingStudent    bool   `json:"continuing_student"`
	StudyPeriodStartDate string `json:"study_period_start_date"`
	StudyPeriodEndDate   string `json:"study_period_end_date"`
	CampusClassNumber    string `json:"campus_class_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	StreamNumber         string `json:"stream_number"`
	Class                string `json:"class"`
	Level                string `json:"level"`
	AcademicYear         string `json:"academic_year"`
	StudyPeriod          string `json:"study_period"`

	// Campus / school
	CampusNumber string `json:"campus_number"`
	CampusName   string `json:"campus_name"`
	SchoolNumber string `json:"school_number"`
}
