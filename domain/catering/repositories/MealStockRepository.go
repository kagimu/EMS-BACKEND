package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetMealStockByMealNumberAndMenuItemNumber(mealNumber, menuItemNumber string, txn ...*gorm.DB) (models.Stocking, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var stocking models.Stocking

	if err := db.
		Joins("JOIN meal_stock ON stockings.number = meal_stock.stocking_number").
		Where("stockings.menu_item_number = ?", menuItemNumber).
		Where("meal_stock.meal_number", mealNumber).First(&stocking).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}

		return stocking, err
	}

	return stocking, nil
}

func SaveMealStocking(stock models.MealStock, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&stock).Error
}

func GetMealStockingByMealNumberAndStockingNumber(mealNumber, stockingNumber string, txn ...*gorm.DB) (models.MealStock, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var mealStock models.MealStock

	if err := db.Where("meal_number = ?", mealNumber).
		Where("stocking_number = ?", stockingNumber).
		First(&mealStock).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealStock, nil
		}

		return mealStock, err
	}

	return mealStock, nil
}
