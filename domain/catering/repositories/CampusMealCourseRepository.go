package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetCampusMealCourseByCampusAndDay(day string, mealCourseFeedingSessionNumber string, campusNumber string, tx ...*gorm.DB) (models.CampusMealCourse, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var campusMealCourse models.CampusMealCourse

	if err := db.Where("day = ?", day).
		Where("campus_number = ?", campusNumber).
		Where("meal_course_feeding_session_number = ?", mealCourseFeedingSessionNumber).
		First(&campusMealCourse).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campusMealCourse, nil
		}

		return campusMealCourse, err
	}

	return campusMealCourse, nil
}

func SaveCampusMealCourse(campusMealCourse models.CampusMealCourse, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Create(&campusMealCourse).Error; err != nil {
		return err
	}

	return nil

}

func GetCampusMealCourseByCampus(mealCourseFeedingSessionNumber string, campusNumber string, tx ...*gorm.DB) ([]models.CampusMealCourse, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var campusMealCourses []models.CampusMealCourse

	if err := db.Where("campus_number = ?", campusNumber).
		Where("meal_course_feeding_session_number = ?", mealCourseFeedingSessionNumber).
		Find(&campusMealCourses).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campusMealCourses, nil
		}

		return campusMealCourses, err
	}

	return campusMealCourses, nil
}

func DeleteCampusMealCourse(campusMealCourse models.CampusMealCourse, tx ...*gorm.DB) error {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Create(&campusMealCourse).Error; err != nil {
		return err
	}

	return nil

}
