package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func SaveMealCourse(mealCourse models.MealCourse) (models.MealCourse, error) {

	db := config.DB

	if err := db.Save(&mealCourse).Error; err != nil {
		return models.MealCourse{}, err
	}

	return mealCourse, nil
}

func GetAllMealCoursesByMealCourseCount(mealCourseCountNumber string, status string) ([]models.MealCourse, error) {

	db := config.DB

	var mealCourses []models.MealCourse

	if err := db.Scopes(filters.ByStatus(status)).Where("meal_course_count_number = ? ", mealCourseCountNumber).Find(&mealCourses).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourses, nil
		}

		return mealCourses, err
	}

	return mealCourses, nil
}

func GetMealCourseByName(name string, mealCourseCountNumber string) (models.MealCourse, error) {

	db := config.DB

	var mealCourse models.MealCourse

	if err := db.
		Where("name = ?", name).
		Where("meal_course_count_number = ?", mealCourseCountNumber).
		First(&mealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourse, nil
		}

		return mealCourse, err
	}

	return mealCourse, nil
}

func GetMealCourseByNumber(mealCourseNumber string) (models.MealCourse, error) {

	db := config.DB

	var mealCourse models.MealCourse

	if err := db.
		Where("number = ?", mealCourseNumber).
		First(&mealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourse, nil
		}

		return mealCourse, err
	}

	return mealCourse, nil
}

func GetMealCourseByNumberAndSchool(mealCourseNumber string, schoolNumber string) (models.MealCourse, error) {

	db := config.DB

	var mealCourse models.MealCourse

	if err := db.
		Joins("JOIN meal_courses ON meal_course_counts.number = meal_courses.meal_course_count_number").
		Where("meal_courses.number = ?", mealCourseNumber).
		Where("school_number = ?", schoolNumber).
		First(&mealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourse, nil
		}

		return mealCourse, err
	}

	return mealCourse, nil
}
