package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
	"time"
)

func SaveMeal(db *gorm.DB, meal models.Meal) error {

	if err := db.Save(meal).Error; err != nil {

		if err = db.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	return nil
}

func GetMealByMealNumberAndDateAndUserNumber(mealCourseNumber string, date time.Time, userNumber string) (models.Meal, error) {
	db := config.DB

	var meal models.Meal

	if err := db.
		Where("meal_course_number = ? ", mealCourseNumber).
		Where("date = ? ", date).
		Where("user_number", userNumber).
		First(&meal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return meal, nil
		}

		return meal, err
	}

	return meal, nil
}
