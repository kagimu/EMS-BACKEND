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
