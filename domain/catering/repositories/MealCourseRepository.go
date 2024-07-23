package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func SaveMealCourse(mealCourse models.MealCourse, tx ...*gorm.DB) (models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&mealCourse).Error; err != nil {
		return models.MealCourse{}, err
	}

	return mealCourse, nil
}

func GetAllMealCoursesByMealCourseCount(mealCourseCountNumber, status string, tx ...*gorm.DB) ([]models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var mealCourses []models.MealCourse

	if err := db.Scopes(filters.ByStatus(status)).Where("meal_course_count_number = ? ", mealCourseCountNumber).Order("code ASC").Find(&mealCourses).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourses, nil
		}

		return mealCourses, err
	}

	return mealCourses, nil
}

func GetMealCourseByName(name, mealCourseCountNumber string, tx ...*gorm.DB) (models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

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

func GetMealCourseByNumber(mealCourseNumber string, tx ...*gorm.DB) (models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

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

func GetMealCourseByNumberAndSchool(mealCourseNumber string, schoolNumber string, tx ...*gorm.DB) (models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var mealCourse models.MealCourse

	if err := db.
		Joins("JOIN meal_course_counts ON meal_courses.meal_course_count_number = meal_course_counts.number").
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

func GetMealCourseByCode(code int, mealCourseCountNumber string, tx ...*gorm.DB) (models.MealCourse, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var mealCourse models.MealCourse

	if err := db.
		Where("code = ?", code).
		Where("meal_course_count_number = ?", mealCourseCountNumber).
		First(&mealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mealCourse, nil
		}

		return mealCourse, err
	}

	return mealCourse, nil
}
