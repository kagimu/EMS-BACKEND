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

func SaveItemAssembly(assembly models.ItemAssembly) error {

	db := config.DB

	if err := db.Save(&assembly).Error; err != nil {
		return err
	}

	return nil
}

func FindAssemblyItemByPrimaryNumberAndSecondaryNumber(primaryItemNumber string, secondaryItemNumber string) (models.ItemAssembly, error) {

	db := config.DB

	var item models.ItemAssembly

	if err := db.
		Where("primary_item_number = ?", primaryItemNumber).
		Where("secondary_item_number = ?", secondaryItemNumber).
		First(&item).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, nil
		}

		return item, err
	}

	return item, nil
}
