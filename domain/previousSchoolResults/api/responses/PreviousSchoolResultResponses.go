package responses

type PreviousSchoolResultResponses struct {
	Number             string  `json:"number"`
	StudentName        string  `json:"student_name"`
	AssessmentAreaName string  `json:"assessment_area_name"`
	Score              float64 `json:"score"`
	Attachment         string  `json:"attachment"`
	ApplicationNumber  string  `json:"application_number"`
}
