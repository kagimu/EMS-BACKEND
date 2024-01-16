package requests

type ClassInfo struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}

type ExcelSchemeOfWorkInfo struct {
	SubTopic      string
	SortingCode   string
	EstimatedTime string
}
