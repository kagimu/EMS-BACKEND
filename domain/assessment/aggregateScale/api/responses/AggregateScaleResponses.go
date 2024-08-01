package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/responses"

type AggregateScaleResponse struct {
	Code         int64   `json:"code"`
	Number       string  `json:"number"`
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	Name         string  `json:"name"`
	Weight       float64 `json:"weight"`
	Illustration string  `json:"illustration"`
	Description  string  `json:"description"`
	ShortName    string  `json:"short_name"`
	ScaleName    string  `json:"scale_name"`
	GradeNumber  string  `json:"grade_number"`
	Status       string  `json:"status"`
}

type ScaleResponse struct {
	ScaleNumber string `json:"scale_number"`
	ScaleName   string `json:"scale_name"`
	Count       int    `json:"count"`
}

type DetailedScaleResponse struct {
	ScaleResponse
	AggregateScaleItems []AggregateScaleResponse          `json:"aggregate_scales"`
	AssessmentGrade     responses.AssessmentGradeResponse `json:"assessment_grade"`
}
