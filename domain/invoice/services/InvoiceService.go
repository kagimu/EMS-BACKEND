package services

import (
	"errors"
	"fmt"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
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

	if err = invoiceRepositories.SaveInvoice(newInvoice); err != nil {
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

	item, err := itemRepositories.GetItemByNumber(invoiceItem.ItemNumber)

	if err != nil {
		return err
	}

	if item.Number == "" {
		return errors.New("invalid item number")
	}

	if checkItemOnInvoice, err := invoiceRepositories.GetInvoiceItemByInvoiceNumberAndItemNumber(newInvoiceNumber, invoiceItem.ItemNumber); err != nil || checkItemOnInvoice.Number != "" {

		if err != nil {
			return err
		}

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
	var _paidAmount float64

	if invoiceAmount, err := invoiceRepositories.GetInvoiceAmount(invoice.Number); err != nil {
		fmt.Println("Getting invoice amount error '" + err.Error() + "' for '" + invoice.RefNumber + "'")
		_invoiceAmount = 0
	} else {
		_invoiceAmount = invoiceAmount
	}

	if paidAmount, err := invoiceRepositories.GetInvoicePaidAmount(invoice.Number); err != nil {
		fmt.Println("Getting invoice paid amount error '" + err.Error() + "' for '" + invoice.RefNumber + "'")
		_paidAmount = 0
	} else {
		_paidAmount = paidAmount
	}

	invoiceResponse := responses.MinInvoiceResponse{
		RefNumber:   invoice.RefNumber,
		Date:        invoice.CreatedAt.Format(time.DateOnly),
		DueDate:     invoice.DueDate.Format(time.DateOnly),
		Number:      invoice.Number,
		TotalAmount: _invoiceAmount,
		PaidAmount:  _paidAmount,
		Balance:     _invoiceAmount - _paidAmount,
		Status:      invoice.Status,
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

	payments, err := invoiceRepositories.GetInvoicePayments(invoice.Number)
	if err != nil {
		return responses.InvoiceResponse{}, err
	}

	var invoicePaymentsResponse []responses.InvoicePaymentResponse

	for _, payment := range payments {
		_invoice := payment.Invoice
		_user := payment.User

		paymentResponse := createPaymentResponse(payment, _user, _invoice)

		invoicePaymentsResponse = append(invoicePaymentsResponse, paymentResponse)

	}

	invoiceResponse := responses.InvoiceResponse{
		MinInvoiceResponse: minInvoice,
		Items:              invoiceItemsResponse,
		Payments:           invoicePaymentsResponse,
	}

	return invoiceResponse, nil
}

func createPaymentResponse(payment models.InvoicePayment, user authModel.User, invoice models.Invoice) responses.InvoicePaymentResponse {
	paymentResponse := responses.InvoicePaymentResponse{
		Amount:           payment.Amount,
		InvoiceRefNumber: invoice.RefNumber,
		RecordedBy:       user.FirstName,
		DateTime:         payment.Date.Format(time.DateOnly),
	}
	return paymentResponse
}

func ReceiveInvoicePayment(invoicePaymentRequest requests.InvoicePaymentRequest, currentUserNumber string) (responses.MinInvoiceResponse, error) {

	var minInvoiceResponse responses.MinInvoiceResponse

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoicePaymentRequest.InvoiceNumber)
	if err != nil {
		return minInvoiceResponse, err
	}

	if invoice.Number == "" {
		return minInvoiceResponse, errors.New("invalid invoice number: " + invoicePaymentRequest.InvoiceNumber)
	}

	var dateTimePaid time.Time
	if dateTimePaid, err = time.Parse(time.DateOnly, invoicePaymentRequest.Date); err != nil {
		return minInvoiceResponse, errors.New("Date error: " + err.Error() + " for " + invoicePaymentRequest.Date)
	}

	saveInvoicePayment := models.InvoicePayment{
		Number:        uuid.NewV4().String(),
		Amount:        invoicePaymentRequest.Amount,
		Date:          dateTimePaid,
		InvoiceNumber: invoice.Number,
		UserNumber:    currentUserNumber,
	}

	if _, err := invoiceRepositories.SavePayment(saveInvoicePayment); err != nil {
		return minInvoiceResponse, err
	}

	minInvoiceResponse = createMinInvoiceResponse(invoice)

	var newStatus string

	if minInvoiceResponse.PaidAmount >= minInvoiceResponse.TotalAmount {
		newStatus = constants.StatusPaid
	} else {
		newStatus = constants.StatusPartiallyPaid
	}

	if newStatus != invoice.Status {
		invoice.Status = newStatus
		if err = invoiceRepositories.UpdateInvoice(invoice); err != nil {
			return minInvoiceResponse, errors.New("Updating invoice failed : " + err.Error())
		}
	}

	return minInvoiceResponse, nil
}

func GetFinancialSummary(userNumber string) (responses.FinancialSummaryResponse, error) {

	paymentsCount, err := invoiceRepositories.CountAllUserPayments(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}

	invoicesCount, err := invoiceRepositories.CountAllUserInvoices(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}

	invoiceAmount, err := invoiceRepositories.GetAllInvoicedAmount(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}

	paymentAmount, err := invoiceRepositories.SumAllUserPayments(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}

	recentInvoice, err := invoiceRepositories.GetMostRecentInvoice(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}
	minInvoice := createMinInvoiceResponse(recentInvoice)

	mostRecentPayment, err := invoiceRepositories.GetMostRecentPayment(userNumber)
	if err != nil {
		return responses.FinancialSummaryResponse{}, err
	}

	mostRecent := responses.MostRecent{
		InvoiceAmount: minInvoice.TotalAmount,
		InvoiceDate:   minInvoice.Date,
		PaidAmount:    mostRecentPayment.Amount,
		PaidDate:      mostRecentPayment.Date.Format(time.DateOnly),
	}

	var financialSummaryResponse = responses.FinancialSummaryResponse{
		TotalPayments:      paymentsCount,
		TotalInvoices:      invoicesCount,
		TotalInvoiceAmount: invoiceAmount,
		TotalPaidAmount:    paymentAmount,
		TotalBalanceAmount: invoiceAmount - paymentAmount,
		MostRecent:         mostRecent,
	}

	return financialSummaryResponse, nil
}

func GetAllPayments(userNumber string, refNumber string, schoolNumber string) ([]responses.InvoicePaymentResponse, error) {

	var invoicePaymentsResponse []responses.InvoicePaymentResponse

	allPayments, err := invoiceRepositories.GetAllPayments(userNumber, schoolNumber, refNumber)

	if err != nil {
		return invoicePaymentsResponse, err
	}

	for _, payment := range allPayments {

		_user, err := repository.GetUserByNumber(payment.UserNumber)
		if err != nil {
			return invoicePaymentsResponse, err
		}

		_invoice, err := invoiceRepositories.FindInvoiceByNumber(payment.InvoiceNumber)
		if err != nil {
			return invoicePaymentsResponse, err
		}

		paymentResponse := createPaymentResponse(payment, _user, _invoice)

		invoicePaymentsResponse = append(invoicePaymentsResponse, paymentResponse)

	}

	return invoicePaymentsResponse, nil
}
