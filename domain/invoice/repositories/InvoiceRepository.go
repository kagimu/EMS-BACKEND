package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
	"time"
)

func FindInvoiceByRefNumber(refNumber string) (invoiceModel.Invoice, error) {

	db := config.DB
	var invoice invoiceModel.Invoice

	if err := db.Where("ref_number = ?", refNumber).First(&invoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoice, nil
		}

		return invoice, err
	}

	return invoice, nil
}

func FindInvoiceByNumber(invoiceNumber string, tx ...*gorm.DB) (invoiceModel.Invoice, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var invoice invoiceModel.Invoice

	if err := db.
		Preload("Statement").
		Where("number = ?", invoiceNumber).First(&invoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoice, nil
		}

		return invoice, err
	}

	return invoice, nil
}

func FindAllInvoices(status string, userNumber string, campusNumber string, schoolNumber string, invoiceItems []string, startDate time.Time, endDate time.Time) ([]invoiceModel.Invoice, error) {

	db := config.DB

	var invoices []invoiceModel.Invoice

	if err := db.
		Joins("JOIN statements ON statements.number =  invoices.statement_number").
		Joins("JOIN invoice_items ON invoice_items.invoice_number =  invoices.number").
		Scopes(
			filters.ByStatus(status),
			filters.ByUserNumber(userNumber),
			filters.BySchoolNumber(schoolNumber),
			filters.ByValueIn("invoice_items.item_number", invoiceItems),
			filters.ByDateBetween("created_at", startDate, endDate),
			filters.ByCampusNumber(campusNumber)).
		Group("invoices.number").
		Find(&invoices).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoices, nil
		}

		return invoices, err
	}
	return invoices, nil
}

func CountAllUserInvoices(userNumber string) (int64, error) {

	db := config.DB

	var invoicesCount int64

	if err := db.
		Joins("JOIN statements ON statements.number =  invoices.statement_number").
		Model(&invoiceModel.Invoice{}).
		Where("user_number = ?", userNumber).
		Count(&invoicesCount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return invoicesCount, err
	}

	return invoicesCount, nil
}

func GetAllInvoicedAmount(userNumber string) (float64, error) {

	db := config.DB

	var totalAmount float64

	if err := db.
		Model(&invoiceModel.Invoice{}).
		Joins("JOIN statements ON statements.number =  invoices.statement_number").
		Joins("JOIN invoice_items ON invoice_items.invoice_number = invoices.number").
		Where("statements.user_number = ? ", userNumber).
		Select("COALESCE(sum(invoice_items.quantity * invoice_items.unit_price), 0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetMostRecentInvoice(userNumber string) (invoiceModel.Invoice, error) {

	db := config.DB
	var invoice invoiceModel.Invoice

	if err := db.
		Joins("JOIN statements ON statements.number =  invoices.statement_number").
		Where("statements.user_number = ? ", userNumber).
		Order("statements.created_at desc").
		First(&invoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModel.Invoice{}, nil
		}

		return invoiceModel.Invoice{}, err
	}

	return invoice, nil

}

func SaveInvoice(invoice invoiceModel.Invoice, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	if err := db.Save(&invoice).Error; err != nil {
		return err
	}

	return nil
}

func UpdateInvoice(invoice invoiceModel.Invoice, tx ...*gorm.DB) error {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&invoice).Error; err != nil {
		return err
	}

	return nil
}

func DeleteInvoice(invoice invoiceModel.Invoice, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&invoice).Error; err != nil {
		return err
	}

	return nil
}
