package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func GetInvoiceItemByInvoiceNumberAndItemNumber(invoiceNumber string, itemNumber string) (invoiceModel.InvoiceItem, error) {

	db := config.DB

	var invoiceItem invoiceModel.InvoiceItem

	if err := db.Where("invoice_number = ? ", invoiceNumber).Where("item_number = ? ", itemNumber).Find(&invoiceItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceItem, nil
		}

		return invoiceItem, err
	}

	return invoiceItem, nil

}

func CreateInvoiceItem(invoiceItem invoiceModel.InvoiceItem) (invoiceModel.InvoiceItem, error) {

	db := config.DB

	if err := db.Create(&invoiceItem).Error; err != nil {
		return invoiceModel.InvoiceItem{}, err
	}

	return invoiceItem, nil

}
