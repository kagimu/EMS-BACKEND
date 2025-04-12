package requests

import "errors"

type StockingRequest struct {
	CampusNumber         string                 `json:"campus_number"`
	ItemPriceAndQuantity []ItemPriceAndQuantity `json:"item_price_and_quantity"`
}

type ItemPriceAndQuantity struct {
	Quantity         float64 `json:"quantity"`
	PurchaseUnitCost float64 `json:"unit_cost"`
	MenuItemNumber   string  `json:"menu_item_number"`
}

func (s *StockingRequest) Validate() error {
	if len(s.ItemPriceAndQuantity) == 0 {
		return errors.New("items prices and quantities are required ")
	}
	if s.CampusNumber == "" {
		return errors.New("campus number is required")
	}

	for _, item := range s.ItemPriceAndQuantity {
		if item.Quantity <= 0 {
			return errors.New("quantity is required ")
		}

		if item.PurchaseUnitCost < 0 {
			return errors.New("purchase unit cost can not be negative")
		}
	}

	return nil
}
