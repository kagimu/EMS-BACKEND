package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	itemRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

func CreateNewInvoice(payload requests.InvoiceRequest, currentUserNumber string) error {

	userNumber := payload.UserNumber

	invoiceUser, err := repository.GetUserByNumber(userNumber)

	if err != nil {
		return err
	}

	if invoiceUser.Number == "" {
		return errors.New("invalid user number")
	}

	invoiceNumber := generateNewInvoiceNumber()

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)

	if err != nil {
		return err
	}

	if invoice.Number != "" {
		return errors.New("try again to generate a new invoice number")
	}

	var dueDate time.Time
	if dueDate, err = time.Parse(time.DateOnly, payload.DueDate); err != nil {
		return errors.New("Due date error: " + payload.DueDate)
	}

	newInvoice := models.Invoice{
		Number:        uuid.NewV4().String(),
		DueDate:       dueDate,
		Description:   payload.Description,
		UserNumber:    userNumber,
		RefNumber:     invoiceNumber,
		AddedByNumber: currentUserNumber,
	}

	if err := config.DB.Create(&newInvoice).Error; err != nil {
		return err
	}

	for _, invoiceItem := range payload.Items {

		var item itemModels.Item

		if item, err = itemRepositories.GetItemByNumber(invoiceItem.ItemNumber); err != nil || item.Number == "" {

			if err != nil {
				return err
			}

			return errors.New("invalid item number")
		}

		if checkItemOnInvoice, err := invoiceRepositories.GetInvoiceItemByInvoiceNumberAndItemNumber(newInvoice.Number, invoiceItem.ItemNumber); err != nil || checkItemOnInvoice.Number != "" {

			if err != nil {
				return err
			}

			fmt.Println("Item '" + invoiceItem.ItemNumber + "' already added to Invoice")

			return errors.New("item already added to Invoice")
		}

		newInvoiceItem := models.InvoiceItem{
			Number:        uuid.NewV4().String(),
			InvoiceNumber: newInvoice.Number,
			ItemNumber:    invoiceItem.ItemNumber,
			Quantity:      invoiceItem.Quantity,
			UnitPrice:     item.Cost,
		}

		if _, err := invoiceRepositories.CreateInvoiceItem(newInvoiceItem); err != nil {
			return err
		}

	}

	return nil

}

func generateNewInvoiceNumber() string {
	now := time.Now()
	year := now.Format("06")  // YY format
	month := now.Format("01") // MM format

	// Generate random four-digit number
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(10000) // Generate random number between 0 and 9999

	// Format invoice number
	return fmt.Sprintf("%s%s%04d", year, month, randomNum)
}
