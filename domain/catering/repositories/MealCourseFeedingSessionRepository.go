package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func SaveMealCourseFeedingSession(mealCourseFeedingSession models.MealCourseFeedingSession) error {

	db := config.DB

	if err := db.Create(&mealCourseFeedingSession).Error; err != nil {
		return err
	}

	return nil

}

func GetMealCourseFeedingSession(mealCourseNumber string, feedingSessionNumber string) (models.MealCourseFeedingSession, error) {

	db := config.DB

	var mealCourseFeedingSession models.MealCourseFeedingSession

	if err := db.
		Where("meal_course_number = ?", mealCourseNumber).
		Where("feeding_session_number = ?", feedingSessionNumber).
		First(&mealCourseFeedingSession).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourseFeedingSession, nil
		}

		return mealCourseFeedingSession, err
	}

	return mealCourseFeedingSession, nil
}

func GetMealCourseFeedingSessionsByMealCourseNumber(mealCourseNumber string) ([]models.MealCourseFeedingSession, error) {

	db := config.DB

	var mealCourseFeedingSession []models.MealCourseFeedingSession

	if err := db.
		Where("meal_course_number = ?", mealCourseNumber).
		Preload("FeedingSession").
		Find(&mealCourseFeedingSession).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourseFeedingSession, nil
		}

		return mealCourseFeedingSession, err
	}

	return mealCourseFeedingSession, nil
}
