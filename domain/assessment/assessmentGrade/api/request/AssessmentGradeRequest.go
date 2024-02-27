package requests

type AssessmentGradeRequest struct {
	Grade          string `json:"grade"`
	Description    string `json:"description"`
	AggregateRange string `json:"aggregate_range"`
	Abbreviation   string `json:"abbreviation"`
	Condition      string `json:"condition"`
	Remarks        string `json:"remarks"`
	LevelNumber    string `json:"level_number"`
}
