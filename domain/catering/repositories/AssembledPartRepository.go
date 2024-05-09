package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
)

func SaveAssembledPart(assembledPart models.AssembledPart) error {

	db := config.DB

	if err := db.Create(&assembledPart).Error; err != nil {
		return err
	}

	return nil

}
