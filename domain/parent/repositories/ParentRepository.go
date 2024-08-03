package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	"gorm.io/gorm"
)

func GetParentByApplicationNumber(applicationNumber string, tx ...*gorm.DB) (models.Parent, error) {
	
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var parent models.Parent

	if err := db.Where("application_number = ?", applicationNumber).First(&parent).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Parent{}, nil
		}

		return models.Parent{}, err
	}

	return parent, nil
}
