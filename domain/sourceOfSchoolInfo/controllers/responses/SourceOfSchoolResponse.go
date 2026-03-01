package responses

type SourceOfSchoolResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	UserNumber   string `json:"user_number"`
	EntityNumber string `json:"entity_number"`
	Code         int64  `json:"code"`
}
