package requests

import "errors"

type InvoicePaymentRequest struct {
	PaymentReference string `json:"payment_reference"`
	Payer            string `json:"payer"`
	Amount           int    `json:"amount"`
	Balance          int    `json:"balance"`
	Date             string `json:"date"`
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
	return nil
}
