package responses

type InvoicePaymentResponse struct {
	Amount           float64 `json:"amount"`
	DateTime         string  `json:"date_time"`
	RecordedBy       string  `json:"recorded_by"`
	InvoiceRefNumber string  `json:"invoice_ref_number"`
}
