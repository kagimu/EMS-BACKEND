package responses

type ClassReportTemplateResponse struct {
	TemplateName string          `json:"template_name"`
	Classes      []ClassResponse `json:"classes"`
}

type ClassResponse struct {
	Number    string `json:"number"`
	ClassName string `json:"class_name"`
	AddedBy   string `json:"added_by"`
}
