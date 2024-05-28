package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetAssembledItemParts(menuItemNumber string) ([]models.ItemAssembly, error) {

	db := config.DB
	var items []models.ItemAssembly

	if err := db.Preload("SecondaryItem").Where("primary_item_number = ? ", menuItemNumber).Find(&items).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return items, nil
		}

		return items, err

	}

	return items, nil

}
