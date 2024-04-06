package requests

import "errors"

type ItemAssesmbly struct {
	PrimaryItemNumber   []string `json:"primary_item_number"`
	SecondaryItemNumber []string `json:"secondary_item_number"`
	Quantity            float64  `json:"quantity"`
}

func (i *ItemAssesmbly) Validate() error {
	if i.Quantity == 0 {
		return errors.New("quantity is required")
	}
	if len(i.PrimaryItemNumber) == 0 {
		return errors.New("primary item number is required")
	}
	if len(i.SecondaryItemNumber) == 0 {
		return errors.New("secondary item number is required")
	}

	return nil
}
