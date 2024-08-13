package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	schoolCurriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

func FindSchoolCurriculumBySchoolNumber(schoolNumber string, tx ...*gorm.DB) ([]schoolCurriculumModel.SchoolCurriculum, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var curriculums []schoolCurriculumModel.SchoolCurriculum

	if err := db.Where("school_number = ?", schoolNumber).Find(&curriculums).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculums, nil
		}

		return curriculums, err
	}

	return curriculums, nil
}

func RemoveSchoolCurriculum(schoolCurriculum schoolCurriculumModel.SchoolCurriculum, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Delete(&schoolCurriculum).Error
}
