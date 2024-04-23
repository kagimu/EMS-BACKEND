package requests

import "errors"

type InvoicePaymentRequest struct {
	InvoiceNumber string  `json:"invoice_number"`
	Payer         string  `json:"payer"`
	Amount        float64 `json:"amount"`
	Balance       float64 `json:"balance"`
	Date          string  `json:"date"`
}

func (i *InvoicePaymentRequest) Validate() error {
	if i.Payer == "" {
		return errors.New("payer is required")
	}
	if i.Amount == 0 {
		return errors.New("amount is required")
	}
	if i.Date == "" {
		return errors.New("date is required")
	}
	if i.InvoiceNumber == "" {
		return errors.New("invoice number is required")
	}

	return nil
}
