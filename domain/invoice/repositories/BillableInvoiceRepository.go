package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SaveBillableInvoice(billableInvoice models.BillableInvoice, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&billableInvoice).Error
}

func GetBillableInvoiceByBillableNumber(billableNumber string, txn ...*gorm.DB) (models.BillableInvoice, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var billableInvoice models.BillableInvoice

	if err := db.Where("billable_number = ?", billableNumber).First(&billableInvoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BillableInvoice{}, nil
		}

		return models.BillableInvoice{}, err

	}

	return billableInvoice, nil
}
