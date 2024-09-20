package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	itemRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

var invoiceTransactionRepository = invoiceRepositories.NewInvoiceTransactionRepository()

var admissionRepository = repositories2.NewAdmissionRepository()

func CreateInvoice(payload requests.InvoiceRequest, currentUserNumber string) error {

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

	dueDate, err := commons.ConvertStringToDate(payload.DueDate)

	var student studentModel.Student

	if student, err = repositories.FindStudentByUserNumber(userNumber); err != nil {
		return err
	}

	if student.Number == "" {
		return errors.New("no student details found")
	}

	var admission admissionModel.Admission

	if admission, err = admissionRepository.FindAdmissionByStudentNumber(student.Number); err != nil {
		return err
	}

	if admission.Number == "" {
		return errors.New("no admission details found")
	}

	invoiceTotal := newInvoiceTotalAmount(payload.Items)

	balance, err := ComputeNextBalance(invoiceTotal, constants.StatementTypeOut, userNumber)

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Amount:        invoiceTotal,
		UserNumber:    payload.UserNumber,
		Description:   "New bill " + invoiceNumber,
		SchoolNumber:  admission.SchoolNumber,
		CampusNumber:  admission.CampusNumber,
		RefNumber:     invoiceNumber,
		Balance:       balance,
		Type:          constants.StatementTypeOut,
		AddedByNumber: currentUserNumber,
	}

	newInvoice := models.Invoice{
		StatementNumber: newStatement.Number,
		Number:          uuid.NewV4().String(),
		DueDate:         dueDate,
		Description:     payload.Description,
		RefNumber:       invoiceNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Error creating new bill statement " + err.Error())
		}

		err = invoiceRepositories.SaveInvoice(newInvoice, tx)
		if err != nil {
			return errors.New("Error creating new bill " + err.Error())
		}

		for _, invoiceItem := range payload.Items {

			if err = doSaveInvoiceItem(invoiceItem, newInvoice.Number, tx); err != nil {
				return err
			}

		}

		return nil
	})

	return err
}

func CreateInvoiceFromBillable(billable models.Billable, student studentModel.Student, clientReference, description, currentUserNumber string, paidAmount float64, datePaid string, txn *gorm.DB) (models.Invoice, error) {

	invoiceNumber := generateNewInvoiceNumber()

	invoice, err := invoiceRepositories.FindInvoiceByRefNumber(invoiceNumber)

	if err != nil {
		return models.Invoice{}, err
	}

	if invoice.Number != "" {
		return models.Invoice{}, errors.New("try again to generate a new invoice number")
	}

	dueDate := time.Now().AddDate(0, 0, 14)

	billableItems, err := invoiceRepositories.FindBillableItems(billable.Number, txn)
	if err != nil {
		return models.Invoice{}, errors.New("Getting billable items error: " + err.Error())
	}

	var invoiceTotal float64

	for _, billItem := range billableItems {
		invoiceTotal = invoiceTotal + (billItem.Quantity * billItem.UnitPrice)
	}

	balance, err := ComputeNextBalance(invoiceTotal, constants.StatementTypeOut, student.UserNumber)

	invoiceStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Amount:        invoiceTotal,
		UserNumber:    student.UserNumber,
		Description:   "New bill " + invoiceNumber,
		SchoolNumber:  billable.SchoolNumber,
		CampusNumber:  billable.CampusNumber,
		RefNumber:     invoiceNumber,
		Balance:       balance,
		Type:          constants.StatementTypeOut,
		AddedByNumber: currentUserNumber,
	}

	newInvoice := models.Invoice{
		StatementNumber: invoiceStatement.Number,
		Number:          uuid.NewV4().String(),
		DueDate:         dueDate,
		Description:     billable.Description,
		RefNumber:       invoiceNumber,
	}

	if err = invoiceRepositories.SaveStatement(invoiceStatement, txn); err != nil {
		return models.Invoice{}, errors.New("Creating statement error: " + err.Error())
	}

	if err = invoiceRepositories.SaveInvoice(newInvoice, txn); err != nil {
		return models.Invoice{}, errors.New("Creating invoice error: " + err.Error())
	}

	err = invoiceRepositories.SaveBillableInvoice(models.BillableInvoice{
		BillableNumber: billable.Number,
		InvoiceNumber:  newInvoice.Number,
	}, txn)

	if err != nil {
		return models.Invoice{}, errors.New("Creating attaching invoice to billable error: " + err.Error())
	}

	for _, invoiceItem := range billableItems {

		if err = doSaveBillableInvoiceItem(invoiceItem, newInvoice.Number, txn); err != nil {
			return models.Invoice{}, err
		}

	}

	if paidAmount > 0 {
		err = addTransactionToInvoice(invoiceStatement, paidAmount, newInvoice, currentUserNumber, clientReference, description, datePaid, txn)
	}

	return newInvoice, err
}

func addTransactionToInvoice(invoiceStatement models.Statement, paidAmount float64, invoice models.Invoice, currentUserNumber, clientReference, description, datePaid string, txn *gorm.DB) error {

	nextBlc, err := ComputeNextBalance(paidAmount, constants.StatementTypeIn, invoiceStatement.UserNumber)
	if err != nil {
		return err
	}

	transactionStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Amount:        paidAmount,
		UserNumber:    invoiceStatement.UserNumber,
		Description:   "Payment received on invoice " + invoiceStatement.RefNumber,
		SchoolNumber:  invoiceStatement.SchoolNumber,
		CampusNumber:  invoiceStatement.CampusNumber,
		RefNumber:     generateNewInvoiceNumber(),
		Type:          constants.StatementTypeIn,
		Balance:       nextBlc,
		AddedByNumber: currentUserNumber,
	}

	invoiceBalance := invoiceStatement.Amount

	payAmount := invoiceBalance

	if paidAmount < invoiceBalance {
		payAmount = paidAmount
	}

	var transactionDate time.Time
	if transactionDate, err = time.Parse(time.DateOnly, datePaid); err != nil {
		return errors.New("Error reading payment date: " + err.Error())
	}

	newTransaction := models.Transaction{
		StatementNumber: transactionStatement.Number,
		Number:          uuid.NewV4().String(),
		Date:            transactionDate,
		Status:          constants.StatusActive,
		AmountInCredit:  paidAmount - payAmount,
		Description:     description,
		ClientReference: clientReference,
	}

	newInvoiceTransaction := models.InvoiceTransaction{
		InvoiceNumber:     invoice.Number,
		TransactionNumber: newTransaction.Number,
		AmountPaidOff:     payAmount,
		AddedByNumber:     currentUserNumber,
	}

	invoiceTransactionRepository.UseTx(txn)

	err = invoiceRepositories.SaveStatement(transactionStatement, txn)
	if err != nil {
		return errors.New("Creating statement error " + err.Error())
	}

	if _, err = invoiceRepositories.SaveTransaction(newTransaction, txn); err != nil {
		return errors.New("Adding new transaction error " + err.Error())
	}

	err = invoiceTransactionRepository.SaveInvoiceTransaction(newInvoiceTransaction)
	if err != nil {
		return errors.New("Attaching transaction to invoice error " + err.Error())
	}

	if payAmount >= invoiceBalance {
		invoice.Status = constants.StatusPaid
	} else {
		invoice.Status = constants.StatusPartiallyPaid
	}

	if err := invoiceRepositories.UpdateInvoice(invoice, txn); err != nil {
		return errors.New("Updating invoice status failed : " + err.Error())
	}

	return err
}

func doSaveInvoiceItem(invoiceItem requests.InvoiceItemRequest, newInvoiceNumber string, tx *gorm.DB) error {

	item, err := itemRepositories.GetItemByNumber(invoiceItem.ItemNumber, tx)

	if err != nil {
		return err
	}

	if item.Number == "" {
		return errors.New("invalid item number")
	}

	if checkItemOnInvoice, err := invoiceRepositories.GetInvoiceItemByInvoiceNumberAndItemNumber(newInvoiceNumber, invoiceItem.ItemNumber, tx); err != nil || checkItemOnInvoice.Number != "" {

		if err != nil {
			return err
		}

		return errors.New("item already added to Invoice")
	}

	newInvoiceItem := models.InvoiceItem{
		Number:          uuid.NewV4().String(),
		InvoiceNumber:   newInvoiceNumber,
		ItemNumber:      invoiceItem.ItemNumber,
		Quantity:        invoiceItem.Quantity,
		UnitPrice:       invoiceItem.UnitCost,
		ActualUnitPrice: item.Cost,
	}

	if _, err := invoiceRepositories.SaveInvoiceItem(newInvoiceItem, tx); err != nil {
		return err
	}

	return nil
}

func newInvoiceTotalAmount(invoiceItem []requests.InvoiceItemRequest) float64 {
	var total float64

	for _, item := range invoiceItem {
		total = total + item.Quantity*item.UnitCost
	}

	return total
}

func doSaveBillableInvoiceItem(invoiceItem models.BillableItem, newInvoiceNumber string, txn *gorm.DB) error {

	item, err := itemRepositories.GetItemByNumber(invoiceItem.ItemNumber, txn)

	if err != nil {
		return errors.New("Checking item error: " + err.Error())
	}

	if item.Number == "" {
		return errors.New("invalid item number")
	}

	if checkItemOnInvoice, err := invoiceRepositories.GetInvoiceItemByInvoiceNumberAndItemNumber(newInvoiceNumber, invoiceItem.ItemNumber, txn); err != nil || checkItemOnInvoice.Number != "" {

		if err != nil {
			return err
		}

		return errors.New("item already added to Invoice")
	}

	newInvoiceItem := models.InvoiceItem{
		Number:          uuid.NewV4().String(),
		InvoiceNumber:   newInvoiceNumber,
		ItemNumber:      invoiceItem.ItemNumber,
		Quantity:        invoiceItem.Quantity,
		UnitPrice:       invoiceItem.UnitPrice,
		ActualUnitPrice: item.Cost,
	}

	if _, err := invoiceRepositories.SaveInvoiceItem(newInvoiceItem, txn); err != nil {
		return err
	}

	return nil
}

func AllInvoices(status string, userNumber string, campusNumber string, schoolNumber string, commaSeperatedItems string, startDate time.Time, endDate time.Time) ([]responses.MinInvoiceResponse, error) {

	var invoiceResponses = make([]responses.MinInvoiceResponse, 0)

	var invoiceItems = make([]string, 0)

	if commaSeperatedItems != "" {
		invoiceItems = strings.Split(commaSeperatedItems, ",")
	}

	invoices, err := invoiceRepositories.FindAllInvoices(status, userNumber, campusNumber, schoolNumber, invoiceItems, startDate, endDate)
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

	statement, _ := invoiceRepositories.GetStatementByNumber(invoice.StatementNumber)

	student, err := repositories.FindStudentByUserNumber(statement.UserNumber)
	if err != nil {
		fmt.Println("Error getting student from user number " + err.Error())
		return responses.MinInvoiceResponse{}
	}

	if student.Number == "" {
		return responses.MinInvoiceResponse{}
	}

	billable, err := invoiceRepositories.GetBillableByInvoiceNumber(invoice.Number)
	if err != nil {
		fmt.Println("Error getting student from user number " + err.Error())
	}

	campusName := ""
	studyPeriodName := ""

	if err != nil && billable.Number != "" {
		campus, err := repositories3.GetCampusByNumber(billable.CampusNumber)
		if err == nil {
			campusName = campus.Name
		}

		studyPeriod, err := repositories4.GetStudyPeriodNameByStudyPeriodNumber(billable.StudyPeriodNumber)
		if err == nil {
			studyPeriodName = studyPeriod.Name
		}
	}

	invoiceResponse := responses.MinInvoiceResponse{
		StudyPeriodName: studyPeriodName,
		RefNumber:       statement.RefNumber,
		Date:            statement.CreatedAt.Format(time.DateOnly),
		DueDate:         invoice.DueDate.Format(time.DateOnly),
		Number:          invoice.Number,
		TotalAmount:     _invoiceAmount,
		PaidAmount:      _paidAmount,
		Balance:         _invoiceAmount - _paidAmount,
		Status:          invoice.Status,
		Adn:             student.ADN,
		Lin:             student.Lin,
		Names:           student.Surname + " " + student.SecondName + " " + student.OtherName,
		CampusName:      campusName,
	}

	return invoiceResponse
}

func generateNewInvoiceNumber() string {
	now := time.Now()
	year := now.Format("06")  // YY format
	month := now.Format("01") // MM format

	// Generate random four-digit number
	rand.NewSource(time.Now().UnixNano())
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

	payments, err := invoiceRepositories.FindAllInvoiceTransactions(invoice.Number)
	if err != nil {
		return responses.InvoiceResponse{}, err
	}

	var invoicePaymentsResponse []responses.InvoicePaymentResponse

	for _, payment := range payments {

		paymentResponse := createPaymentResponse(payment)

		invoicePaymentsResponse = append(invoicePaymentsResponse, paymentResponse)

	}

	invoiceResponse := responses.InvoiceResponse{
		MinInvoiceResponse: minInvoice,
		Items:              invoiceItemsResponse,
		Payments:           invoicePaymentsResponse,
	}

	return invoiceResponse, nil
}

func createPaymentResponse(invoiceTransaction models.InvoiceTransaction) responses.InvoicePaymentResponse {

	_invoiceStatement, _ := invoiceRepositories.FindInvoiceByNumber(invoiceTransaction.InvoiceNumber)
	transaction, _ := invoiceRepositories.GetTransactionByNumber(invoiceTransaction.TransactionNumber)
	transactionStatement := transaction.Statement
	addedBy, _ := repository.GetUserByNumber(transactionStatement.AddedByNumber)

	paymentResponse := responses.InvoicePaymentResponse{
		Amount:           invoiceTransaction.AmountPaidOff,
		InvoiceRefNumber: _invoiceStatement.RefNumber,
		RecordedBy:       addedBy.FirstName,
		DateTime:         invoiceTransaction.CreatedAt.Format(time.DateOnly),
		Status:           constants.StatusActive,
		PaymentNumber:    transactionStatement.RefNumber,
	}
	return paymentResponse
}

func ReceiveInvoicePayment(invoicePaymentRequest requests.InvoicePaymentRequest, currentUserNumber string) (responses.MinInvoiceResponse, error) {

	if len(invoicePaymentRequest.Credits) >= 1 {
		return ReceiveInvoicePaymentByCredit(invoicePaymentRequest.InvoiceNumber, currentUserNumber, invoicePaymentRequest.Credits)
	}

	var minInvoiceResponse responses.MinInvoiceResponse

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoicePaymentRequest.InvoiceNumber)
	if err != nil {
		return minInvoiceResponse, err
	}

	if invoice.Number == "" {
		return minInvoiceResponse, errors.New("invalid invoice number: " + invoicePaymentRequest.InvoiceNumber)
	}

	invoiceStatement := invoice.Statement

	var dateTimePaid time.Time
	if dateTimePaid, err = time.Parse(time.DateOnly, invoicePaymentRequest.Date); err != nil {
		return minInvoiceResponse, errors.New("Date error: " + err.Error() + " for " + invoicePaymentRequest.Date)
	}

	invoiceAmount, _ := invoiceRepositories.GetInvoiceAmount(invoice.Number)

	paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)

	if paidAmount >= invoiceAmount {
		return minInvoiceResponse, errors.New("This invoice is fully paid ")
	}

	invoiceBalance := invoiceAmount - paidAmount

	payAmount := invoiceBalance
	payRequestAmount := invoicePaymentRequest.Amount

	if payRequestAmount < invoiceBalance {
		payAmount = payRequestAmount
	}

	trnRefNo := generateNewInvoiceNumber()

	nextBlc, err := ComputeNextBalance(payRequestAmount, constants.StatementTypeIn, invoiceStatement.UserNumber)
	if err != nil {
		return minInvoiceResponse, err
	}

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		SchoolNumber:  invoiceStatement.SchoolNumber,
		AddedByNumber: currentUserNumber,
		RefNumber:     trnRefNo,
		CampusNumber:  invoiceStatement.CampusNumber,
		UserNumber:    invoiceStatement.UserNumber,
		Type:          constants.StatementTypeIn,
		Description:   "Payment received on invoice " + invoiceStatement.RefNumber,
		Amount:        payRequestAmount,
		Balance:       nextBlc,
	}

	newTransaction := models.Transaction{
		StatementNumber: newStatement.Number,
		Number:          uuid.NewV4().String(),
		Date:            dateTimePaid,
		Status:          constants.StatusActive,
		AmountInCredit:  payRequestAmount - payAmount,
	}

	newInvoiceTransaction := models.InvoiceTransaction{
		InvoiceNumber:     invoice.Number,
		TransactionNumber: newTransaction.Number,
		AmountPaidOff:     payAmount,
		AddedByNumber:     currentUserNumber,
	}

	var newStatus string

	if payAmount >= invoiceBalance {
		newStatus = constants.StatusPaid
	} else {
		newStatus = constants.StatusPartiallyPaid
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		tx = invoiceTransactionRepository.UseTx(tx)

		err := invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		if _, err = invoiceRepositories.SaveTransaction(newTransaction, tx); err != nil {
			return errors.New("Adding new transaction error " + err.Error())
		}

		err = invoiceTransactionRepository.SaveInvoiceTransaction(newInvoiceTransaction)
		if err != nil {
			return errors.New("Attaching transaction to invoice error " + err.Error())
		}

		if newStatus != invoice.Status {
			invoice.Status = newStatus
			if err = invoiceRepositories.UpdateInvoice(invoice, tx); err != nil {
				return errors.New("Updating invoice failed : " + err.Error())
			}
		}

		return nil
	})

	return createMinInvoiceResponse(invoice), err
}
func ReceiveInvoicePaymentByCredit(invoiceNumber, currentUserNumber string, credits []string) (responses.MinInvoiceResponse, error) {

	var minInvoiceResponse responses.MinInvoiceResponse

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return minInvoiceResponse, err
	}

	if invoice.Number == "" {
		return minInvoiceResponse, errors.New("invalid invoice number: " + invoiceNumber)
	}

	invoiceAmount, _ := invoiceRepositories.GetInvoiceAmount(invoice.Number)

	paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)

	if paidAmount >= invoiceAmount {
		return minInvoiceResponse, errors.New("This invoice is fully paid ")
	}

	invoiceBalance := invoiceAmount - paidAmount

	var totalCreditAmount float64

	var creditTransactions = make([]models.Transaction, 0)

	for _, transactionNumber := range credits {

		availableCredit, err := GetTransactionCredit(transactionNumber)
		if err != nil {
			return minInvoiceResponse, err
		}

		if availableCredit.AmountInCredit > 0 {
			creditTransactions = append(creditTransactions, availableCredit)
			totalCreditAmount = totalCreditAmount + availableCredit.AmountInCredit
		}

	}

	if totalCreditAmount <= 0 {
		return minInvoiceResponse, errors.New("All the selected transactions have no credit ")
	}

	var newStatus string

	if totalCreditAmount >= invoiceBalance {
		newStatus = constants.StatusPaid
	} else {
		newStatus = constants.StatusPartiallyPaid
	}

	var payAmount float64

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		tx = invoiceTransactionRepository.UseTx(tx)

		for _, creditTransaction := range creditTransactions {

			payAmount = invoiceBalance

			if creditTransaction.AmountInCredit < invoiceBalance {
				payAmount = creditTransaction.AmountInCredit
			}

			err = invoiceTransactionRepository.SaveInvoiceTransaction(models.InvoiceTransaction{
				InvoiceNumber:     invoice.Number,
				TransactionNumber: creditTransaction.Number,
				AmountPaidOff:     payAmount,
				AddedByNumber:     currentUserNumber,
			})

			creditTransaction.AmountInCredit = creditTransaction.AmountInCredit - payAmount
			_, err := invoiceRepositories.SaveTransaction(creditTransaction, tx)
			if err != nil {
				return errors.New("Updating credit error ")
			}

			invoiceBalance = invoiceBalance - payAmount
		}

		if err != nil {
			return errors.New("Attaching transaction to invoice error " + err.Error())
		}

		if newStatus != invoice.Status {
			invoice.Status = newStatus
			if err = invoiceRepositories.UpdateInvoice(invoice, tx); err != nil {
				return errors.New("Updating invoice failed : " + err.Error())
			}
		}

		return nil
	})

	return createMinInvoiceResponse(invoice), err
}

func GetFinancialSummary(userNumber string) (responses.FinancialSummaryResponse, error) {

	var financialSummaryResponse = responses.FinancialSummaryResponse{}

	var errorDesc string

	paymentsCount, err := invoiceRepositories.CountAllUserPayments(userNumber)
	if err == nil {
		financialSummaryResponse.TotalPayments = paymentsCount
	} else {
		errorDesc = err.Error()
	}

	invoicesCount, err := invoiceRepositories.CountAllUserInvoices(userNumber)
	if err == nil {
		financialSummaryResponse.TotalInvoices = invoicesCount
	} else {
		errorDesc = err.Error()
	}

	invoiceAmount, err := invoiceRepositories.GetAllInvoicedAmount(userNumber)
	if err == nil {
		financialSummaryResponse.TotalInvoiceAmount = invoiceAmount
	} else {
		errorDesc = err.Error()
	}

	paymentAmount, err := invoiceRepositories.SumAllUserPayments(userNumber)
	if err == nil {
		financialSummaryResponse.TotalPaidAmount = paymentAmount
	} else {
		errorDesc = err.Error()
	}

	financialSummaryResponse.TotalBalanceAmount = financialSummaryResponse.TotalInvoiceAmount - financialSummaryResponse.TotalPaidAmount

	mostRecent := responses.MostRecent{}

	recentInvoice, err := invoiceRepositories.GetMostRecentInvoice(userNumber)
	if err == nil {
		minInvoice := createMinInvoiceResponse(recentInvoice)
		mostRecent.InvoiceAmount = minInvoice.TotalAmount
		mostRecent.InvoiceDate = minInvoice.Date
	} else {
		errorDesc = err.Error()
	}

	mostRecentPayment, err := invoiceRepositories.GetMostRecentTransaction(userNumber)
	if err == nil {
		mostRecent.PaidAmount = mostRecentPayment.Amount
		mostRecent.PaidDate = mostRecentPayment.CreatedAt.Format(time.DateOnly)
	} else {
		errorDesc = err.Error()
	}

	financialSummaryResponse.MostRecent = mostRecent

	return financialSummaryResponse, errors.New(errorDesc)
}

func GetAllPayments(userNumber string, refNumber string, schoolNumber string, startDate time.Time, endDate time.Time) ([]responses.TransactionResponse, error) {

	var invoicePaymentsResponse []responses.TransactionResponse

	transactions, err := invoiceRepositories.GetAllPayments(userNumber, schoolNumber, refNumber, startDate, endDate)

	if err != nil {
		return invoicePaymentsResponse, err
	}

	for _, transaction := range transactions {

		paymentResponse := createTransactionResponse(transaction)

		invoicePaymentsResponse = append(invoicePaymentsResponse, paymentResponse)

	}

	return invoicePaymentsResponse, nil
}

func createTransactionResponse(transaction models.Transaction) responses.TransactionResponse {
	statement, _ := invoiceRepositories.GetStatementByNumber(transaction.StatementNumber)
	addedBy, _ := repository.GetUserByNumber(statement.AddedByNumber)
	paidBy, _ := repository.GetUserByNumber(statement.UserNumber)
	return responses.TransactionResponse{
		Amount:         statement.Amount,
		ReceiptNumber:  statement.RefNumber,
		DateReceived:   statement.CreatedAt.Format(time.DateTime),
		DatePaid:       transaction.Date.Format(time.DateOnly),
		ReceivedBy:     addedBy.FirstName,
		PaidBy:         paidBy.FirstName,
		Description:    statement.Description,
		AmountInCredit: transaction.AmountInCredit,
	}
}

func DeleteInvoice(invoiceNumber, currentUser string) error {

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return errors.New("checking invoice number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid invoice number")
	}

	if invoice.Status != constants.StatusNotPaid {
		return errors.New("you can only delete a not paid invoice")
	}

	invoiceAmount, err := invoiceRepositories.GetInvoiceAmount(invoiceNumber)

	invoiceStatement, err := invoiceRepositories.GetStatementByNumber(invoice.StatementNumber)
	if err != nil {
		return errors.New("Checking statement error " + err.Error())
	}

	nextBalance, err := ComputeNextBalance(invoiceAmount, constants.StatementTypeIn, invoiceStatement.UserNumber)
	if err != nil {
		return errors.New("Computing next balance error " + err.Error())
	}

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Type:          constants.StatementTypeIn,
		Amount:        invoiceAmount,
		Balance:       nextBalance,
		SchoolNumber:  invoiceStatement.SchoolNumber,
		AddedByNumber: currentUser,
		RefNumber:     generateNewInvoiceNumber(),
		CampusNumber:  invoiceStatement.CampusNumber,
		UserNumber:    invoiceStatement.UserNumber,
		Description:   "Invoice deleted " + invoiceStatement.RefNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		err = invoiceRepositories.DeleteInvoice(invoice, tx)
		if err != nil {
			return errors.New("deleting invoice failed : " + err.Error())
		}

		return nil
	})

	return err
}

func CancelInvoice(invoiceNumber string, currentUserNumber string) error {

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return errors.New("checking invoice number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid invoice number")
	}

	invoiceAmount, err := invoiceRepositories.GetInvoiceAmount(invoice.Number)
	if err != nil {
		return errors.New("Checking amount error " + err.Error())
	}

	invoiceTransactions, err := invoiceTransactionRepository.FindInvoiceTransactions(invoiceNumber, "")
	if err != nil {
		return errors.New("Checking payments error " + err.Error())
	}

	invoiceStatement, err := invoiceRepositories.GetStatementByNumber(invoice.StatementNumber)
	if err != nil {
		return errors.New("Checking statement error " + err.Error())
	}

	nextBalance, err := ComputeNextBalance(invoiceAmount, constants.StatementTypeIn, invoiceStatement.UserNumber)
	if err != nil {
		return errors.New("Computing next balance error " + err.Error())
	}

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Type:          constants.StatementTypeIn,
		Amount:        invoiceAmount,
		Balance:       nextBalance,
		SchoolNumber:  invoiceStatement.SchoolNumber,
		AddedByNumber: currentUserNumber,
		RefNumber:     generateNewInvoiceNumber(),
		CampusNumber:  invoiceStatement.CampusNumber,
		UserNumber:    invoiceStatement.UserNumber,
		Description:   "Invoice cancelled " + invoiceStatement.RefNumber,
	}

	invoice.Status = constants.StatusCancelled

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		for _, invoiceTransaction := range invoiceTransactions {
			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing invoice payments error " + err.Error())
			}

			transaction, _ := invoiceRepositories.GetTransactionByNumber(invoiceTransaction.TransactionNumber, tx)
			transaction.AmountInCredit = transaction.AmountInCredit + invoiceTransaction.AmountPaidOff
			_, err = invoiceRepositories.SaveTransaction(transaction, tx)
			if err != nil {
				return errors.New("Updating credit error " + err.Error())
			}
		}

		err := invoiceRepositories.UpdateInvoice(invoice, tx)
		if err != nil {
			return errors.New("updating invoice failed : " + err.Error())
		}

		err = invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		return nil
	})

	return err
}

func CancelPayment(transactionNumber string, currentUser string) error {

	invoicePayment, err := invoiceRepositories.GetTransactionByNumber(transactionNumber)
	if err != nil {
		return errors.New("checking payment number error: " + err.Error())
	}

	if invoicePayment.Number == "" {
		return errors.New("invalid transaction number")
	}

	transactionStatement, err := invoiceRepositories.GetStatementByNumber(invoicePayment.StatementNumber)
	if err != nil {
		return errors.New("Checking transaction statement error " + err.Error())
	}

	invoiceTransactions, err := invoiceTransactionRepository.FindInvoiceTransactions("", transactionNumber)
	if err != nil {
		return errors.New("Checking invoice transactions error " + err.Error())
	}

	nextBlc, err := ComputeNextBalance(transactionStatement.Amount, constants.StatementTypeOut, transactionStatement.UserNumber)
	if err != nil {
		return errors.New("Checking next user balance error " + err.Error())
	}

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		SchoolNumber:  transactionStatement.SchoolNumber,
		AddedByNumber: currentUser,
		RefNumber:     generateNewInvoiceNumber(),
		CampusNumber:  transactionStatement.CampusNumber,
		UserNumber:    transactionStatement.UserNumber,
		Type:          constants.StatementTypeOut,
		Description:   "Payment cancelled " + transactionStatement.RefNumber,
		Amount:        transactionStatement.Amount,
		Balance:       nextBlc,
	}

	invoicePayment.Status = constants.StatusCancelled
	invoicePayment.AmountInCredit = transactionStatement.Amount

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		for _, invoiceTransaction := range invoiceTransactions {

			paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoiceTransaction.InvoiceNumber, tx)

			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing invoice transaction error " + err.Error())
			}

			invoice, _ := invoiceRepositories.FindInvoiceByNumber(invoiceTransaction.InvoiceNumber, tx)

			if (paidAmount - invoiceTransaction.AmountPaidOff) > 0 {
				invoice.Status = constants.StatusPartiallyPaid
			} else {
				invoice.Status = constants.StatusNotPaid
			}

			err = invoiceRepositories.SaveInvoice(invoice, tx)
			if err != nil {
				return errors.New("Updating invoice status error " + err.Error())
			}
		}

		_, err := invoiceRepositories.SaveTransaction(invoicePayment, tx)
		if err != nil {
			return errors.New("Updating transaction failed : " + err.Error())
		}

		err = invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		return nil
	})

	return err
}

func DeletePayment(transactionNumber, currentUser string) error {

	transaction, err := invoiceRepositories.GetTransactionByNumber(transactionNumber)
	if err != nil {
		return errors.New("checking payment number error: " + err.Error())
	}

	if transaction.Number == "" {
		return errors.New("invalid payment number")
	}

	invoiceTransactions, err := invoiceTransactionRepository.FindInvoiceTransactions("", transactionNumber)
	if err != nil {
		return errors.New("Checking invoice transactions error " + err.Error())
	}

	oldStatement, err := invoiceRepositories.GetStatementByNumber(transaction.StatementNumber)
	if err != nil {
		return errors.New("Checking statement error " + err.Error())
	}

	balance, err := ComputeNextBalance(oldStatement.Amount, constants.StatementTypeOut, oldStatement.UserNumber)
	if err != nil {
		return err
	}

	statement := models.Statement{
		Number:        uuid.NewV4().String(),
		Amount:        oldStatement.Amount,
		UserNumber:    oldStatement.UserNumber,
		Description:   "Payment deleted " + oldStatement.RefNumber,
		SchoolNumber:  oldStatement.SchoolNumber,
		CampusNumber:  oldStatement.CampusNumber,
		RefNumber:     generateNewInvoiceNumber(),
		Balance:       balance,
		Type:          constants.StatementTypeOut,
		AddedByNumber: currentUser,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		for _, invoiceTransaction := range invoiceTransactions {

			paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoiceTransaction.InvoiceNumber, tx)

			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing invoice transaction error " + err.Error())
			}

			invoice, _ := invoiceRepositories.FindInvoiceByNumber(invoiceTransaction.InvoiceNumber, tx)

			if (paidAmount - invoiceTransaction.AmountPaidOff) > 0 {
				invoice.Status = constants.StatusPartiallyPaid
			} else {
				invoice.Status = constants.StatusNotPaid
			}

			err = invoiceRepositories.SaveInvoice(invoice, tx)
			if err != nil {
				return errors.New("Updating invoice status error " + err.Error())
			}
		}

		err := invoiceRepositories.DeleteTransaction(transaction, tx)
		if err != nil {
			return errors.New("deleting transaction failed : " + err.Error())
		}

		err = invoiceRepositories.SaveStatement(statement, tx)

		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		return nil
	})

	return err
}

func ComputeNextBalance(newAmount float64, statementType string, userNumber string) (float64, error) {

	if statementType != constants.StatementTypeIn && statementType != constants.StatementTypeOut {
		return 0, errors.New("Invalid statement type: '" + statementType + "'")
	}

	recentStatement, err := invoiceRepositories.GetRecentStatementByUserNumber(userNumber)
	if err != nil {
		return 0, errors.New("Getting recent statement error " + err.Error())
	}

	balance := recentStatement.Balance - newAmount

	if statementType == "IN" {
		balance = recentStatement.Balance + newAmount
	}

	return balance, nil
}

func UpdateInvoice(updateInvoiceRequest requests.UpdateInvoiceRequest, currentUser string) error {

	invoice, err := invoiceRepositories.FindInvoiceByNumber(updateInvoiceRequest.InvoiceNumber)
	if err != nil {
		return errors.New("checking invoice number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid invoice number")
	}

	paidAmount, err := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)
	if err != nil {
		return errors.New("checking paid amount error: " + err.Error())
	}

	if paidAmount > 0 {
		return errors.New("you cannot update an invoice with payments")
	}

	dueDate, err := commons.ConvertStringToDate(updateInvoiceRequest.DueDate)
	if err != nil {
		return err
	}

	items, err := invoiceRepositories.GetInvoiceItems(invoice.Number)
	if err != nil {
		return errors.New("getting current invoice items failed : " + err.Error())
	}

	invoice.Description = updateInvoiceRequest.Description
	invoice.DueDate = dueDate

	tx := config.DB.Begin()
	err = invoiceRepositories.UpdateInvoice(invoice, tx)
	if err != nil {
		tx.Rollback()
		return errors.New("updating invoice failed : " + err.Error())
	}

	for _, item := range items {
		err = invoiceRepositories.DeleteInvoiceItem(item, tx)
		if err != nil {
			tx.Rollback()
			return errors.New("removing invoice item failed : " + err.Error())
		}
	}

	for _, itemUpdate := range updateInvoiceRequest.Items {

		item, err := itemRepositories.GetItemByNumber(itemUpdate.ItemNumber, tx)
		if err != nil {
			tx.Rollback()
			return errors.New("checking item failed : " + err.Error())
		}

		if item.Number == "" {
			tx.Rollback()
			return errors.New("invalid item number")
		}

		newInvoiceItem := models.InvoiceItem{
			Number:          uuid.NewV4().String(),
			InvoiceNumber:   invoice.Number,
			ItemNumber:      itemUpdate.ItemNumber,
			Quantity:        itemUpdate.Quantity,
			UnitPrice:       itemUpdate.UnitCost,
			ActualUnitPrice: item.Cost,
		}

		if _, err := invoiceRepositories.SaveInvoiceItem(newInvoiceItem, tx); err != nil {
			tx.Rollback()
			return err
		}

	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("commit error : " + err.Error())
		return errors.New("completing invoice editing failed")
	}

	return nil
}

func UserStatement(userNumber string, startDate time.Time, endDate time.Time) ([]responses.StatementResponse, error) {

	var statementResponse = make([]responses.StatementResponse, 0)

	allStatements, err := invoiceRepositories.FindAllStatements(userNumber, startDate, endDate)
	if err != nil {
		return statementResponse, errors.New("Error getting statements " + err.Error())
	}

	for _, statement := range allStatements {

		receivedBy, _ := repository.GetUserByNumber(statement.AddedByNumber)

		statementResponse = append(statementResponse, responses.StatementResponse{
			Amount:       statement.Amount,
			Balance:      statement.Balance,
			Type:         statement.Type,
			Description:  statement.Description,
			DateReceived: statement.CreatedAt.Format(time.DateTime),
			ReceivedBy:   receivedBy.FirstName,
		})
	}

	return statementResponse, nil
}

func UserCredits(userNumber string) ([]responses.TransactionResponse, error) {

	var creditsResponse = make([]responses.TransactionResponse, 0)

	userTransactions, err := invoiceRepositories.AllUserCredit(userNumber)
	if err != nil {
		return creditsResponse, errors.New("Getting user transactions error " + err.Error())
	}

	for _, transaction := range userTransactions {
		creditsResponse = append(creditsResponse, createTransactionResponse(transaction))
	}

	return creditsResponse, nil
}

func GetTransactionCredit(transactionNumber string) (models.Transaction, error) {

	transaction, err := invoiceRepositories.GetTransactionByNumber(transactionNumber)
	if err != nil {
		return transaction, errors.New("Error getting transaction " + err.Error())
	}

	if transaction.Number == "" {
		return transaction, errors.New("Invalid transaction number ")
	}

	creditsUsed, err := invoiceRepositories.TransactionUsedCredit(transactionNumber)
	if err != nil {
		return transaction, errors.New("Error getting transaction used credit " + err.Error())
	}

	statement := transaction.Statement

	transaction.AmountInCredit = statement.Amount - creditsUsed

	return transaction, nil
}
