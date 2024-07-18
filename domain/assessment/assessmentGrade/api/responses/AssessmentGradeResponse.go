package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"

type AssessmentGradeResponse struct {
	AssessmentGradeNumber string                        `json:"assessment_grade_number"`
	AssessmentAreaCount   int                           `json:"assessment_area_count"`
	AssessmentGradeItems  []AssessmentGradeItemResponse `json:"assessment_grade_items"`
	Classes               []responses.ClassResponse     `json:"classes"`
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
