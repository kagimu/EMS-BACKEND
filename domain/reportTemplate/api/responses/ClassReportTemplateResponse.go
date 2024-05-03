package responses

type ClassReportTemplateResponse struct {
	TemplateName string
	Classes      []ClassResponse
}

type ClassResponse struct {
	Number    string
	ClassName string
	AddedBy   string
}
