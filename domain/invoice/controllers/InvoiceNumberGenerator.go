package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
)

func InvoiceNumberGenerator() (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	saveInvoiceNumber := models.InvoiceNumber{
		Year:  year,
		Month: month,
	}

	if err := config.DB.Create(&saveInvoiceNumber).Error; err != nil {
		return "", err
	}

	savedID := saveInvoiceNumber.ID
	result := fmt.Sprintf("CST%s%05d", saveInvoiceNumber.Year, savedID)

	return result, nil
}
