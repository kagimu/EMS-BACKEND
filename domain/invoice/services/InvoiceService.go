package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	itemRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
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

	invoice, err := invoiceRepositories.FindInvoiceByRefNumber(invoiceNumber)

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

	var student studentModel.Student

	if student, err = repositories.FindStudentByUserNumber(userNumber); err != nil {
		return err
	}

	if student.Number == "" {
		return errors.New("no student details found")
	}

	var admission admissionModel.Admission

	if admission, err = repositories2.FindAdmissionByStudentNumber(student.Number); err != nil {
		return err
	}

	if admission.Number == "" {
		return errors.New("no admission details found")
	}

	newInvoice := models.Invoice{
		Number:        uuid.NewV4().String(),
		DueDate:       dueDate,
		Description:   payload.Description,
		UserNumber:    userNumber,
		RefNumber:     invoiceNumber,
		AddedByNumber: currentUserNumber,
		SchoolNumber:  admission.SchoolNumber,
		CampusNumber:  admission.CampusNumber,
	}

	if err := config.DB.Create(&newInvoice).Error; err != nil {
		return err
	}

	for _, invoiceItem := range payload.Items {

		if err = doSaveInvoiceItem(invoiceItem, newInvoice.Number); err != nil {
			return err
		}

	}

	return nil

}

func doSaveInvoiceItem(invoiceItem requests.InvoiceItemRequest, newInvoiceNumber string) error {

	var item itemModels.Item

	if item, err := itemRepositories.GetItemByNumber(invoiceItem.ItemNumber); err != nil || item.Number == "" {

		if err != nil {
			return err
		}

		return errors.New("invalid item number")
	}

	if checkItemOnInvoice, err := invoiceRepositories.GetInvoiceItemByInvoiceNumberAndItemNumber(newInvoiceNumber, invoiceItem.ItemNumber); err != nil || checkItemOnInvoice.Number != "" {

		if err != nil {
			return err
		}

		fmt.Println("Item '" + invoiceItem.ItemNumber + "' already added to Invoice")

		return errors.New("item already added to Invoice")
	}

	newInvoiceItem := models.InvoiceItem{
		Number:        uuid.NewV4().String(),
		InvoiceNumber: newInvoiceNumber,
		ItemNumber:    invoiceItem.ItemNumber,
		Quantity:      invoiceItem.Quantity,
		UnitPrice:     item.Cost,
	}

	if _, err := invoiceRepositories.CreateInvoiceItem(newInvoiceItem); err != nil {
		return err
	}

	return nil
}

func AllInvoices(status string, userNumber string, campusNumber string, schoolNumber string) ([]responses.MinInvoiceResponse, error) {

	var invoices []models.Invoice

	var invoiceResponses []responses.MinInvoiceResponse

	invoices, err := invoiceRepositories.FindAllInvoices(status, userNumber, campusNumber, schoolNumber)
	if err != nil {
		return invoiceResponses, err
	}

	for _, invoice := range invoices {
		invoiceResponses = append(invoiceResponses, createMinInvoiceResponse(invoice))
	}

	return invoiceResponses, nil
}

func createMinInvoiceResponse(invoice models.Invoice) responses.MinInvoiceResponse {

	var _invoiceAmount float64

	if invoiceAmount, err := invoiceRepositories.GetInvoiceAmount(invoice.Number); err != nil {
		fmt.Println("Getting invoice amount error '" + err.Error() + "' for '" + invoice.RefNumber + "'")
		_invoiceAmount = 0
	} else {
		_invoiceAmount = invoiceAmount
	}

	invoiceResponse := responses.MinInvoiceResponse{
		RefNumber:   invoice.RefNumber,
		Date:        invoice.CreatedAt,
		DueDate:     invoice.DueDate,
		Number:      invoice.Number,
		TotalAmount: _invoiceAmount,
		PaidAmount:  0,
		Balance:     _invoiceAmount,
	}
	return invoiceResponse
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

func InvoiceDetails(refNumber string) (responses.InvoiceResponse, error) {

	invoice, err := invoiceRepositories.FindInvoiceByRefNumber(refNumber)

	if err != nil {
		return responses.InvoiceResponse{}, err
	}

	if invoice.Number == "" {
		return responses.InvoiceResponse{}, errors.New("invalid invoice ref number")
	}

	minInvoice := createMinInvoiceResponse(invoice)

	items, err := invoiceRepositories.GetInvoiceItems(invoice.Number)
	if err != nil {
		return responses.InvoiceResponse{}, err
	}

	var invoiceItemsResponse []responses.InvoiceItemResponse

	for _, item := range items {
		_item := item.Item

		itemResponse := responses.InvoiceItemResponse{
			ItemName:   _item.Name,
			ItemNumber: item.Number,
			UnitCost:   item.UnitPrice,
			Quantity:   item.Quantity,
			TotalCost:  item.UnitPrice * item.Quantity,
		}

		invoiceItemsResponse = append(invoiceItemsResponse, itemResponse)

	}

	invoiceResponses := responses.InvoiceResponse{
		minInvoice,
		invoiceItemsResponse,
	}

	return invoiceResponses, nil
}
