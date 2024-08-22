package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SaveBillableItem(billable models.BillableItem, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&billable).Error
}

func FindBillableItems(billableNumber string, trn ...*gorm.DB) ([]models.BillableItem, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billableItems []models.BillableItem

	if err := db.
		Where("billable_number = ?", billableNumber).Find(&billableItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billableItems, nil
		}

		return nil, err
	}

	return billableItems, nil
}

func RemoveBillableItem(billableItem models.BillableItem, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&billableItem).Error
}

func GetDuplicatedBillableItem(billableNumber, itemNumber string, trn ...*gorm.DB) (models.BillableItem, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var billableItem models.BillableItem

	if err := db.
		Where("item_number = ? AND billable_number = ?", itemNumber, billableNumber).First(&billableItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return billableItem, nil
		}

		return billableItem, err
	}

	return billableItem, nil
}
