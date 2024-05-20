package requests

import "errors"

type FeedingSessionItemRequest struct {
	FeedingSessionNumber string             `json:"feeding_session_number"`
	MenuItemQuantity     []MenuItemQuantity `json:"menu_item_quantity"`
}

type MenuItemQuantity struct {
	MenuItemNumber string  `json:"menu_item_number"`
	Quantity       float64 `json:"quantity"`
}

func (f *FeedingSessionItemRequest) Validate() error {
	if len(f.MenuItemQuantity) == 0 {
		return errors.New("menu item quantity is required")
	}

	return nil
}
