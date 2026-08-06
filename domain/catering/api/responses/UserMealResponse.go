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

// LearnerMealHistoryResponse is used for the GET catering/meals/learners/history endpoint.
// It returns exactly the data needed for the Meal Delivery History table.
type LearnerMealHistoryResponse struct {
	Name               string `json:"name"`
	FeedingGroupName   string `json:"feeding_group_name"`
	MealName           string `json:"meal_name"`
	Date               string `json:"date"`
	Time               string `json:"time"`
	Served             bool   `json:"served"`
}

// StaffMealResponse is used by GET /catering/meals/staff.
// It mirrors UserMealResponse but targets feeding_group_staffs and includes
// staff-specific identity fields (STN, job position, contact info).
// Required filters: school_number, meal_course_number.
// Optional filter:  feeding_group_number.
type StaffMealResponse struct {
	// Meal / session fields (mirrors UserMealResponse)
	UserNumber           string               `json:"user_number"`
	StaffNumber          string               `json:"staff_number"`
	Name                 string               `json:"name"`
	FeedingGroupNumber   string               `json:"feeding_group_number"`
	FeedingGroupName     string               `json:"feeding_group_name"`
	Meals                []SingleUserMealItem `json:"meals"`
	FeedingSessionNumber string               `json:"feeding_session_number"`
	Served               bool                 `json:"served"`

	// Staff identity
	Gender      string `json:"gender"`
	STN         string `json:"stn"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	JobPosition string `json:"job_position"`

	// Campus / school
	SchoolNumber string `json:"school_number"`
	CampusNumber string `json:"campus_number"`
	CampusName   string `json:"campus_name"`
}

// StaffMealHistoryResponse is used for GET /catering/meals/staff/history.
type StaffMealHistoryResponse struct {
	STN              string `json:"stn"`
	Name             string `json:"name"`
	Position         string `json:"position"`
	FeedingGroupName string `json:"feeding_group_name"`
	MealName         string `json:"meal_name"`
	Date             string `json:"date"`
	Time             string `json:"time"`
	Served           bool   `json:"served"`
}
