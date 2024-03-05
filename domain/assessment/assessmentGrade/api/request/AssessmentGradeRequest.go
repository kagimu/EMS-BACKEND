package requests

type AssessmentGradeRequest struct {
	Grade         string `json:"grade"`
	Description   string `json:"description"`
	MinimumPoints int64  `json:"minimum_points"`
	MaximumPoints int64  `json:"maximum_points"`
	Name          string `json:"name"`
	Abbreviation  string `json:"abbreviation"`
	Condition     string `json:"condition"`
	Remarks       string `json:"remarks"`
	LevelNumber   string `json:"level_number"`
}
