package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
	"time"
)

func SaveStocking(stocking models.Stocking, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&stocking).Error
}

func GetMealStockingByMealNumberAndDateUserNumber(mealCourseNumber string, date time.Time, userNumber string, menuItemNumber string) (models.Stocking, error) {

	db := config.DB

	var stocking models.Stocking

	formattedDate := date.Format(time.DateOnly)

	if err := db.
		Joins("JOIN meal_stock ON stockings.number = meal_stock.stocking_number").
		Joins("JOIN meals ON meal_stock.meal_number = meals.number").
		Where("meals.meal_course_number = ?", mealCourseNumber).
		Where("DATE(meals.date) = ?", formattedDate).
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
