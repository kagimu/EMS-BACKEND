package responses

type ReportTemplateResponse struct {
	Name             string `json:"name"`
	SchoolName       string `json:"school_name"`
	SchoolNumber     string `json:"school_number"`
	ReportTemplateId string `json:"report_template_id"`
	Status           string `json:"status"`
	AddedBy          string `json:"added_by"`
}
