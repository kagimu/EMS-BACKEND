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

	if err := db.Omit("Invoice", "Item").Save(&invoiceItem).Error; err != nil {
		return invoiceModel.InvoiceItem{}, err
	}

	return invoiceItem, nil

}

func GetInvoiceAmount(invoiceNumber string, tx ...*gorm.DB) (float64, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

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

func GetInvoiceItemPaidAmount(invoiceItemNumber string, tx ...*gorm.DB) (float64, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var paidAmount float64

	if err := db.Model(&invoiceModel.InvoiceItemTransaction{}).
		Where("invoice_item_number = ? ", invoiceItemNumber).
		Select("COALESCE(sum(amount_paid_off),0)").
		Scan(&paidAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return paidAmount, err
	}

	return paidAmount, nil
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

func GetInvoiceItemTransactionsHistory(invoiceNumber string, tx ...*gorm.DB) ([]invoiceModel.InvoiceItemTransaction, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var history []invoiceModel.InvoiceItemTransaction
	err := db.Preload("Transaction.Statement.AddedBy").
		Preload("InvoiceItem.Item").
		Joins("JOIN invoice_items on invoice_items.number = invoice_item_transactions.invoice_item_number").
		Where("invoice_items.invoice_number = ?", invoiceNumber).
		Find(&history).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return history, err
	}

	return history, nil
}

func GetAllItemTransactionsHistory(campusNumber string, tx ...*gorm.DB) ([]invoiceModel.InvoiceItemTransaction, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var history []invoiceModel.InvoiceItemTransaction
	query := db.Preload("Transaction.Statement.AddedBy").
		Preload("InvoiceItem.Item")

	if campusNumber != "" {
		query = query.Joins("JOIN transactions on transactions.number = invoice_item_transactions.transaction_number").
			Joins("JOIN statements on statements.number = transactions.statement_number").
			Where("statements.campus_number = ?", campusNumber)
	}

	err := query.Find(&history).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return history, err
	}

	return history, nil
}
