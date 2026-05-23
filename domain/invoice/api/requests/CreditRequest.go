package requests

import "errors"

type CreateCreditRequest struct {
	UserNumber      string  `json:"user_number"`
	Amount          float64 `json:"amount"`
	Description     *string `json:"description"`
	Date            string  `json:"date"`
	ClientReference *string `json:"client_reference"`
}

func (r *CreateCreditRequest) Validate() error {
	if r.UserNumber == "" {
		return errors.New("choose the user receiving credit")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if r.Date == "" {
		return errors.New("choose the date")
	}
	return nil
}
