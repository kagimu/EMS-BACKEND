package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func SaveFeedingGroupMealCourse(feedingGroupMealCourse models.FeedingGroupMealCourse) error {

	db := config.DB

	if err := db.Create(&feedingGroupMealCourse).Error; err != nil {
		return err
	}

	return nil

}

func GetAllFeedingGroupMealCoursesForDay(feedingGroupNumber string, day string) ([]models.FeedingGroupMealCourse, error) {

	db := config.DB
	var feedingGroupMealCourse []models.FeedingGroupMealCourse

	if err := db.
		Preload("MealCourse").
		Where("feeding_group_number = ?", feedingGroupNumber).
		Where("day = ?", day).
		Find(&feedingGroupMealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupMealCourse, nil
		}

		return feedingGroupMealCourse, err
	}

	return feedingGroupMealCourse, nil

}

func GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber string, day string, feedingGroupNumber string) (models.FeedingGroupMealCourse, error) {

	db := config.DB
	var feedingGroupMealCourse models.FeedingGroupMealCourse

	if err := db.Where("feeding_group_number = ?", feedingGroupNumber).
		Where("day = ?", day).
		Where("meal_course_number = ?", mealCourseNumber).
		First(&feedingGroupMealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupMealCourse, nil
		}

		return feedingGroupMealCourse, err
	}

	return feedingGroupMealCourse, nil
}

func GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber string, day string) ([]models.FeedingGroupMealCourse, error) {

	db := config.DB
	var feedingGroupMealCourse []models.FeedingGroupMealCourse

	if err := db.
		Preload("FeedingGroup").
		Where("day = ?", day).
		Where("meal_course_number = ?", mealCourseNumber).
		Find(&feedingGroupMealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupMealCourse, nil
		}

		return feedingGroupMealCourse, err
	}

	return feedingGroupMealCourse, nil
}
