package requests

import "errors"

type ItemAssembly struct {
	PrimaryItemNumber     string                  `json:"primary_item_number"`
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_item_quantity"`
	Name                  string                  `json:"name"`
	Code                  uint64                  `json:"code"`
	Type                  string                  `json:"type"`
	UnitOfMeasureNumber   string                  `json:"unit_of_measure_number"`
	CampusNumber          string                  `json:"campus_number"`
}

type SecondaryItemQuantity struct {
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
}

func (i *ItemAssembly) Validate() error {
	if i.PrimaryItemNumber == "" {
		return errors.New("primary item number is required")
	}
	if len(i.SecondaryItemQuantity) == 0 {
		return errors.New("secondary item quantity is required")
	}

	return nil
}
