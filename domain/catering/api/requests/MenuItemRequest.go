package requests

import (
	"errors"
)

type MenuItemRequest struct {
	BasicMenuItemRequest
}

type UpdateMenuItemRequest struct {
	BasicMenuItemRequest
	MenuItemNumber string `json:"menu_item_number"`
}

type BasicMenuItemRequest struct {
	Name                        string  `json:"name"`
	Code                        int64   `json:"code"`
	UnitOfMeasureNumber         string  `json:"unit_of_measure_number"`
	ConversionToPurchaseUnit    float64 `json:"conversion_to_purchase_unit"`
	PurchaseUnitOfMeasureNumber string  `json:"purchase_unit_of_measure_number"`
	SchoolNumber                string  `json:"school_number"`
	BasePrice                   float64 `json:"base_price"`
	Status                      string  `json:"status"`
}

func (r *BasicMenuItemRequest) Validate() error {

	if r.SchoolNumber == "" {
		return errors.New("choose the school")
	}

	if r.Name == "" {
		return errors.New("provide item name")
	}

	if r.Code <= 0 {
		return errors.New("sorting code is required")
	}

	if r.BasePrice <= 0 {
		return errors.New("provide base item unit price")
	}

	if r.ConversionToPurchaseUnit <= 0 {
		return errors.New("conversion to purchase unit can not be less or equal to 0")
	}

	if r.PurchaseUnitOfMeasureNumber == "" && r.ConversionToPurchaseUnit > 0 {
		return errors.New("set purchase unit of measure")
	}

	if r.PurchaseUnitOfMeasureNumber == r.UnitOfMeasureNumber {
		r.ConversionToPurchaseUnit = 1
	}

	return nil
}

func (r *MenuItemRequest) Validate() error {
	return r.BasicMenuItemRequest.Validate()
}

func (r *UpdateMenuItemRequest) Validate() error {

	if r.MenuItemNumber == "" {
		return errors.New("set the menu item number to update")
	}

	err := r.BasicMenuItemRequest.Validate()
	if err != nil {
		return err
	}

	return nil
}
