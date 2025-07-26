package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	learningAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	"gorm.io/gorm"
)

func FindLearningAreasByCurriculumNumber(curriculumNumber string, txn ...*gorm.DB) ([]learningAreaModel.LearningArea, error) {

	var learningAreas []learningAreaModel.LearningArea

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	if err := db.Where("curriculum_number = ?", curriculumNumber).Find(&learningAreas).Error; err != nil {
		return learningAreas, err
	}

	return learningAreas, nil
}

func GetLearningAreaByNumber(learningAreaNumber string, trn ...*gorm.DB) (learningAreaModel.LearningArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var learningArea learningAreaModel.LearningArea

	if err := db.Where("number = ?", learningAreaNumber).Unscoped().First(&learningArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return learningArea, nil
		}

		return learningArea, err
	}

	return learningArea, nil

}
