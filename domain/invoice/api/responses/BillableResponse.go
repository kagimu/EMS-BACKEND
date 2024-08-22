package responses

type BillableResponse struct {
	BillableNumber    string  `json:"billable_number"`
	Status            string  `json:"status"`
	ClassName         string  `json:"class_name"`
	ClassNumber       string  `json:"class_number"`
	StudyPeriodNumber string  `json:"study_period_number"`
	StudyPeriodName   string  `json:"study_period_name"`
	CampusNumber      string  `json:"campus_number"`
	CampusName        string  `json:"campus_name"`
	Type              string  `json:"type"`
	TotalAmount       float64 `json:"total_amount"`
}

type DetailedBillableResponse struct {
	BillableResponse
	BillableItems []InvoiceItemResponse `json:"billable_items"`
}
