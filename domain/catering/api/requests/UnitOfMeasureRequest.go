package requests

import (
	"errors"
)

type UnitOfMeasureRequest struct {
	Name         string `json:"name"`
	Code         uint64 `json:"code"`
	Abbreviation string `json:"abbreviation"`
	SchoolNumber string `json:"school_number"`
}

type EditUnitOfMeasureRequest struct {
	UnitOfMeasureRequest
	UnitOfMeasureNumber string `json:"unit_of_measure_number"`
	Status              string `json:"status"`
}

func (u *UnitOfMeasureRequest) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Code == 0 {
		return errors.New("code is required")
	}
	if u.Abbreviation == "" {
		return errors.New("abbreviation is required")
	}

	return nil
}

func (u *EditUnitOfMeasureRequest) Validate() error {

	if u.UnitOfMeasureNumber == "" {
		return errors.New("set the unit of measure number to edit")
	}

	if err := u.UnitOfMeasureRequest.Validate(); err != nil {
		return err
	}

	return nil
}
