package services

import (
	"errors"
	"fmt"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/responses"
)

type ItemService struct {
}

func (receiver ItemService) AllServiceItems(schoolNumber string) ([]responses.ServiceResponse, error) {
	result := make([]responses.ServiceResponse, 0)

	items, err := repositories.FetchAllServiceItems(schoolNumber)
	if err != nil {
		return result, errors.New(fmt.Sprintf("Error %s fetching service items", err.Error()))
	}

	for _, item := range items {

		campusResponse := make([]responses.CampusServiceResponse, 0)

		itemCampuses, err := repositories.FetchAllCampusItems(item.Number)
		if err != nil {
			return result, errors.New(fmt.Sprintf("Error %s fetching campus items", err.Error()))
		}

		for _, itemCampus := range itemCampuses {
			campusResponse = append(campusResponse, responses.CampusServiceResponse{
				CampusNumber: itemCampus.Number,
				Name:         itemCampus.Name,
			})
		}

		result = append(result, receiver.CreateServiceResponse(item, campusResponse))
	}

	return result, nil
}

func (receiver ItemService) CreateServiceResponse(item itemModels.Item, campusServiceResponse []responses.CampusServiceResponse) responses.ServiceResponse {
	return responses.ServiceResponse{
		ItemNumber: item.Number,
		Item: responses.ItemResponse{
			Number:          item.Number,
			Name:            item.Name,
			Description:     item.Description,
			Type:            item.Type,
			Cost:            item.Cost,
			SerialCode:      item.SerialCode,
			PaymentMode:     item.PaymentMode,
			Mandatory:       item.Mandatory,
			SchoolNumber:    item.SchoolNumber,
			UserNumber:      item.UserNumber,
			RepeatCycle:     item.RepeatCycle,
			OneTimePurchase: item.OneTimePurchase,
		},
		Campuses: campusServiceResponse,
	}
}

func NewItemService() ItemService {
	return ItemService{}
}
