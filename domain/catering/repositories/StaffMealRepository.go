package repositories

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetStaffMealByCourseDateAndStaff(mealCourseNumber string, date time.Time, staffNumber string) (models.StaffMeal, error) {
	db := config.DB
	var meal models.StaffMeal
	formattedDate := date.Format(time.DateOnly)

	if err := db.
		Where("meal_course_number = ?", mealCourseNumber).
		Where("DATE(date) = ?", formattedDate).
		Where("staff_number = ?", staffNumber).
		First(&meal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return meal, nil
		}
		return meal, err
	}

	return meal, nil
}

func SaveStaffMeal(meal models.StaffMeal, tx *gorm.DB) error {
	db := config.DB
	if tx != nil {
		db = tx
	}

	if err := db.Save(&meal).Error; err != nil {
		return err
	}

	return nil
}

func SaveStaffMealStock(mealStock models.StaffMealStock, tx *gorm.DB) error {
	db := config.DB
	if tx != nil {
		db = tx
	}

	if err := db.Save(&mealStock).Error; err != nil {
		return err
	}

	return nil
}

func GetStaffMealStockByMealAndMenuItem(staffMealNumber string, menuItemNumber string, tx *gorm.DB) (models.Stocking, error) {
	db := config.DB
	if tx != nil {
		db = tx
	}

	var stocking models.Stocking

	if err := db.
		Joins("JOIN staff_meal_stock ON stockings.number = staff_meal_stock.stocking_number").
		Where("staff_meal_stock.staff_meal_number = ?", staffMealNumber).
		Where("stockings.menu_item_number = ?", menuItemNumber).
		First(&stocking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}
		return stocking, err
	}

	return stocking, nil
}

func GetStaffMealStockingByMealNumberAndDateStaffNumber(mealCourseNumber string, date time.Time, staffNumber string, menuItemNumber string) (models.Stocking, error) {
	db := config.DB
	var stocking models.Stocking
	formattedDate := date.Format(time.DateOnly)

	if err := db.
		Joins("JOIN staff_meal_stock ON stockings.number = staff_meal_stock.stocking_number").
		Joins("JOIN staff_meals ON staff_meal_stock.staff_meal_number = staff_meals.number").
		Where("staff_meals.meal_course_number = ?", mealCourseNumber).
		Where("DATE(staff_meals.date) = ?", formattedDate).
		Where("staff_meals.staff_number = ?", staffNumber).
		Where("stockings.menu_item_number = ?", menuItemNumber).
		First(&stocking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stocking, nil
		}
		return stocking, err
	}

	return stocking, nil
}
