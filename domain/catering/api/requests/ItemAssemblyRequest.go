package requests

import "errors"

type ItemAssesmbly struct {
	PrimaryItemNumber     string                  `json:"primary_item_number"`
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_item_quantity"`
}

type SecondaryItemQuantity struct {
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
}

func (i *ItemAssesmbly) Validate() error {
	if i.PrimaryItemNumber == "" {
		return errors.New("primary item number is required")
	}
	if len(i.SecondaryItemQuantity) == 0 {
		return errors.New("secondary item quantity is required")
	}

	return nil
}
