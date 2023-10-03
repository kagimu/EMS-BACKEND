package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

type ClassLevelInfo struct {
	Number       string          `json:"number"`
	ClassNumber  string          `json:"class_number"`
	Name         string          `json:"name"`
	Abbreviation string          `json:"abbreviation"`
	Streams      []models.Stream `json:"streams"`
}
