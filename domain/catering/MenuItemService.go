package catering

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
)

func GetMenuItemParts(menuItemNumber string) ([]responses.Part, error) {

	var parts []responses.Part

	itemAssembles, err := repositories.GetAssembledItemParts(menuItemNumber)
	if err != nil {
		return nil, err
	}

	for _, itemAssemble := range itemAssembles {

		secondaryItem := itemAssemble.SecondaryItem

		part := responses.Part{
			Name:     secondaryItem.Name,
			Quantity: itemAssemble.Quantity,
			Code:     secondaryItem.Code,
			Number:   secondaryItem.Number,
		}

		parts = append(parts, part)

	}

	return parts, nil
}
