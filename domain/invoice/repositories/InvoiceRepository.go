package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
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

func FindInvoiceByNumber(invoiceNumber string) (invoiceModel.Invoice, error) {

	db := config.DB
	var invoice invoiceModel.Invoice

	if err := db.Where("number = ?", invoiceNumber).First(&invoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoice, nil
		}

		return invoice, err
	}

	return invoice, nil
}

func FindAllInvoices(status string, userNumber string, campusNumber string, schoolNumber string) ([]invoiceModel.Invoice, error) {

	db := config.DB

	var invoices []invoiceModel.Invoice

	if err := db.Scopes(
		filters.ByStatus(status),
		filters.ByUserNumber(userNumber),
		filters.BySchoolNumber(schoolNumber),
		filters.ByCampusNumber(campusNumber)).
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

	if err := db.Model(&invoiceModel.Invoice{}).
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

	if err := db.Model(&invoiceModel.Invoice{}).
		Joins("JOIN invoice_items ON invoices.number = invoice_items.invoice_number").
		Where("invoices.user_number = ? ", userNumber).
		Select("sum(invoice_items.quantity * invoice_items.unit_price)").
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
		Where("user_number = ? ", userNumber).
		Order("created_at desc").
		First(&invoice).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoice, nil
		}

		return invoice, err
	}

	return invoice, nil

}

func SaveInvoice(invoice invoiceModel.Invoice) error {

	if err := config.DB.Create(&invoice).Error; err != nil {
		return err
	}

	return nil
}

func UpdateInvoice(invoice invoiceModel.Invoice) error {

	if err := config.DB.Save(&invoice).Error; err != nil {
		return err
	}

	return nil
}
