package responses

type ReportTemplateResponse struct {
	Name             string `json:"name"`
	ReportTemplateId string `json:"reportTemplateId"`
	Status           string `json:"status"`
	AddedBy          string `json:"added_by"`
}
