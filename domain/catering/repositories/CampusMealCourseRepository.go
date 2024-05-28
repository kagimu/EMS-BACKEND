package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetCampusMealCourseByCampusAndDay(day string, mealCourseFeedingSessionNumber string, campusNumber string) (models.CampusMealCourse, error) {

	db := config.DB

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

func SaveCampusMealCourse(campusMealCourse models.CampusMealCourse) error {

	db := config.DB

	if err := db.Create(&campusMealCourse).Error; err != nil {
		return err
	}

	return nil

}
