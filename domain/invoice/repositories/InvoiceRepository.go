package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

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
