package requests

import "errors"

type ItemAssembly struct {
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_items"`
	MenuItemRequest
}

type SecondaryItemQuantity struct {
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
}

type UpdateItemAssembly struct {
	UpdateMenuItemRequest
	SecondaryItemQuantity []SecondaryItemQuantity `json:"secondary_items"`
}

func (r *ItemAssembly) Validate() error {

	if r.SchoolNumber == "" {
		return errors.New("choose the school")
	}

	if r.Name == "" {
		return errors.New("provide item name")
	}

	if r.Code <= 0 {
		return errors.New("sorting code is required")
	}

	if len(r.SecondaryItemQuantity) < 1 {
		return errors.New("provide at least one secondary item")
	}

	for _, secItem := range r.SecondaryItemQuantity {

		if secItem.SecondaryItemNumber == "" {
			return errors.New("one of the secondary items has no number set")
		}

		if secItem.Quantity <= 0 {
			return errors.New("one of the secondary items has no quantity set")
		}

	}

	return r.MenuItemRequest.Validate()
}

func (r *UpdateItemAssembly) Validate() error {

	if r.MenuItemNumber == "" {
		return errors.New("choose the item to edit")
	}

	if r.Name == "" {
		return errors.New("provide item name")
	}

	if r.Code <= 0 {
		return errors.New("sorting code is required")
	}

	if len(r.SecondaryItemQuantity) < 2 {
		return errors.New("provide at least two secondary items")
	}

	for _, secItem := range r.SecondaryItemQuantity {

		if secItem.SecondaryItemNumber == "" {
			return errors.New("one of the secondary items has no number set")
		}

		if secItem.Quantity <= 0 {
			return errors.New("one of the secondary items has no quantity set")
		}

	}

	return nil
}
