package requests

import (
	"errors"
)

type UnitOfMeasureRequest struct {
	Name string `json:"name"`
	Code uint64 `json:"code"`
}

func (u *UnitOfMeasureRequest) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Code == 0 {
		return errors.New("code is required")
	}

	return nil
}
