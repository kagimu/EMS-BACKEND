package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	cateringModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetItemCurrentStock(menuItemNumber string) (float64, error) {

	db := config.DB

	var totalQuantity float64

	if err := db.Model(&cateringModels.Stocking{}).Where("menu_item_number = ? ", menuItemNumber).Select("sum(quantity)").Scan(&totalQuantity).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, err
	}

	return totalQuantity, nil
}
