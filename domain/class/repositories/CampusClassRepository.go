package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
)

func QueryCampusClassByNumber(campusClassNumber string) (campusClassModel.CampusClass, error) {

	var campusClass campusClassModel.CampusClass

	db := config.DB

	err := db.Where("number = ? ", campusClassNumber).First(campusClass).Error

	if err != nil {
		return campusClass, err
	}

	return campusClass, nil
}
