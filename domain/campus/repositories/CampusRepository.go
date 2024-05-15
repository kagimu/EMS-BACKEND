package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

func GetCampusByNumber(campusNumber string) (campusModel.Campus, error) {

	db := config.DB
	var campus campusModel.Campus
	if err := db.Where("number = ?", campusNumber).First(&campus).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campus, nil
		}

		return campus, err
	}
	return campus, nil
}
