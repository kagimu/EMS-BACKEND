package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type MenuItemRequest struct {
	Name                string  `json:"name"`
	Code                int64   `json:"code"`
	Type                string  `json:"type"`
	UnitOfMeasureNumber string  `json:"unit_of_measure_number"`
	SchoolNumber        string  `json:"school_number"`
	BasePrice           float64 `json:"base_price"`
}

func (r *MenuItemRequest) Validate() error {

	if r.SchoolNumber == "" {
		return errors.New("choose the school")
	}

	if r.Name == "" {
		return errors.New("provide item name")
	}

	if r.Code <= 0 {
		return errors.New("sorting code is required")
	}

	if r.Type == "" {
		return errors.New("type of the item not selected")
	}

	if r.Type != constants.TypeAssembled && r.Type != constants.TypeItem {
		return errors.New("the item type must be ASSEMBLED, or ITEM")
	}

	if r.BasePrice <= 0 && r.Type != constants.TypeAssembled {
		return errors.New("provide base item unit price")
	}

	return nil
}
