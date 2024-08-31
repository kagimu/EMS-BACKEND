package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

func GetCampusLevelByNumber(campusLevelNumber string, txn ...*gorm.DB) (models.CampusLevel, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var campusLevel models.CampusLevel
	if err := db.Where("number = ?", campusLevelNumber).First(&campusLevel).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.CampusLevel{}, nil
		}

		return models.CampusLevel{}, err
	}

	return campusLevel, nil
}
