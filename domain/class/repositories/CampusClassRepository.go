package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

func GetCampusClassByNumber(campusClassNumber string) (campusClassModel.CampusClass, error) {

	var campusClass campusClassModel.CampusClass

	db := config.DB

	if err := db.Where("number = ? ", campusClassNumber).First(&campusClass).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campusClass, nil
		}

		return campusClass, err
	}

	/*if err != nil {
		return campusClass, err
	}*/

	return campusClass, nil
}
