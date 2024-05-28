package requests

import "errors"

type InvoicePaymentRequest struct {
	InvoiceNumber string  `json:"invoice_number"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
}

func (i *InvoicePaymentRequest) Validate() error {

	if i.Amount <= 0 {
		return errors.New("amount should be greater than zero")
	}
	if i.Date == "" {
		return errors.New("date when paid is required")
	}
	if i.InvoiceNumber == "" {
		return errors.New("choose the invoice they are paying against")
	}

	return nil
}
