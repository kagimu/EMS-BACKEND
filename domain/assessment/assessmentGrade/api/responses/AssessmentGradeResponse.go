package responses

type AssessmentGradeResponse struct {
	Number           string  `json:"number"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	MinimumPoints    float64 `json:"minimum_points"`
	MaximumPoints    float64 `json:"maximum_points"`
	Abbreviation     string  `json:"abbreviation"`
	IllustrationPath string  `json:"illustration_path"`
	Remarks          string  `json:"remarks"`
	CampusClassName  string  `json:"campus_class_name"`
}
