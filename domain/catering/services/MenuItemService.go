package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

func GetMenuItemParts(menuItemNumber string) ([]responses.Part, error) {

	var parts = make([]responses.Part, 0)

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
			//return nil, errors.New("unit of measure number is empty")
			unitOfMeasure = models.UnitOfMeasure{}
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

	err := db.Transaction(func(transaction *gorm.DB) error {

		err := doRunOnTransaction(transaction, stockingRequest, currentUserNumber)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func doRunOnTransaction(transaction *gorm.DB, stockingRequest requests.StockingRequest, currentUserNumber string) error {

	for _, stock := range stockingRequest.ItemPriceAndQuantity {

		menuItem, err := repositories.GetMenuItemByNumber(stock.MenuItemNumber, transaction)

		if err != nil {
			return errors.New("getting menu item by number error: " + err.Error())
		}

		if menuItem.Name == "" {
			return errors.New("invalid menu item ID")
		}

		if menuItem.Type == constants.TypeAssembled {
			err = assembleItem(transaction, currentUserNumber, stock, stockingRequest.CampusNumber)
		} else {
			_, err = saveNonAssembled(transaction, stockingRequest.CampusNumber, stock.UnitCost, menuItem.Number, stock.Quantity, currentUserNumber)
		}

		return err
	}

	return nil

}

func assembleItem(transaction *gorm.DB, currentUserNumber string, stock requests.ItemPriceAndQuantity, campusNumber string) error {

	itemAssembles, err := repositories.GetAssembledItemParts(stock.MenuItemNumber, transaction)
	if err != nil {
		return errors.New("getting assembled parts error: " + err.Error())
	}

	if len(itemAssembles) == 0 {
		return errors.New("item does not have parts set")
	}

	//Check stock for each part item
	var stockReducingParts []models.Stocking

	var unitCost float64

	for _, itemAssemble := range itemAssembles {

		secondaryItem := itemAssemble.SecondaryItem

		requiredQty := itemAssemble.Quantity * stock.Quantity
		availableQty, err := repositories.GetMenuItemCurrentStock(secondaryItem.Number, transaction)

		if err != nil {
			return errors.New("getting item current stock error: " + err.Error())
		}

		if availableQty < requiredQty {
			return errors.New("Not enough available quantity for: " + secondaryItem.Name + " " + strconv.FormatFloat(requiredQty, 'f', 2, 64) + " is required but we have " + strconv.FormatFloat(availableQty, 'f', 2, 64))
		}

		recentStock, err := repositories.GetRecentItemStock(itemAssemble.SecondaryItemNumber, transaction)

		if err != nil {
			return errors.New("getting recent stock for " + secondaryItem.Name + " error: " + err.Error())
		}

		if recentStock.Number == "" {
			return errors.New("no recent stock for " + itemAssemble.SecondaryItemName)
		}

		unitCost = recentStock.UnitCost + unitCost

		negativeStock, err := saveNonAssembled(transaction, campusNumber, recentStock.UnitCost, itemAssemble.SecondaryItemNumber, -1*requiredQty, currentUserNumber)
		if err != nil {
			return errors.New("first save of none assembled item error: " + err.Error())
		}

		stockReducingParts = append(stockReducingParts, negativeStock)

	}

	newStock, err := saveNonAssembled(transaction, campusNumber, unitCost, stock.MenuItemNumber, stock.Quantity, currentUserNumber)

	if err != nil {
		return errors.New("saving none assembled item error: " + err.Error())
	}

	//Create stocking entries that reduced the stock to create new stock
	for _, stockReducingPart := range stockReducingParts {

		var assembledItem = models.AssembledPart{
			FromStockingNumber: stockReducingPart.Number,
			NewStockingNumber:  newStock.Number,
		}

		err = repositories.SaveAssembledPart(assembledItem, transaction)
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

	if _, err := repositories.CreateStockingEntry(saveStock, transaction); err != nil {
		return saveStock, err
	}

	return saveStock, nil
}

func MenuProjectedCalender(startDate time.Time, campusNumber string, endDate time.Time) (responses.ProjectionCalenderResponse, error) {

	var projectionCalenderResponse responses.ProjectionCalenderResponse

	campus, err := repositories2.GetCampusByNumber(campusNumber)
	if err != nil {
		return projectionCalenderResponse, err
	}

	schoolNumber := campus.SchoolNumber
	mealCourseCount, err := repositories.GetMealCourseCountBySchoolNumber(schoolNumber)
	if err != nil {
		return projectionCalenderResponse, err
	}

	if mealCourseCount.Number == "" {
		return projectionCalenderResponse, errors.New("the school haven't set up meals")
	}

	mealCourses, err := repositories.GetAllMealCoursesByMealCourseCount(mealCourseCount.Number, constants.StatusActive)
	if err != nil {
		return projectionCalenderResponse, err
	}

	calenderDays := make([]responses.ProjectionCalender, 0)

	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {

		_date := date.Format(time.DateOnly)
		day := date.Format("Mon")

		calendarItem := responses.ProjectionCalender{
			Date: _date,
			Day:  day,
		}

		calendarItem.Meals = getMealItemsAndQuantitiesPerDayForAllUsers(day, mealCourses, campusNumber)

		calenderDays = append(calenderDays, calendarItem)
	}

	projectionCalenderResponse.CalenderItems = calenderDays

	var summary = make([]responses.ProjectionMenuItem, 0)
	var totalCost float64

	if len(calenderDays) > 0 {
		summary, totalCost = buildProjectionSummary(calenderDays)
	}

	projectionCalenderResponse.ItemsSummary = summary

	projectionCalenderResponse.TotalCost = totalCost

	return projectionCalenderResponse, nil
}

func buildProjectionSummary(calenderDays []responses.ProjectionCalender) ([]responses.ProjectionMenuItem, float64) {

	var items []responses.ProjectionMenuItem

	var totalCost float64

	for _, calendarItem := range calenderDays {

		for _, meal := range calendarItem.Meals {

			for _, menuItem := range meal.Items {
				items = updateItemOrAppend(items, menuItem)
				totalCost += menuItem.TotalCost
			}

		}

	}

	return items, totalCost

}

func updateItemOrAppend(items []responses.ProjectionMenuItem, newItem responses.ProjectionMenuItem) []responses.ProjectionMenuItem {

	for index, item := range items {

		if item.MenuItemNumber == newItem.MenuItemNumber {
			items[index].Quantity += newItem.Quantity
			items[index].TotalCost += newItem.TotalCost
			items[index].PurchaseQuantity += newItem.RateToPurchaseUnit * newItem.Quantity
			return items
		}

	}

	return append(items, newItem)
}

func getMealItemsAndQuantitiesPerDayForAllUsers(day string, courses []models.MealCourse, campusNumber string) []responses.ProjectionCalenderMeal {

	meals := make([]responses.ProjectionCalenderMeal, 0)

	for _, course := range courses {

		var mealProjection responses.ProjectionCalenderMeal
		mealProjection.MealName = course.Name

		dayItems, err := GetUserMeals(course.Number, campusNumber, "", day)
		if err != nil {
			log.Println("No meals for day " + day + " on " + course.Number + ", error: " + err.Error())
		} else if len(dayItems) > 0 {

			totalUsers := len(dayItems)
			mealProjection.NumberOfPeople = totalUsers
			getItems := dayItems[0].Meals //Since they are all served the same items for the meal. Getting items at index zero is enough

			items := make([]responses.ProjectionMenuItem, 0)

			for _, _item := range getItems {

				_baseItem, err := repositories.GetMenuItemByNumber(_item.MenuItemNumber)
				if err != nil {
					fmt.Println("Getting base price item error: " + err.Error())
					continue
				}

				_totalQuantity := _item.Quantity * float64(totalUsers)

				if _baseItem.Type == constants.TypeAssembled {

					items = splitAssembledItem(_baseItem, items, _totalQuantity)

				} else {
					items = append(items, createProjectionMenuItem(_baseItem, _totalQuantity))
				}

			}

			mealProjection.Items = items

		}

		meals = append(meals, mealProjection)

	}

	return meals
}

func createProjectionMenuItem(menuItem models.MenuItem, totalQuantity float64) responses.ProjectionMenuItem {

	projectionMenuItem := responses.ProjectionMenuItem{
		Name:               menuItem.Name,
		Quantity:           totalQuantity,
		RateToPurchaseUnit: menuItem.ConversionToPurchaseUnit,
		PurchaseQuantity:   totalQuantity * menuItem.ConversionToPurchaseUnit,
		MenuItemNumber:     menuItem.Number,
		TotalCost:          totalQuantity * determineUnitCost(menuItem),
	}

	unitOfMeasure, err := repositories.GetUnitOfMeasureByNumber(menuItem.UnitOfMeasureNumber)
	if err != nil {
		fmt.Println("Error getting unit of measure: " + err.Error())
		return projectionMenuItem
	}

	projectionMenuItem.UnitOfMeasure = unitOfMeasure.Abbreviation

	if menuItem.PurchaseUnitOfMeasureNumber != "" {
		unitOfMeasure, err = repositories.GetUnitOfMeasureByNumber(menuItem.PurchaseUnitOfMeasureNumber)
		if err != nil {
			fmt.Println("Error getting purchase unit of measure: " + err.Error())
			return projectionMenuItem
		}

		projectionMenuItem.PurchaseUnitOfMeasure = unitOfMeasure.Abbreviation
		projectionMenuItem.RateToPurchaseUnit = menuItem.ConversionToPurchaseUnit
	}

	return projectionMenuItem
}

func splitAssembledItem(item models.MenuItem, items []responses.ProjectionMenuItem, requiredQty float64) []responses.ProjectionMenuItem {

	parts, err := repositories.GetAssembledItemParts(item.Number)
	if err != nil {
		fmt.Println("Getting items parts error: " + err.Error())
		return items
	}

	for _, part := range parts {

		partsQty := part.Quantity * requiredQty

		menuItem, err := repositories.GetMenuItemByNumber(part.SecondaryItemNumber)
		if err != nil {
			fmt.Println("Getting base price item error: " + err.Error())
			continue
		}

		if menuItem.Type == constants.TypeAssembled {
			items = splitAssembledItem(menuItem, items, partsQty)
		} else {

			items = append(items, createProjectionMenuItem(menuItem, partsQty))

		}

	}

	return items
}

func determineUnitCost(item models.MenuItem) float64 {

	if item.Number == "" {
		return 0.0
	}

	if item.Type == constants.TypeItem {
		return item.BaseUnitCost
	}

	var unitCost float64
	parts, err := repositories.GetAssembledItemParts(item.Number)
	if err != nil {
		log.Println("Getting items parts for unit cost error: " + err.Error())
		return 0.0
	}

	for _, part := range parts {
		partItem, err := repositories.GetMenuItemByNumber(part.SecondaryItemNumber)
		if err != nil {
			continue
		}
		unitCost += part.Quantity * determineUnitCost(partItem)
	}

	return unitCost
}

func SaveMenuItem(itemType string, payload requests.MenuItemRequest, currentUserNumber string) (models.MenuItem, error) {

	school, err := repositories3.GetSchoolByNumber(payload.SchoolNumber)

	if err != nil {
		return models.MenuItem{}, errors.New("checking school number error: " + err.Error())
	}

	checkName, err := repositories.GetMenuItemByNameAndSchool(school.Number, payload.Name)
	if err != nil {
		return models.MenuItem{}, errors.New("checking menu item name error: " + err.Error())
	}

	if checkName.Number != "" {
		return models.MenuItem{}, errors.New("menu item name '" + payload.Name + "' already exists")
	}

	checkCode, err := repositories.GetMenuItemByCodeAndSchool(school.Number, payload.Code)
	if err != nil {
		return models.MenuItem{}, errors.New("checking menu item code error: " + err.Error())
	}

	if checkCode.Number != "" {
		return models.MenuItem{}, errors.New("item code already used")
	}

	saveMenuItem := models.MenuItem{
		Number:                      uuid.NewV4().String(),
		Name:                        payload.Name,
		Code:                        payload.Code,
		Type:                        itemType,
		SchoolNumber:                school.Number,
		UnitOfMeasureNumber:         payload.UnitOfMeasureNumber,
		UserNumber:                  currentUserNumber,
		BaseUnitCost:                payload.BasePrice,
		Status:                      constants.StatusActive,
		PurchaseUnitOfMeasureNumber: payload.PurchaseUnitOfMeasureNumber,
		ConversionToPurchaseUnit:    payload.ConversionToPurchaseUnit,
	}

	err = repositories.SaveMenuItem(saveMenuItem)
	if err != nil {
		return models.MenuItem{}, errors.New("saving menu item error: " + err.Error())
	}

	return saveMenuItem, nil
}

func GetPaginatedMenuItems(schoolNumber string, campusNumber string, status string, itemType string, context *gin.Context) (pagination.PaginatedResponse, error) {

	response := make([]responses.MenuItemResponse, 0)

	var paginatedResponse pagination.PaginatedResponse
	paginatedResponse.Content = response

	if campusNumber != "" {
		campus, err := repositories2.GetCampusByNumber(campusNumber)
		if err != nil {
			return paginatedResponse, errors.New("validating campus number error: " + err.Error())
		}

		if campus.Number == "" {
			return paginatedResponse, errors.New("invalid campus number: " + campusNumber)
		}

		schoolNumber = campus.SchoolNumber
	}

	menuItems, err := repositories.GetPaginatedMenuItems(schoolNumber, status, itemType, paginatedResponse, context)

	if err != nil {
		return pagination.PaginatedResponse{}, errors.New("getting all menu items error: " + err.Error())
	}

	for _, menuItem := range menuItems {

		itemResp, err := createMenuItemResponse(menuItem)
		if err != nil {
			return pagination.PaginatedResponse{}, err
		}

		response = append(response, itemResp)

	}

	paginatedResponse.Content = response

	return paginatedResponse, nil
}

func GetAllMenuItems(schoolNumber string, status string, itemType string) ([]responses.MenuItemResponse, error) {

	response := make([]responses.MenuItemResponse, 0)

	menuItems, err := repositories.GetAllMenuItems(schoolNumber, status, itemType)

	for _, menuItem := range menuItems {

		itemResp, err := createMenuItemResponse(menuItem)
		if err != nil {
			return response, err
		}

		response = append(response, itemResp)

	}

	return response, err
}

func createMenuItemResponse(menuItem models.MenuItem) (responses.MenuItemResponse, error) {

	currentQuantity, err := repositories.GetMenuItemCurrentStock(menuItem.Number)
	if err != nil {
		return responses.MenuItemResponse{}, errors.New("getting current stock error: " + err.Error())
	}

	parts, err := GetMenuItemParts(menuItem.Number)
	if err != nil {
		return responses.MenuItemResponse{}, errors.New("getting parts error: " + err.Error())
	}

	purchaseUnitOfMeasure := models.UnitOfMeasure{}
	if menuItem.PurchaseUnitOfMeasureNumber != "" {
		purchaseUnitOfMeasure, err = repositories.GetUnitOfMeasureByNumber(menuItem.PurchaseUnitOfMeasureNumber)
		if err != nil {
			return responses.MenuItemResponse{}, errors.New("getting purchase unit of measure error: " + err.Error())
		}
	}

	return responses.MenuItemResponse{
		Number:                      menuItem.Number,
		Name:                        menuItem.Name,
		Code:                        menuItem.Code,
		Type:                        menuItem.Type,
		SchoolName:                  menuItem.School.Name,
		SchoolNumber:                menuItem.SchoolNumber,
		UnitOfMeasureNumber:         menuItem.UnitOfMeasureNumber,
		UnitOfMeasureName:           menuItem.UnitOfMeasure.Name,
		CurrentQuantity:             currentQuantity,
		PurchaseUnitOfMeasureNumber: purchaseUnitOfMeasure.Number,
		PurchaseUnitOfMeasureName:   purchaseUnitOfMeasure.Name,
		Parts:                       parts,
		UnitCost:                    determineUnitCost(menuItem),
	}, nil
}

func MenuItemDetails(menuItemNumber string) (responses.MenuItemResponse, error) {

	item, err := repositories.GetMenuItemByNumber(menuItemNumber)

	if err != nil {
		return responses.MenuItemResponse{}, errors.New("getting menu item by number error: " + err.Error())
	}

	if item.Number == "" {
		return responses.MenuItemResponse{}, errors.New("invalid menu item number" + menuItemNumber)
	}

	return createMenuItemResponse(item)
}

func UpdateMenuItem(payload requests.UpdateMenuItemRequest, currentUserNumber string) error {

	menuItem, err := repositories.GetMenuItemByNumber(payload.MenuItemNumber)
	if err != nil {
		return errors.New("checking menu item number error: " + err.Error())
	}

	if menuItem.Number == "" {
		return errors.New("invalid menu item number" + payload.MenuItemNumber)
	}

	school, err := repositories3.GetSchoolByNumber(payload.SchoolNumber)

	if err != nil {
		return errors.New("checking school number error: " + err.Error())
	}

	checkName, err := repositories.GetMenuItemByNameAndSchool(school.Number, payload.Name)
	if err != nil {
		return errors.New("checking menu item name error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != menuItem.Number {
		return errors.New("menu item name '" + payload.Name + "' already exists")
	}

	checkCode, err := repositories.GetMenuItemByCodeAndSchool(school.Number, payload.Code)
	if err != nil {
		return errors.New("checking menu item code error: " + err.Error())
	}

	if checkCode.Number != "" && checkCode.Number != menuItem.Number {
		return errors.New("item code already used")
	}

	menuItem.Status = payload.Status
	menuItem.Name = payload.Name
	menuItem.Code = payload.Code
	menuItem.UnitOfMeasureNumber = payload.UnitOfMeasureNumber
	menuItem.UserNumber = currentUserNumber
	menuItem.BaseUnitCost = payload.BasePrice
	menuItem.PurchaseUnitOfMeasureNumber = payload.PurchaseUnitOfMeasureNumber
	menuItem.ConversionToPurchaseUnit = payload.ConversionToPurchaseUnit

	err = repositories.SaveMenuItem(menuItem)
	if err != nil {
		return errors.New("updating menu item error: " + err.Error())
	}

	return nil
}

func DeleteMenuItem(menuItemNumber string) error {
	menuItem, err := repositories.GetMenuItemByNumber(menuItemNumber)
	if err != nil {
		return errors.New("checking menu item number error: " + err.Error())
	}

	if menuItem.Number == "" {
		return errors.New("invalid menu item number" + menuItemNumber)
	}

	parts, err := repositories.GetAssembledItemParts(menuItem.Number)
	if err != nil {
		return errors.New("checking menu item parts before delete error: " + err.Error())
	}

	tx := config.DB.Begin()

	for _, part := range parts {
		err = repositories.DeleteItemAssembly(part, tx)
		if err != nil {
			return errors.New("clearing assembled parts error: " + err.Error())
		}
	}

	err = repositories.DeleteMenuItem(menuItem, tx)
	if err != nil {
		fmt.Println("delete menu item error: " + err.Error())
		return errors.New("item can not be removed, you may need to delete all other dependencies first")
	}

	tx.Commit()

	return nil
}

func CreateItemAssembly(currentUser string, payload requests.ItemAssembly) error {

	menuItem, err := SaveMenuItem(constants.TypeAssembled, payload.MenuItemRequest, currentUser)
	if err != nil {
		return errors.New("saving assembly menu item error: " + err.Error())
	}

	if menuItem.Number == "" {
		return errors.New("saved assembly item returned no number")
	}

	for _, secondaryItem := range payload.SecondaryItemQuantity {

		duplicate, err := repositories.FindAssemblyItemByPrimaryNumberAndSecondaryNumber(menuItem.Number, secondaryItem.SecondaryItemNumber)
		if err != nil {
			fmt.Println("checking assembly error: " + err.Error())
			continue
		}

		if duplicate.Number != "" {
			continue
		}

		saveItemAssembly := models.ItemAssembly{
			Number:              uuid.NewV4().String(),
			PrimaryItemNumber:   menuItem.Number,
			SecondaryItemNumber: secondaryItem.SecondaryItemNumber,
			Quantity:            secondaryItem.Quantity,
			UserNumber:          currentUser,
		}

		err = repositories.SaveItemAssembly(saveItemAssembly)

	}

	return nil
}

func UpdateItemAssembly(currentUser string, payload requests.UpdateItemAssembly) error {

	err := UpdateMenuItem(payload.UpdateMenuItemRequest, currentUser)
	if err != nil {
		return errors.New("updating assembled menu item error: " + err.Error())
	}

	oldParts, err := repositories.GetAssembledItemParts(payload.MenuItemNumber)
	if err != nil {
		return errors.New("getting assembled item parts error: " + err.Error())
	}

	for _, oldPart := range oldParts {
		err = repositories.DeleteItemAssembly(oldPart)
		if err != nil {
			fmt.Println("deleting assembled item error: " + err.Error())
			continue
		}
	}

	for _, secondaryItem := range payload.SecondaryItemQuantity {

		duplicate, err := repositories.FindAssemblyItemByPrimaryNumberAndSecondaryNumber(payload.MenuItemNumber, secondaryItem.SecondaryItemNumber)
		if err != nil {
			fmt.Println("checking assembly error on update: " + err.Error())
			continue
		}

		if duplicate.Number != "" {
			continue
		}

		saveItemAssembly := models.ItemAssembly{
			Number:              uuid.NewV4().String(),
			PrimaryItemNumber:   payload.MenuItemNumber,
			SecondaryItemNumber: secondaryItem.SecondaryItemNumber,
			Quantity:            secondaryItem.Quantity,
			UserNumber:          currentUser,
		}

		err = repositories.SaveItemAssembly(saveItemAssembly)

	}

	return nil
}
