package requests

type LevelDelRequest struct {
	Number string `json:"number"`
}

type LevelEditInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	LevelNumber string `json:"level_number"`
}
