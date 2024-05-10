package catering

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
)

func GetMenuItemParts(menuItemNumber string) ([]responses.Part, error) {

	var parts []responses.Part

	itemAssembles, err := repositories.GetAssembledItemParts(menuItemNumber)
	if err != nil {
		return nil, err
	}

	for _, itemAssemble := range itemAssembles {

		secondaryItem := itemAssemble.SecondaryItem

		unitOfMeasure, err := repositories.GetUnitOfMeasureByNumber(secondaryItem.UnitOfMeasureNumber)

		if err != nil {
			return nil, err
		}

		if unitOfMeasure.Number == "" {
			return nil, errors.New("unit of measure number is empty")
		}

		part := responses.Part{
			Name:                      secondaryItem.Name,
			Quantity:                  itemAssemble.Quantity,
			Code:                      secondaryItem.Code,
			Number:                    secondaryItem.Number,
			UnitOfMeasureAbbreviation: unitOfMeasure.Abbreviation,
		}

		parts = append(parts, part)

	}

	return parts, nil
}

func ProcessNewStock(stockingRequest requests.StockingRequest, currentUserNumber string) error {

	db := config.DB

	if err := db.Transaction(func(transaction *gorm.DB) error {
		return doRunOnTransaction(transaction, stockingRequest, currentUserNumber)
	}); err != nil {
		return err
	}

	return nil
}

func doRunOnTransaction(transaction *gorm.DB, stockingRequest requests.StockingRequest, currentUserNumber string) error {

	for _, stock := range stockingRequest.ItemPriceAndQuantity {

		menuItem, err := repositories.GetMenuItemByNumber(stock.MenuItemNumber)

		if err != nil {
			return err
		}

		if menuItem.Name == "" {
			return errors.New("invalid menu item ID")
		}

		if menuItem.Type == constants.TYPE_ASSEMBLED {
			err = assembleItem(transaction, currentUserNumber, stock, stockingRequest.CampusNumber)
		} else {
			_, err = saveNonAssembled(transaction, stockingRequest.CampusNumber, stock.UnitCost, menuItem.Number, stock.Quantity, currentUserNumber)
		}

		return err
	}

	return nil

}

func assembleItem(transaction *gorm.DB, currentUserNumber string, stock requests.ItemPriceAndQuantity, campusNumber string) error {

	itemAssembles, err := repositories.GetAssembledItemParts(stock.MenuItemNumber)
	if err != nil {
		return err
	}

	if len(itemAssembles) == 0 {
		return errors.New("item does not have parts set")
	}

	//Check stock for each part item
	var stockReducingParts []models.Stocking

	for _, itemAssemble := range itemAssembles {

		secondaryItem := itemAssemble.SecondaryItem

		requiredQty := itemAssemble.Quantity * stock.Quantity
		availableQty, err := repositories.GetMenuItemCurrentStock(secondaryItem.Number)

		if err != nil {
			return err
		}

		if availableQty < requiredQty {
			return errors.New("Not enough available quantity for: " + secondaryItem.Name + " " + strconv.FormatFloat(requiredQty, 'f', 2, 64) + " is required but we have " + strconv.FormatFloat(availableQty, 'f', 2, 64))
		}

		var mostRecentStockUnitCost float64

		mostRecentStockUnitCost = 1200

		negativeStock, err := saveNonAssembled(transaction, campusNumber, mostRecentStockUnitCost, itemAssemble.SecondaryItemNumber, -1*requiredQty, currentUserNumber)
		if err != nil {
			return err
		}

		stockReducingParts = append(stockReducingParts, negativeStock)

	}

	newStock, err := saveNonAssembled(transaction, campusNumber, stock.UnitCost, stock.MenuItemNumber, stock.Quantity, currentUserNumber)

	if err != nil {
		return err
	}

	//Create stocking entries that reduced the stock to create new stock
	for _, stockReducingPart := range stockReducingParts {

		var assembledItem = models.AssembledPart{
			FromStockingNumber: stockReducingPart.Number,
			NewStockingNumber:  newStock.Number,
		}

		err = repositories.SaveAssembledPart(assembledItem)
		if err != nil {
			return err
		}

	}

	return nil
}

func saveNonAssembled(transaction *gorm.DB, campusNumber string, unitCost float64, menuItem string, quantity float64, currentUserNumber string) (models.Stocking, error) {

	saveStock := models.Stocking{
		Number:         uuid.NewV4().String(),
		CampusNumber:   campusNumber,
		UnitCost:       unitCost,
		MenuItemNumber: menuItem,
		Quantity:       quantity,
		UserNumber:     currentUserNumber,
	}

	if _, err := repositories.CreateStockingEntry(transaction, saveStock); err != nil {
		return saveStock, err
	}

	return saveStock, nil
}
