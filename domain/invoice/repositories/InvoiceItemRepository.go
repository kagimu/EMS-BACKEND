package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func GetInvoiceItemByInvoiceNumberAndItemNumber(invoiceNumber string, itemNumber string, txn ...*gorm.DB) (invoiceModel.InvoiceItem, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var invoiceItem invoiceModel.InvoiceItem

	if err := db.Where("invoice_number = ? ", invoiceNumber).Where("item_number = ? ", itemNumber).First(&invoiceItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceItem, nil
		}

		return invoiceItem, err
	}

	return invoiceItem, nil

}

func SaveInvoiceItem(invoiceItem invoiceModel.InvoiceItem, tx ...*gorm.DB) (invoiceModel.InvoiceItem, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&invoiceItem).Error; err != nil {
		return invoiceModel.InvoiceItem{}, err
	}

	return invoiceItem, nil

}

func GetInvoiceAmount(invoiceNumber string) (float64, error) {

	db := config.DB

	var totalAmount float64

	if err := db.Model(&invoiceModel.InvoiceItem{}).
		Where("invoice_number = ? ", invoiceNumber).
		Select("COALESCE(sum(quantity * unit_price),0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetInvoiceItems(invoiceNumber string, tx ...*gorm.DB) ([]invoiceModel.InvoiceItem, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var invoiceItems []invoiceModel.InvoiceItem

	if err := db.Preload("Item").Where("invoice_number = ?", invoiceNumber).Find(&invoiceItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceItems, nil
		}

		return invoiceItems, err
	}

	return invoiceItems, nil
}

func DeleteInvoiceItem(item invoiceModel.InvoiceItem, tx *gorm.DB) error {

	if err := tx.Delete(&item).Error; err != nil {
		return err
	}
	return nil
}
