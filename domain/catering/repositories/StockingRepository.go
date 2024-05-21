package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
	"time"
)

func SaveStocking(tx *gorm.DB, stocking models.Stocking) error {

	if err := tx.Save(&stocking).Error; err != nil {

		if err = tx.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	return nil
}

func GetMealStockingByMealNumberAndDateUserNumber(mealCourseNumber string, date time.Time, userNumber string, menuItemNumber string) (models.Stocking, error) {

	db := config.DB

	var stocking models.Stocking

	if err := db.
		Joins("JOIN meal_stock ON stockings.number = meal_stock.stocking_number").
		Joins("JOIN meals ON meal_stock.meal_number = meals.number").
		Where("meals.meal_course_number = ?", mealCourseNumber).
		Where("meals.date = ?", date).
		Where("meals.user_number = ?", userNumber).
		Where("stockings.menu_item_number = ?", menuItemNumber).
		First(&stocking).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}

		return stocking, err
	}

	return stocking, nil
}
