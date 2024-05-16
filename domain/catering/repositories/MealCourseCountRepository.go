package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func SaveMealCourseCount(mealCourseCount models.MealCourseCount) (models.MealCourseCount, error) {
	db := config.DB

	if err := db.Save(&mealCourseCount).Error; err != nil {
		return models.MealCourseCount{}, err
	}

	return mealCourseCount, nil
}

func GetMealCourseCountBySchoolNumber(schoolNumber string) (models.MealCourseCount, error) {

	db := config.DB

	var mealCourseCount models.MealCourseCount

	if err := db.Where("school_number = ?", schoolNumber).First(&mealCourseCount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourseCount, nil
		}

		return mealCourseCount, err
	}

	return mealCourseCount, nil
}
