package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

func (receiver ItemService) AllItems(schoolNumber, itemType string) ([]responses.ItemResponse, error) {
	itemResponses := make([]responses.ItemResponse, 0)

	items, err := repositories.FetchAllItems(schoolNumber, itemType)
	if err != nil {
		return itemResponses, errors.New(fmt.Sprintf("Error %s fetching items", err.Error()))
	}

	for _, item := range items {
		itemResponses = append(itemResponses, receiver.createItemResponse(item))
	}

	return itemResponses, nil
}

func (receiver ItemService) createItemResponse(item itemModels.Item) responses.ItemResponse {
	return responses.ItemResponse{
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
		Prefix:          item.Prefix,
	}
}

func (receiver ItemService) SaveService(user *models.User, payload requests2.ServiceRequest) error {

	checkSerial, err := repositories.GetItemBySerialCode(payload.SerialCode, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if checkSerial.Number != "" {
		return errors.New("Item serial code already used ")
	}

	itemNumber := uuid.NewV4().String()
	saveProductItem := itemModels.Item{
		Number:          itemNumber,
		Cost:            payload.Cost,
		Description:     payload.Description,
		OneTimePurchase: payload.OneTimePurchase,
		UserNumber:      user.Number,
		SerialCode:      payload.SerialCode,
		Mandatory:       payload.Mandatory,
		RepeatCycle:     payload.RepeatCycle,
		PaymentMode:     payload.PaymentMode,
		Name:            payload.Name,
		SchoolNumber:    payload.SchoolNumber,
		Type:            "Service",
	}

	saveService := models2.Service{
		ItemNumber: itemNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := repositories.SaveItem(saveProductItem, tx)

		err = repositories.SaveServiceItem(saveService, tx)

		for _, campusNumber := range payload.CampusNumbers {
			campusService := models2.CampusItem{
				ItemNumber:   itemNumber,
				CampusNumber: campusNumber,
			}

			check, err := repositories.FindCampusItem(campusNumber, itemNumber)
			if err != nil {
				return err
			}

			if check.Number == "" {
				err = repositories.SaveCampusItem(campusService, tx)
			}
		}

		return err
	})

	if err != nil {
		return errors.New(fmt.Sprintf("Error %s saving service item", err.Error()))
	}

	return nil
}

func (receiver ItemService) UpdateService(serviceNumber string, user *models.User, payload requests2.ServiceRequest) error {

	product, err := repositories.GetItemByNumber(serviceNumber)
	if err != nil {
		return err
	}

	if product.Number == "" {
		return errors.New("Invalid product number ")
	}

	byName, err := repositories.GetItemByName(payload.Name, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if byName.Number != "" && byName.Number != product.Number {
		return errors.New("Item with name already exists ")
	}

	checkSerial, err := repositories.GetItemBySerialCode(payload.SerialCode, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if checkSerial.Number != "" && checkSerial.Number != product.Number {
		return errors.New("Item serial code already used ")
	}

	saveProductItem := itemModels.Item{
		Number:          product.Number,
		Cost:            payload.Cost,
		Description:     payload.Description,
		OneTimePurchase: payload.OneTimePurchase,
		UserNumber:      user.Number,
		SerialCode:      payload.SerialCode,
		Mandatory:       payload.Mandatory,
		RepeatCycle:     payload.RepeatCycle,
		PaymentMode:     payload.PaymentMode,
		Name:            payload.Name,
		SchoolNumber:    payload.SchoolNumber,
		Type:            "Service",
		CreatedAt:       product.CreatedAt,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := repositories.SaveItem(saveProductItem, tx)
		if err != nil {
			return err
		}

		err = repositories.DeleteCampusItem(product.Number, tx)
		if err != nil {
			return err
		}

		for _, campusNumber := range payload.CampusNumbers {
			campusService := models2.CampusItem{
				ItemNumber:   product.Number,
				CampusNumber: campusNumber,
			}

			check, err := repositories.FindCampusItem(campusNumber, product.Number)
			if err != nil {
				return err
			}

			if check.Number == "" {
				err = repositories.SaveCampusItem(campusService, tx)
			}
		}

		return err
	})

	if err != nil {
		return errors.New(fmt.Sprintf("Error %s updating service item", err.Error()))
	}

	return nil
}

func NewItemService() ItemService {
	return ItemService{}
}
