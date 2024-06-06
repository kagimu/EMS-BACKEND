package requests

import "errors"

type UserEntityRequest struct {
	UserNumber   string `json:"user_number"`
	EntityNumber string `json:"entity_number"`
}

func (u *UserEntityRequest) Validate() error {
	if u.UserNumber == "" {
		return errors.New("user number is required")
	}

	if u.EntityNumber == "" {
		return errors.New("entity number is required")
	}

	return nil
}
