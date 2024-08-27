package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	curriculumLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

func GetCurriculumLevelByNumber(curriculumLevelNumber string, txn ...*gorm.DB) (curriculumLevelModels.CurriculumLevel, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var curriculumLevel curriculumLevelModels.CurriculumLevel

	if err := db.Where("number = ?", curriculumLevelNumber).First(&curriculumLevel).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculumLevel, nil
		}

		return curriculumLevel, err
	}

	return curriculumLevel, nil
}
