package responses

type AggregateScaleResponse struct {
	Code                int64   `json:"code"`
	Number              string  `json:"number"`
	MinValue            float64 `json:"min_value"`
	MaxValue            float64 `json:"max_value"`
	Name                string  `json:"name"`
	Weight              float64 `json:"weight"`
	Illustration        string  `json:"illustration"`
	Description         string  `json:"description"`
	ClassName           string  `json:"class_name"`
	CampusName          string  `json:"campus_name"`
	GradingStandardName string  `json:"grading_standard_name"`
}
