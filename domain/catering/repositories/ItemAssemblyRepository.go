package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetAssembledItemParts(menuItemNumber string, tx ...*gorm.DB) ([]models.ItemAssembly, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

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

func DeleteItemAssembly(itemAssemble models.ItemAssembly, tx ...*gorm.DB) error {
	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&itemAssemble).Error; err != nil {

		if tx != nil {
			db.Rollback()
		}

		return err
	}

	return nil
}
