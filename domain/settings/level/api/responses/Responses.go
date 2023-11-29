package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"

type (
	LevelList struct {
		Number     string `json:"number"`
		Name       string `json:"name"`
		Code       int64  `json:"code"`
		UserNumber string `json:"user_number"`
	}

	ClassLevelList struct {
		Number        string                    `json:"number"`
		Name          string                    `json:"name"`
		Code          int64                     `json:"code"`
		UserNumber    string                    `json:"user_number"`
		ClassDataList []responses.ClassDataList `json:"class"`
	}

	CurriculumLevelList struct {
		Number           string `json:"number"`
		Name             string `json:"name"`
		Code             int64  `json:"code"`
		UserNumber       string `json:"user_number"`
		CurriculumName   string `json:"curriculum_name"`
		CurriculumNumber string `json:"curriculum_number"`
	}
)
