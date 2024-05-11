package responses

import "time"

type InvoicePaymentResponse struct {
	Amount           float64   `json:"amount"`
	DateTime         time.Time `json:"date_time"`
	RecordedBy       string    `json:"recorded_by"`
	InvoiceRefNumber string    `json:"invoice_ref_number"`
}
