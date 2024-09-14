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

func SaveTransaction(invoicePayment invoiceModels.Transaction, tx ...*gorm.DB) (invoiceModels.Transaction, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&invoicePayment).Error; err != nil {
		return invoicePayment, err
	}

	return invoicePayment, nil
}

func GetInvoicePaidAmount(invoiceNumber string, tx ...*gorm.DB) (float64, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var totalAmount float64

	if err := db.Model(&invoiceModels.InvoiceTransaction{}).
		Where("invoice_number = ? ", invoiceNumber).
		Select("COALESCE(sum(amount_paid_off),0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetInvoicePayments(invoiceNumber string) ([]invoiceModels.Transaction, error) {

	db := config.DB

	var invoiceItems []invoiceModels.Transaction

	if err := db.
		Joins("").
		Preload("User").Preload("Invoice").Where("invoice_number = ?", invoiceNumber).Find(&invoiceItems).Error; err != nil {

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

	if err := db.Model(&invoiceModels.Statement{}).
		Joins("join transactions on transactions.statement_number = statements.number").
		Where("statements.user_number = ?", userNumber).
		Where("transactions.status = ?", constants.StatusActive).
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

	if err := db.Model(&invoiceModels.Statement{}).
		Joins("JOIN transactions ON transactions.statement_number = statements.number").
		Where("statement.user_number = ? ", userNumber).
		Where("transactions.status = ? ", constants.StatusActive).
		Select("COALESCE(sum(statement.amount),0)").
		Scan(&totalAmount).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalAmount, err
	}

	return totalAmount, nil
}

func GetMostRecentPayment(userNumber string) (invoiceModels.Transaction, error) {

	db := config.DB

	var invoicePayment invoiceModels.Transaction

	if err := db.
		Joins("JOIN statements ON statements.number = transactions.statement_number").
		Where("statements.user_number = ? ", userNumber).
		Where("transactions.status = ? ", constants.StatusActive).
		Order("statements.created_at desc").
		First(&invoicePayment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModels.Transaction{}, nil
		}

		return invoiceModels.Transaction{}, err
	}

	return invoicePayment, nil

}

func GetAllPayments(userNumber string, schoolNumber string, refNumber string, startDate time.Time, endDate time.Time) ([]invoiceModels.Transaction, error) {

	db := config.DB

	var invoicePayment []invoiceModels.Transaction

	if err := db.
		Joins("JOIN statements ON statements.number = transactions.statement_number").
		Scopes(
			filters.ByAnyField("statements.user_number", userNumber),
			filters.BySchoolNumber(schoolNumber),
			filters.ByRefNumber(refNumber),
			filters.ByDateBetween("transactions.date", startDate, endDate)).
		Order("statements.created_at desc").
		Find(&invoicePayment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoicePayment, nil
		}

		return invoicePayment, err
	}

	return invoicePayment, nil
}

func GetTransactionByNumber(transactionNumber string) (invoiceModels.Transaction, error) {

	db := config.DB

	var invoicePayment invoiceModels.Transaction

	if err := db.Where("number = ?", transactionNumber).First(&invoicePayment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoicePayment, nil
		}
		return invoicePayment, err
	}

	return invoicePayment, nil
}

func DeleteTransaction(transactionNumber invoiceModels.Transaction, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&transactionNumber).Error; err != nil {
		return err
	}

	return nil
}
