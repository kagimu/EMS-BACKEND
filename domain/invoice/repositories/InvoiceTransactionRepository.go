package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

type InvoiceTransactionRepository interface {
	UseTx(tx *gorm.DB) *gorm.DB
	SaveInvoiceTransaction(invoiceTransaction models.InvoiceTransaction) error
	FindInvoiceTransactions(invoiceNumber, transactionNumber string) ([]models.InvoiceTransaction, error)
	DeleteTransaction(invoiceTransaction models.InvoiceTransaction) error
}

type invoiceTransactionRepository struct {
	Db *gorm.DB
}

func (invoiceTransactionRepository *invoiceTransactionRepository) DeleteTransaction(invoiceTransaction models.InvoiceTransaction) error {
	return invoiceTransactionRepository.Db.Delete(&invoiceTransaction).Error
}

func (invoiceTransactionRepository *invoiceTransactionRepository) FindInvoiceTransactions(invoiceNumber, transactionNumber string) ([]models.InvoiceTransaction, error) {

	db := invoiceTransactionRepository.Db

	var transactions []models.InvoiceTransaction

	if err := db.
		Scopes(filters.ByAnyField("invoice_number", invoiceNumber), filters.ByAnyField("transaction_number", transactionNumber)).
		Find(&transactions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]models.InvoiceTransaction, 0), nil
		}

		return transactions, err
	}

	return transactions, nil
}

func (invoiceTransactionRepository *invoiceTransactionRepository) SaveInvoiceTransaction(invoiceTransaction models.InvoiceTransaction) error {
	return invoiceTransactionRepository.Db.Save(&invoiceTransaction).Error
}

func (invoiceTransactionRepository *invoiceTransactionRepository) UseTx(tx *gorm.DB) *gorm.DB {
	invoiceTransactionRepository.Db = tx
	return invoiceTransactionRepository.Db
}

func NewInvoiceTransactionRepository() InvoiceTransactionRepository {
	return &invoiceTransactionRepository{
		Db: config.DB,
	}
}
