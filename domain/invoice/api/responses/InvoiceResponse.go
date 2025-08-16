package responses

type InvoiceResponse struct {
	MinInvoiceResponse
	Items    []InvoiceItemResponse    `json:"items"`
	Payments []InvoicePaymentResponse `json:"payments"`
}

type MinInvoiceResponse struct {
	Number          string  `json:"number"`
	RefNumber       string  `json:"ref_number"`
	Date            string  `json:"date"`
	DueDate         string  `json:"due_date"`
	IssueDate       string  `json:"issue_date"`
	TotalAmount     float64 `json:"total_amount"`
	PaidAmount      float64 `json:"paid_amount"`
	Balance         float64 `json:"balance"`
	Status          string  `json:"status"`
	Adn             string  `json:"adn"`
	Lin             string  `json:"lin"`
	Names           string  `json:"names"`
	CampusName      string  `json:"campus_name"`
	StudyPeriodName string  `json:"study_period_name"`
	StudentNumber   string  `json:"student_number"`
}
