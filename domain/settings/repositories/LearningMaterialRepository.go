package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	learningMaterialModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	"gorm.io/gorm"
)

func FindLearningMaterialByNumber(learningMaterialNumber string, trn ...*gorm.DB) (learningMaterialModel.LearningMaterial, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var learningMaterial learningMaterialModel.LearningMaterial

	if err := db.Where("number = ?", learningMaterialNumber).First(&learningMaterial).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return learningMaterial, nil
		}

		return learningMaterial, err
	}

	return learningMaterial, nil
}
