package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
	"time"
)

func SaveMeal(meal models.Meal, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&meal).Error
}

func GetMealByMealNumberAndDateAndUserNumber(mealCourseNumber string, date time.Time, userNumber string, tx ...*gorm.DB) (models.Meal, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var meal models.Meal

	formattedDate := date.Format(time.DateOnly)

	if err := db.
		Where("meal_course_number = ? ", mealCourseNumber).
		Where("DATE(date) = ? ", formattedDate).
		Where("user_number", userNumber).
		First(&meal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return meal, nil
		}

		return meal, err
	}

	return meal, nil
}
