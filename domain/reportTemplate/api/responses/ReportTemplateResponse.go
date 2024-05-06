package responses

type ReportTemplateResponse struct {
	Name             string `json:"name"`
	SchoolName       string `json:"schoolName"`
	SchoolNumber     string `json:"schoolNumber"`
	ReportTemplateId string `json:"reportTemplateId"`
	Status           string `json:"status"`
	AddedBy          string `json:"added_by"`
}
