package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

type InvoiceTransactionRepository interface {
	UseTx(tx *gorm.DB) *gorm.DB
	SaveInvoiceTransaction(invoiceTransaction models.InvoiceTransaction) error
}

type invoiceTransactionRepository struct {
	Db *gorm.DB
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
