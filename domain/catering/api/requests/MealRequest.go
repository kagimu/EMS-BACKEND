package requests

import "errors"

type MealRequest struct {
	UserNumber           string `json:"user_number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	StockingNumber       string `json:"stocking_number"`
	MenuItemNumber       string `json:"menu_item_number"`
	StaffNumber          string `json:"staff_number"`
}

func (m *MealRequest) Validate() error {
	if m.FeedingSessionNumber == "" {
		return errors.New("feeding session number is required")
	}
	if m.StockingNumber == "" {
		return errors.New("stocking number is required")
	}
	if m.MenuItemNumber == "" {
		return errors.New("menu item number is required")
	}
	if m.StaffNumber == "" {
		return errors.New("staff number is required")
	}
	if m.UserNumber == "" {
		return errors.New("user number is required")
	}

	return nil
}
