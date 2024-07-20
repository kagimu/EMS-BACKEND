package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

func CreateStream(stream *models.Stream, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	if err := db.Create(&stream).Error; err != nil {
		return err
	}
	return nil
}
