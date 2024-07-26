package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	cateringModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetMenuItemCurrentStock(menuItemNumber string, tx ...*gorm.DB) (float64, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var totalQuantity float64

	if err := db.Model(&cateringModels.Stocking{}).Where("menu_item_number = ? ", menuItemNumber).Select("COALESCE(sum(quantity),0)").Scan(&totalQuantity).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return totalQuantity, nil
}

func SaveStockingEntry(stocking cateringModels.Stocking, tx ...*gorm.DB) (cateringModels.Stocking, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&stocking).Error; err != nil {
		return stocking, err
	}

	return stocking, nil
}

func GetRecentItemStock(menuItemNumber string, tx ...*gorm.DB) (cateringModels.Stocking, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var stocking cateringModels.Stocking

	if err := db.Where("menu_item_number = ?", menuItemNumber).
		Order("created_at desc").
		First(&stocking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}

		return stocking, err
	}

	return stocking, nil
}

func AllStocking(schoolNumber, campusNumber, menuItemNumber string, tx ...*gorm.DB) ([]cateringModels.Stocking, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var stockings []cateringModels.Stocking

	if err := db.
		Joins("join menu_items on menu_items.number = stockings.menu_item_number").
		Scopes(filters.ByAnyField("menu_item_number", menuItemNumber), filters.BySchoolNumber(schoolNumber), filters.ByCampusNumber(campusNumber)).
		Order("stockings.created_at desc").
		Find(&stockings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stockings, nil
		}

		return stockings, err
	}

	return stockings, nil
}

func GetStockItemByNumber(stockItemNumber string, tx ...*gorm.DB) (cateringModels.Stocking, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var stocking cateringModels.Stocking

	if err := db.Where("number = ?", stockItemNumber).
		First(&stocking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}

		return stocking, err
	}

	return stocking, nil
}
