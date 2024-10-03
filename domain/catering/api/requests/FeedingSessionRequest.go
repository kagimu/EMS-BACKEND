package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type UpdateFeedingSessionRequest struct {
	FeedingSessionRequest
	FeedingSessionNumber string `json:"feeding_session_number"`
}

type FeedingSessionRequest struct {
	Name         string             `json:"name"`
	Status       string             `json:"status"`
	SchoolNumber string             `json:"school_number"`
	MenuItems    []MenuItemQuantity `json:"menu_items"`
}

type MenuItemQuantity struct {
	MenuItemNumber string  `json:"menu_item_number"`
	Quantity       float64 `json:"quantity"`
}

func (r *FeedingSessionRequest) Validate() error {
	if r.Name == "" {
		return errors.New("provide name of session")
	}
	if r.Status == "" {
		return errors.New("choose the status of session")
	}
	if r.Status != constants.StatusActive && r.Status != constants.StatusInactive {
		return errors.New("the status of this session must be one of [ACTIVE, INACTIVE]")
	}
	if r.SchoolNumber == "" {
		return errors.New("choose the school and set the number")
	}

	if len(r.MenuItems) == 0 {
		return errors.New("choose at least one menu item")
	}

	for _, menuItem := range r.MenuItems {

		if menuItem.MenuItemNumber == "" {
			return errors.New("choose the menu item")
		}

		if menuItem.Quantity <= 0 {
			return errors.New("menu item quantity is required")
		}
	}

	return nil
}

func (r *UpdateFeedingSessionRequest) Validate() error {

	err := r.FeedingSessionRequest.Validate()
	if err != nil {
		return err
	}

	if r.FeedingSessionNumber == "" {
		return errors.New("choose feeding session to update")
	}

	return nil
}
