package requests

type ClassInfo struct {
	Name             string `json:"name"`
	Code             string `json:"code"`
	Abbreviation     string `json:"abbreviation"`
	Level            string `json:"level"`
	CurriculumNumber string `json:"curriculum_number"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}
