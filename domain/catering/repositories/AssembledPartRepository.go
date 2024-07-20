package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func SaveAssembledPart(assembledPart models.AssembledPart, transaction ...*gorm.DB) error {

	db := config.DB

	if len(transaction) > 0 {
		db = transaction[0]
	}

	if err := db.Create(&assembledPart).Error; err != nil {
		return err
	}

	return nil

}
