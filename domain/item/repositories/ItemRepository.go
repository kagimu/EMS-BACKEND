package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

func GetItemByNumber(itemNumber string) (itemModels.Item, error) {

	db := config.DB

	var item itemModels.Item

	if err := db.Where("number = ?", itemNumber).First(&item).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, nil
		}

		return item, err
	}

	return item, nil
}
