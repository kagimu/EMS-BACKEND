package responses

type InvoicePaymentResponse struct {
	Amount           float64 `json:"amount"`
	DateTime         string  `json:"date_time"`
	PaymentDate      string  `json:"payment_date"`
	RecordedBy       string  `json:"recorded_by"`
	InvoiceRefNumber string  `json:"invoice_ref_number"`
	Status           string  `json:"status"`
	PaymentNumber    string  `json:"payment_number"`
}
