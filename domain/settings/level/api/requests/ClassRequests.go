package requests

type ClassInfo struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}

type LevelDelRequest struct {
	Number string `json:"number"`
}
