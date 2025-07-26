package requests

import "errors"

type InvoicePaymentRequest struct {
	InvoiceNumber string   `json:"invoice_number"`
	Amount        float64  `json:"amount"`
	Date          string   `json:"date"`
	Credits       []string `json:"credits"`
}

func (i *InvoicePaymentRequest) Validate() error {

	if len(i.Credits) == 0 {
		if i.Amount <= 0 {
			return errors.New("amount should be greater than zero")
		}
		if i.Date == "" {
			return errors.New("date when paid is required")
		}
	} else {
		for _, creditNumber := range i.Credits {
			if creditNumber == "" {
				return errors.New("You have empty credit ID ")
			}
		}
	}

	if i.InvoiceNumber == "" {
		return errors.New("choose the promissory note they are paying against")
	}

	return nil
}
