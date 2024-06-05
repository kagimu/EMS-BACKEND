package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	invoiceModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
	"time"
)

func SavePayment(invoicePayment invoiceModels.InvoicePayment) (invoiceModels.InvoicePayment, error) {

	if err := config.DB.Save(&invoicePayment).Error; err != nil {
		return invoicePayment, err
	}

	return invoicePayment, nil
}

func GetInvoicePaidAmount(invoiceNumber string, tx ...*gorm.DB) (float64, error) {

	db := config.DB

	if tx != nil {
		db = tx[0]
	}

	var totalAmount float64

	if err := db.Model(&invoiceModels.InvoicePayment{}).
		Where("invoice_number = ? ", invoiceNumber).
		Where("status = ? ", constants.StatusActive).
		Select("COALESCE(sum(amount),0)").
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

func CountAllUserPayments(userNumber string) (int64, error) {

	db := config.DB

	var countPayments int64

	if err := db.Model(&invoiceModels.Invoice{}).
		Joins("JOIN invoice_payments ON invoices.number = invoice_payments.invoice_number").
		Where("invoices.user_number = ?", userNumber).
		Where("invoice_payments.status = ?", constants.StatusActive).
		Count(&countPayments).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return countPayments, err
	}

	return countPayments, nil
}

func SumAllUserPayments(userNumber string) (float64, error) {

	db := config.DB

	var totalAmount float64

	if err := db.Model(&invoiceModels.Invoice{}).
		Joins("JOIN invoice_payments ON invoices.number = invoice_payments.invoice_number").
		Where("invoices.user_number = ? ", userNumber).
		Where("invoice_payments.status = ? ", constants.StatusActive).
		Select("COALESCE(sum(invoice_payments.amount),0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetMostRecentPayment(userNumber string) (invoiceModels.InvoicePayment, error) {

	db := config.DB

	var invoicePayment invoiceModels.InvoicePayment

	if err := db.
		Joins("JOIN invoices ON invoices.number = invoice_payments.invoice_number").
		Where("invoices.user_number = ? ", userNumber).
		Where("invoice_payments.status = ? ", constants.StatusActive).
		Order("invoice_payments.created_at desc").
		First(&invoicePayment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModels.InvoicePayment{}, nil
		}

		return invoiceModels.InvoicePayment{}, err
	}

	return invoicePayment, nil

}

func GetAllPayments(userNumber string, schoolNumber string, refNumber string, startDate time.Time, endDate time.Time) ([]invoiceModels.InvoicePayment, error) {

	db := config.DB

	var invoicePayment []invoiceModels.InvoicePayment

	if err := db.
		Joins("JOIN invoices ON invoices.number = invoice_payments.invoice_number").
		Scopes(
			filters.ByAnyField("invoices.user_number", userNumber),
			filters.BySchoolNumber(schoolNumber),
			filters.ByRefNumber(refNumber),
			filters.ByDateBetween("invoice_payments.date", startDate, endDate)).
		Order("invoice_payments.created_at desc").
		Find(&invoicePayment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoicePayment, nil
		}

		return invoicePayment, err
	}

	return invoicePayment, nil
}

func GetInvoicePaymentByNumber(paymentNumber string) (invoiceModels.InvoicePayment, error) {

	db := config.DB

	var invoicePayment invoiceModels.InvoicePayment

	if err := db.Where("number = ?", paymentNumber).First(&invoicePayment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoicePayment, nil
		}
		return invoicePayment, err
	}

	return invoicePayment, nil
}

func DeleteInvoicePayment(invoicePayment invoiceModels.InvoicePayment) error {

	if err := config.DB.Delete(&invoicePayment).Error; err != nil {
		return err
	}

	return nil

}
