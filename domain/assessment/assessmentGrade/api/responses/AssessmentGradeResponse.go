package responses

type AssessmentGradeResponse struct {
	AssessmentGradeCount int                           `json:"assessment_grade_count"`
	AssessmentAreaCount  int                           `json:"assessment_area_count"`
	AssessmentGrades     []AssessmentGradeItemResponse `json:"assessment_grades"`
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
