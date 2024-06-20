package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

func GetCampusByNumber(campusNumber string, tx ...*gorm.DB) (campusModel.Campus, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var campus campusModel.Campus
	if err := db.Where("number = ?", campusNumber).First(&campus).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campus, nil
		}

		return campus, err
	}
	return campus, nil
}
