package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

func GetLevelByNumber(levelNumber string, trn ...*gorm.DB) (levelModels.Level, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var level levelModels.Level

	if err := db.Where("number = ?", levelNumber).First(&level).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return levelModels.Level{}, nil
		}

		return levelModels.Level{}, err
	}

	return level, nil
}
