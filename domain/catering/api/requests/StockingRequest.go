package requests

import (
	"errors"
	"time"
)

type StockingRequest struct {
	CampusNumber         string                 `json:"campus_number"`
	DateOfPurchase       string                 `json:"date_of_purchase"`
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

	if s.DateOfPurchase == "" {
		return errors.New("date of purchase is required")
	}

	_, err := time.Parse(time.DateOnly, s.DateOfPurchase)
	if err != nil {
		return errors.New("invalid date of purchase format, expected YYYY-MM-DD")
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
