package responses

type InvoiceResponse struct {
	MinInvoiceResponse
	Items    []InvoiceItemResponse    `json:"items"`
	Payments []InvoicePaymentResponse `json:"payments"`
}

type MinInvoiceResponse struct {
	Number      string  `json:"number"`
	RefNumber   string  `json:"ref_number"`
	Date        string  `json:"date"`
	DueDate     string  `json:"due_date"`
	TotalAmount float64 `json:"total_amount"`
	PaidAmount  float64 `json:"paid_amount"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
}
