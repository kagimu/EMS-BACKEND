package responses

type AssessmentGradeResponse struct {
	AssessmentGradeNumber string                        `json:"assessment_grade_number"`
	ScaleNumber           string                        `json:"scale_number"`
	AssessmentAreaCount   int                           `json:"assessment_area_count"`
	AssessmentGradeItems  []AssessmentGradeItemResponse `json:"assessment_grade_items"`
}

type AssessmentGradeItemResponse struct {
	Code          int64   `json:"code"`
	Description   string  `json:"description"`
	MinimumPoints float64 `json:"minimum_points"`
	MaximumPoints float64 `json:"maximum_points"`
	Name          string  `json:"name"`
	Abbreviation  string  `json:"abbreviation"`
	Remarks       string  `json:"remarks"`
}
