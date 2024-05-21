package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetMealStockByMealNumberAndMenuItemNumber(mealNumber string, menuItemNumber string) (models.Stocking, error) {
	db := config.DB

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

func SaveMealStocking(tx *gorm.DB, stock models.MealStock) error {

	if err := tx.Save(&stock).Error; err != nil {

		if err = tx.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	return nil
}

func GetMealStockingByMealNumberAndStockingNumber(mealNumber string, stockingNumber string) (models.MealStock, error) {

	db := config.DB

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
