package requests

import "errors"

type MenuItemRequest struct {
	Name                string `json:"name"`
	Code                uint64 `json:"code"`
	Type                string `json:"type"`
	UnitOfMeasureNumber string `json:"unit_of_measure_number"`
}

func (r *MenuItemRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Code == 0 {
		return errors.New("code is required")
	}
	if r.Type == "" {
		return errors.New("type is required")
	}
	return nil
}
