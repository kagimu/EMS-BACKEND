package responses

type GradingStandardAssessmentAreaResponse struct {
	Code                int64  `json:"code"`
	AssessmentAreaName  string `json:"assessment_area_name"`
	GradingStandardName string `json:"grading_standard_name"`
}

type AssessmentAreaResponse struct {
	Number                string `json:"number"`
	Name                  string `json:"name"`
	SubjectNumber         string `json:"subject_number"`
	UserNumber            string `json:"user_number"`
	CurriculumClassNumber string `json:"curriculum_class_number"`
	CurriculumNumber      string `json:"curriculum_number"`
	CurriculumName        string `json:"curriculum_name"`
}
