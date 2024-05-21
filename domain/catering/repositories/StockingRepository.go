package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func SaveStocking(tx *gorm.DB, stocking models.Stocking) error {

	if err := tx.Save(&stocking).Error; err != nil {

		if err = tx.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	return nil
}
