package requests

type ClassInfo struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
}

type EditClassInfo struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
	ClassNumber  string `json:"class_number"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}

type ClassDelRequest struct {
	Number string `json:"number"`
}
