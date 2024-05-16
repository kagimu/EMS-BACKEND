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

func GetFeedingSessionForCampusAndMealCourseAndDay(day string, mealCourseNumber string, campusNumber string) (session models.MealCourseFeedingSession, err error) {

	db := config.DB

	var mealCourseFeedingSession models.MealCourseFeedingSession

	if err := db.
		Joins("JOIN campus_meal_courses ON meal_course_feeding_sessions.number = campus_meal_courses.meal_course_feeding_session_number").
		Joins("JOIN meal_courses ON meal_course_feeding_sessions.meal_course_number = meal_courses.number").
		Where("campus_meal_courses.day = ? ", day).
		Where("meal_courses.number = ? ", mealCourseNumber).
		Where("campus_meal_courses.campus_number = ? ", campusNumber).
		Preload("FeedingSession").
		First(&mealCourseFeedingSession).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourseFeedingSession, nil
		}

		return mealCourseFeedingSession, err
	}

	return mealCourseFeedingSession, nil
}
