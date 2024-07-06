package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"gorm.io/gorm"
)

func SaveAssessmentArea(assessmentArea assessmentAreaModels.AssessmentArea, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&assessmentArea).Error; err != nil {
		return assessmentArea, err
	}

	return assessmentArea, nil
}

func FindAssessmentAreaBySubjectNumberAndCurriculumClassNumber(subjectNumber string, curriculumClassNUmber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ? AND subject_number = ? ", curriculumClassNUmber, subjectNumber).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil
}

func FindAssessmentAreaByNameAndCurriculumClassNumber(otherName string, curriculumClassNUmber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ? AND name = ? ", curriculumClassNUmber, otherName).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil
}

func GetAssessmentAreaByNumber(assessmentAreaNumber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("number = ?", assessmentAreaNumber).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil

}

func FindAssessmentAreasByCurriculumClassNumber(curriculumClassNumber string, trn ...*gorm.DB) ([]assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentAreas []assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ?", curriculumClassNumber).Find(&assessmentAreas).Error; err != nil {

		return assessmentAreas, err
	}

	return assessmentAreas, nil
}
