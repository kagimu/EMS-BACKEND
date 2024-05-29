package requests

import "errors"

type ItemAssembly struct {
	PrimaryItemNumber     string                  `json:"primary_item_number"`
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_items"`
	Name                  string                  `json:"name"`
	Code                  int64                   `json:"code"`
	Type                  string                  `json:"type"`
	UnitOfMeasureNumber   string                  `json:"unit_of_measure_number"`
	CampusNumber          string                  `json:"campus_number"`
}

type SecondaryItemQuantity struct {
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
}

type UpdateItemAssembly struct {
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_items"`
	PrimaryItemNumber     string                  `json:"primary_item_number"`
}

func (i *ItemAssembly) Validate() error {
	if len(i.SecondaryItemQuantity) == 0 {
		return errors.New("secondary item quantity is required")
	}

	return nil
}
