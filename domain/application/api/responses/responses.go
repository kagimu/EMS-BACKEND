package responses

type ApplicationResponse struct {
	Number                    string                   `json:"number"`
	Status                    string                   `json:"status"`
	ApplicationFeePayment     bool                     `json:"application_fee_payment"`
	ApplicationFeeInvalidDate string                   `json:"application_fee_invalid_date"`
	UserNumber                string                   `json:"user_number"`
	InterviewMode             string                   `json:"interview_mode"`
	Student                   StudentResponse          `json:"student"`
	Parent                    ParentResponse           `json:"parent"`
	PreviousSchools           []PreviousSchoolResponse `json:"previous_schools"`
}

type StudentResponse struct {
	Number                    string `json:"number"`
	APN                       string `json:"apn"`
	ADN                       string `json:"adn"`
	APNStatus                 string `json:"apn_status"`
	ADNStatus                 string `json:"adn_status"`
	Surname                   string `json:"surname" `
	SecondName                string `json:"second_name"`
	OtherName                 string `json:"other_name"`
	Photo                     string `json:"photo"`
	AcademicReport            string `json:"academic_report"`
	Gender                    string `json:"gender"`
	ReasonForChoosingSchool   string `json:"reason_for_choosing_school"`
	ReasonForAdmissionDenial  string `json:"reason_for_admission_denial"`
	AttentionNeeded           string `json:"attention_needed"`
	DateOfBirth               string `json:"date_of_birth"`
	PreferredTypeOfEnrollment string `json:"preferred_type_of_enrollment"`
	PreferredCampusNumber     string `json:"preferred_campus_number"`
	PreferredCampus           string `json:"preferred_campus"`
	CampusClassNumber         string `json:"campus_class_number"`
	CampusClass               string `json:"campus_class"`
	SourceOfSchoolInfo        string `json:"source_of_school_info"`
	Status                    string `json:"status"`
	SchoolNumber              string `json:"school_number"`
	CampusCurriculumNumber    string `json:"campus_curriculum_number"`
	LevelNumber               string `json:"level_number"`
	CurriculumLevelNumber     string `json:"curriculum_level_number"`
}
type PreviousSchoolResponse struct {
	Name    string                        `json:"name"`
	Classes []PreviousSchoolClassResponse `json:"classes"`
}

type PreviousSchoolClassResponse struct {
	PreviousSchoolNumber string `json:"previous_school_number"`
	Year                 string `json:"year"`
	Class                string `json:"class" `
}

type ParentResponse struct {
	Surname           string `json:"surname"`
	SecondName        string `json:"second_name"`
	OtherNames        string `json:"other_names"`
	TelephoneNumber   string `json:"telephone_number"`
	AlternativeNumber string `json:"alternative_number"`
	Email             string `json:"email"`
	AlternativeEmail  string `json:"alternative_email"`
}
