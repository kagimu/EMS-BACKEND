package responses

import "time"

// AdmissionRequirementResponse is the API response for a single admission requirement.
type AdmissionRequirementResponse struct {
	Number              string    `json:"number"`
	ItemNumber          string    `json:"item_number"`
	ItemName            string    `json:"item_name"`
	ItemCost            float64   `json:"item_cost"`
	CampusClassNumber   string    `json:"campus_class_number"`
	CampusNumber        string    `json:"campus_number"`
	StudyPeriodNumber   string    `json:"study_period_number"`
	AcademicYearNumber  string    `json:"academic_year_number"`
	SchoolNumber        string    `json:"school_number"`
	SchoolName          string    `json:"school_name"`
	CampusName          string    `json:"campus_name"`
	CurriculumName      string    `json:"curriculum_name"`
	ClassName           string    `json:"class_name"`
	StudyPeriodName     string    `json:"study_period_name"`
	UserNumber          string    `json:"user_number"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
