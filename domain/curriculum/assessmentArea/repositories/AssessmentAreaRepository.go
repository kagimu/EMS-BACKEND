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
