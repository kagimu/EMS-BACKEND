package responses

type AssessmentAreaResponse struct {
	Code                int64  `json:"code"`
	AssessmentAreaName  string `json:"assessment_area_name"`
	GradingStandardName string `json:"grading_standard_name"`
}
