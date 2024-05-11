package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SavePayment(invoicePayment invoiceModels.InvoicePayment) (invoiceModels.InvoicePayment, error) {

	if err := config.DB.Create(&invoicePayment).Error; err != nil {
		return invoicePayment, err
	}

	return invoicePayment, nil
}

func GetInvoicePaidAmount(invoiceNumber string) (float64, error) {

	db := config.DB

	var totalAmount float64

	if err := db.Model(&invoiceModels.InvoicePayment{}).
		Where("invoice_number = ? ", invoiceNumber).
		Select("sum(amount)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetInvoicePayments(invoiceNumber string) ([]invoiceModels.InvoicePayment, error) {

	db := config.DB

	var invoiceItems []invoiceModels.InvoicePayment

	if err := db.Preload("User").Preload("Invoice").Where("invoice_number = ?", invoiceNumber).Find(&invoiceItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceItems, nil
		}

		return invoiceItems, err
	}

	return invoiceItems, nil
}
