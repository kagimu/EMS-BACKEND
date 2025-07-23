package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/api/requests"
	productModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ItemProduct struct {
}

func (receiver ItemProduct) AllProductItems(schoolNumber string) ([]responses.ServiceResponse, error) {
	result := make([]responses.ServiceResponse, 0)

	items, err := repositories.FetchAllProductItems(schoolNumber)
	if err != nil {
		return result, errors.New(fmt.Sprintf("Error %s fetching product items", err.Error()))
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

func (receiver ItemProduct) CreateServiceResponse(item itemModels.Item, campusServiceResponse []responses.CampusServiceResponse) responses.ServiceResponse {
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
			Code:            item.Code,
		},
		Campuses: campusServiceResponse,
	}
}

func (receiver ItemProduct) SaveProduct(user *models.User, payload requests.ProductRequest) error {

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
		Type:            "Product",
		Code:            payload.Code,
	}

	saveProduct := productModels.Product{
		ItemNumber: itemNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := repositories.SaveItem(saveProductItem, tx)

		err = repositories.SaveProductItem(saveProduct, tx)

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
		return errors.New(fmt.Sprintf("Error %s saving product items", err.Error()))
	}

	return nil
}

func (receiver ItemProduct) UpdateProduct(productNumber string, user *models.User, payload requests.ProductRequest) error {

	product, err := repositories.GetItemByNumber(productNumber)
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
		Type:            "Product",
		CreatedAt:       product.CreatedAt,
		Code:            payload.Code,
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
		return errors.New(fmt.Sprintf("Error %s saving product items", err.Error()))
	}

	return nil
}

func NewItemProduct() ItemProduct {
	return ItemProduct{}
}
