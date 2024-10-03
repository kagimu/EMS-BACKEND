package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	"gorm.io/gorm"
)

func GetCurriculumByCampusCurriculumNumber(campusCurriculumNumber, campusNumber string, tx ...*gorm.DB) (models.Curriculum, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var curriculum models.Curriculum

	if err := db.Joins("JOIN school_curriculums ON curriculums.number = school_curriculums.curriculum_number").
		Joins("JOIN campus_curriculums ON school_curriculums.number = campus_curriculums.school_curriculum_number").
		Where("campus_curriculums.number = ? AND campus_curriculums.campus_number = ?", campusCurriculumNumber, campusNumber).
		First(&curriculum).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculum, nil
		}

		return curriculum, err
	}

	return curriculum, nil
}
