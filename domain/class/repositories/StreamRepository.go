package repositories

import (
	"errors"
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

func FetchStreamByCampusClassNumber(campusClassNumber string, trn ...*gorm.DB) ([]models.Stream, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var streams []models.Stream

	if err := db.Where("class_level_number = ?", campusClassNumber).Find(&streams).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return streams, nil
		}

		return streams, err

	}

	return streams, nil
}
