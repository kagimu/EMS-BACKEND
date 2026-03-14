package services

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
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
)

var invoiceTransactionRepository = invoiceRepositories.NewInvoiceTransactionRepository()

var admissionRepository = repositories2.NewAdmissionRepository()

func CreateInvoice(payload requests.InvoiceRequest, currentUserNumber string) error {

	dueDate, err := commons.ConvertStringToDate(payload.DueDate)
	if err != nil {
		return err
	}

	issueDate, err := commons.ConvertStringToDate(payload.IssueDate)
	if err != nil {
		return err
	}

	if dueDate.Before(issueDate) {
		return errors.New("due date should be after the date of issue")
	}

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
		return errors.New("try again to generate a new promissory note number")
	}

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

	newStatement := models.Statement{
		Number:        uuid.NewV4().String(),
		Amount:        invoiceTotal,
		UserNumber:    payload.UserNumber,
		Description:   "New bill " + invoiceNumber,
		SchoolNumber:  admission.SchoolNumber,
		CampusNumber:  admission.CampusNumber,
		RefNumber:     invoiceNumber,
		Balance:       0,
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

		balance, err := ComputeNextBalance(invoiceTotal, constants.StatementTypeOut, userNumber, tx)
		if err != nil {
			return err
		}

		newStatement.Balance = balance

		err = invoiceRepositories.SaveStatement(newStatement, tx)
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

func CreateInvoiceFromBillable(billable models.Billable, student studentModel.Student, clientReference, description, currentUserNumber string, paidAmount float64, datePaid string, balanceDueDate string, txn *gorm.DB) (models.Invoice, error) {

	invoiceNumber := generateNewInvoiceNumber()

	invoice, err := invoiceRepositories.FindInvoiceByRefNumber(invoiceNumber)

	if err != nil {
		return models.Invoice{}, err
	}

	if invoice.Number != "" {
		return models.Invoice{}, errors.New("try again to generate a new promissory note number")
	}

	dueDate := time.Now().AddDate(0, 0, 14)
	if balanceDueDate != "" {
		if parsed, err := commons.ConvertStringToDate(balanceDueDate); err == nil {
			dueDate = parsed
		}
	}

	billableItems, err := invoiceRepositories.FindBillableItems(billable.Number, txn)
	if err != nil {
		return models.Invoice{}, errors.New("Getting billable items error: " + err.Error())
	}

	var invoiceTotal float64

	for _, billItem := range billableItems {
		invoiceTotal = invoiceTotal + (billItem.Quantity * billItem.UnitPrice)
	}

	balance, err := ComputeNextBalance(invoiceTotal, constants.StatementTypeOut, student.UserNumber, txn)

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
		IssueDate:     time.Now(),
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
		return models.Invoice{}, errors.New("Creating promissory note error: " + err.Error())
	}

	err = invoiceRepositories.SaveBillableInvoice(models.BillableInvoice{
		BillableNumber: billable.Number,
		InvoiceNumber:  newInvoice.Number,
	}, txn)

	if err != nil {
		return models.Invoice{}, errors.New("Creating attaching promissory note to billable error: " + err.Error())
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

	_, err := ReceiveEnrollmentInvoicePayment(requests.InvoicePaymentRequest{
		InvoiceNumber: invoice.Number,
		Amount:        paidAmount,
		Date:          datePaid,
	}, currentUserNumber, txn)

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

func AllInvoices(status string, userNumber string, campusNumber string, schoolNumber string, commaSeperatedItems string, startDate time.Time, endDate time.Time, user *userModel.User) ([]responses.MinInvoiceResponse, error) {

	var invoiceResponses = make([]responses.MinInvoiceResponse, 0)

	var invoiceItems = make([]string, 0)

	if commaSeperatedItems != "" {
		invoiceItems = strings.Split(commaSeperatedItems, ",")
	}

	invoices, err := invoiceRepositories.FindAllInvoices(status, userNumber, campusNumber, schoolNumber, invoiceItems, startDate, endDate, user)
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
		fmt.Println("Getting promissory note amount error '" + err.Error() + "' for '" + invoice.RefNumber + "'")
		_invoiceAmount = 0
	} else {
		_invoiceAmount = invoiceAmount
	}

	if paidAmount, err := invoiceRepositories.GetInvoicePaidAmount(invoice.Number); err != nil {
		fmt.Println("Getting promissory note paid amount error '" + err.Error() + "' for '" + invoice.RefNumber + "'")
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
		IssueDate:       statement.IssueDate.Format(time.DateOnly),
		Number:          invoice.Number,
		TotalAmount:     _invoiceAmount,
		PaidAmount:      _paidAmount,
		Balance:         _invoiceAmount - _paidAmount,
		Status:          invoice.Status,
		Adn:             student.ADN,
		Lin:             student.Lin,
		Names:           student.Surname + " " + student.SecondName + " " + student.OtherName,
		CampusName:      campusName,
		StudentNumber:   student.Number,
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
		return responses.InvoiceResponse{}, errors.New("invalid promissory note ref number")
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
		PaymentDate:      transactionStatement.IssueDate.Format(time.DateOnly),
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
		return minInvoiceResponse, errors.New("invalid promissory note number while receiving payment: " + invoicePaymentRequest.InvoiceNumber)
	}

	invoiceStatement := invoice.Statement

	paymentDate, err := validatePaymentDate(invoicePaymentRequest.Date, invoiceStatement.IssueDate)
	if err != nil {
		return minInvoiceResponse, err
	}

	invoiceAmount, _ := invoiceRepositories.GetInvoiceAmount(invoice.Number)

	paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)

	if paidAmount >= invoiceAmount {
		return minInvoiceResponse, errors.New("This promissory note is fully paid ")
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
		Description:   "Payment received on promissory note " + invoiceStatement.RefNumber,
		Amount:        payRequestAmount,
		Balance:       nextBlc,
		IssueDate:     paymentDate,
	}

	newTransaction := models.Transaction{
		StatementNumber: newStatement.Number,
		Number:          uuid.NewV4().String(),
		Date:            paymentDate,
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
			return errors.New("Attaching transaction to promissory note error " + err.Error())
		}

		if newStatus != invoice.Status {
			invoice.Status = newStatus
			if err = invoiceRepositories.UpdateInvoice(invoice, tx); err != nil {
				return errors.New("Updating promissory note failed, after receiving payment: " + err.Error())
			}
		}

		return nil
	})

	return createMinInvoiceResponse(invoice), err
}
func ReceiveEnrollmentInvoicePayment(invoicePaymentRequest requests.InvoicePaymentRequest, currentUserNumber string, tx *gorm.DB) (responses.MinInvoiceResponse, error) {

	if len(invoicePaymentRequest.Credits) >= 1 {
		return ReceiveInvoicePaymentByCredit(invoicePaymentRequest.InvoiceNumber, currentUserNumber, invoicePaymentRequest.Credits)
	}

	var minInvoiceResponse responses.MinInvoiceResponse

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoicePaymentRequest.InvoiceNumber, tx)
	if err != nil {
		return minInvoiceResponse, err
	}

	if invoice.Number == "" {
		return minInvoiceResponse, errors.New("invalid promissory note number while receiving payment: " + invoicePaymentRequest.InvoiceNumber)
	}

	invoiceStatement := invoice.Statement

	dateTimePaid, err := validatePaymentDate(invoicePaymentRequest.Date, invoiceStatement.IssueDate)
	if err != nil {
		return minInvoiceResponse, err
	}

	invoiceAmount, _ := invoiceRepositories.GetInvoiceAmount(invoice.Number, tx)

	paidAmount := 0.0

	if paidAmount >= invoiceAmount {
		return minInvoiceResponse, errors.New("This promissory note is fully paid ")
	}

	invoiceBalance := invoiceAmount - paidAmount

	payAmount := invoiceBalance
	payRequestAmount := invoicePaymentRequest.Amount

	if payRequestAmount < invoiceBalance {
		payAmount = payRequestAmount
	}

	trnRefNo := generateNewInvoiceNumber()

	nextBlc, err := ComputeNextBalance(payRequestAmount, constants.StatementTypeIn, invoiceStatement.UserNumber, tx)
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
		Description:   "Payment received on promissory note " + invoiceStatement.RefNumber,
		Amount:        payRequestAmount,
		Balance:       nextBlc,
		IssueDate:     dateTimePaid,
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

	tx = invoiceTransactionRepository.UseTx(tx)

	err = invoiceRepositories.SaveStatement(newStatement, tx)

	if err != nil {
		return minInvoiceResponse, errors.New("Creating statement error " + err.Error())
	}

	if _, err = invoiceRepositories.SaveTransaction(newTransaction, tx); err != nil {
		return minInvoiceResponse, errors.New("Adding new transaction error " + err.Error())
	}

	err = invoiceTransactionRepository.SaveInvoiceTransaction(newInvoiceTransaction)
	if err != nil {
		return minInvoiceResponse, errors.New("Attaching transaction to promissory note error " + err.Error())
	}

	if newStatus != invoice.Status {
		invoice.Status = newStatus
		if err = invoiceRepositories.UpdateInvoice(invoice, tx); err != nil {
			return minInvoiceResponse, errors.New("Updating promissory note failed while enrolling: " + err.Error())
		}
	}

	return createMinInvoiceResponse(invoice), err
}

func validatePaymentDate(date string, dateOfIssue time.Time) (time.Time, error) {
	var paymentDate time.Time
	if paymentDate, err := time.Parse(time.DateOnly, date); err != nil {
		return paymentDate, errors.New("Date error: " + err.Error() + " for " + date)
	}

	if paymentDate.Before(dateOfIssue) {
		return paymentDate, errors.New(fmt.Sprintf("The payment date '%s' can not be before the issue date '%s'", paymentDate.Format(time.DateOnly), dateOfIssue.Format(time.DateOnly)))
	}

	return paymentDate, nil
}

func ReceiveInvoicePaymentByCredit(invoiceNumber, currentUserNumber string, credits []string) (responses.MinInvoiceResponse, error) {

	var minInvoiceResponse responses.MinInvoiceResponse

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return minInvoiceResponse, err
	}

	if invoice.Number == "" {
		return minInvoiceResponse, errors.New("invalid promissory note number: " + invoiceNumber)
	}

	invoiceAmount, _ := invoiceRepositories.GetInvoiceAmount(invoice.Number)

	paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)

	if paidAmount >= invoiceAmount {
		return minInvoiceResponse, errors.New("This promissory note is fully paid ")
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
			return errors.New("Attaching transaction to promissory note error " + err.Error())
		}

		if newStatus != invoice.Status {
			invoice.Status = newStatus
			if err = invoiceRepositories.UpdateInvoice(invoice, tx); err != nil {
				return errors.New("Updating promissory note failed while using credit: " + err.Error())
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

	totalBalance := financialSummaryResponse.TotalInvoiceAmount - financialSummaryResponse.TotalPaidAmount
	if totalBalance >= 0 {
		financialSummaryResponse.TotalBalanceAmount = totalBalance
	} else {
		financialSummaryResponse.TotalCreditAmount = -1 * totalBalance
	}

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
		mostRecent.PaidDate = mostRecentPayment.IssueDate.Format(time.DateOnly)
		mostRecent.DateReceived = mostRecentPayment.CreatedAt.Format(time.DateOnly)
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
		Amount:            statement.Amount,
		ReceiptNumber:     statement.RefNumber,
		DateReceived:      statement.CreatedAt.Format(time.DateTime),
		DatePaid:          transaction.Date.Format(time.DateOnly),
		ReceivedBy:        addedBy.FirstName,
		PaidBy:            paidBy.FirstName,
		Description:       statement.Description,
		AmountInCredit:    transaction.AmountInCredit,
		TransactionNumber: transaction.Number,
		Status:            transaction.Status,
	}
}

func DeleteInvoice(invoiceNumber, currentUser string) error {

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return errors.New("checking promissory note number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid promissory note number")
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
		Description:   "promissory note deleted " + invoiceStatement.RefNumber,
		IssueDate:     time.Now(),
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err := invoiceRepositories.SaveStatement(newStatement, tx)
		if err != nil {
			return errors.New("Creating statement error " + err.Error())
		}

		err = invoiceRepositories.DeleteInvoice(invoice, tx)
		if err != nil {
			return errors.New("deleting promissory note failed : " + err.Error())
		}

		return nil
	})

	return err
}

func CancelInvoice(invoiceNumber string, currentUserNumber string) error {

	invoice, err := invoiceRepositories.FindInvoiceByNumber(invoiceNumber)
	if err != nil {
		return errors.New("checking promissory note number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid promissory note number")
	}

	if invoice.Status == constants.StatusCancelled {
		return errors.New("This promissory note is already cancelled ")
	}

	invoiceAmount, err := invoiceRepositories.GetInvoiceAmount(invoice.Number)
	if err != nil {
		return errors.New("Checking amount error " + err.Error())
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
		Description:   "Promissory note cancelled " + invoiceStatement.RefNumber,
		IssueDate:     time.Now(),
	}

	invoice.Status = constants.StatusCancelled

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		invoiceTransactions, err := invoiceTransactionRepository.FindInvoiceTransactions(invoiceNumber, "")
		if err != nil {
			return errors.New("Checking payments error " + err.Error())
		}

		for _, invoiceTransaction := range invoiceTransactions {
			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing promissory note payments error " + err.Error())
			}

			transaction, _ := invoiceRepositories.GetTransactionByNumber(invoiceTransaction.TransactionNumber, tx)
			transaction.AmountInCredit = transaction.AmountInCredit + invoiceTransaction.AmountPaidOff

			updatedInvoice, err := invoiceRepositories.SaveTransaction(transaction, tx)
			if err != nil {
				return errors.New("Updating credit error " + err.Error())
			}

			fmt.Println("Promissory note updated status to: " + updatedInvoice.Status)
		}

		err = invoiceRepositories.UpdateInvoice(invoice, tx)
		if err != nil {
			return errors.New("updating promissory note failed after cancellation: " + err.Error())
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

	transaction, err := invoiceRepositories.GetTransactionByNumber(transactionNumber)
	if err != nil {
		return errors.New("checking payment number error: " + err.Error())
	}

	if transaction.Number == "" {
		return errors.New("invalid transaction number")
	}

	if transaction.Status == constants.StatusCancelled {
		return errors.New("This transaction was already cancelled ")
	}

	transactionStatement, err := invoiceRepositories.GetStatementByNumber(transaction.StatementNumber)
	if err != nil {
		return errors.New("Checking transaction statement error " + err.Error())
	}

	if transactionStatement.Number == "" {
		return errors.New("Transaction has no corresponding statement ")
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
		IssueDate:     time.Now(),
	}

	transaction.Status = constants.StatusCancelled
	transaction.AmountInCredit = transactionStatement.Amount

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		invoiceTransactions, err := invoiceTransactionRepository.FindInvoiceTransactions("", transactionNumber)
		if err != nil {
			return errors.New("Checking promissory note transactions error " + err.Error())
		}

		for _, invoiceTransaction := range invoiceTransactions {

			paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoiceTransaction.InvoiceNumber, tx)

			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing promissory note transaction error " + err.Error())
			}

			invoice, _ := invoiceRepositories.FindInvoiceByNumber(invoiceTransaction.InvoiceNumber, tx)

			if (paidAmount - invoiceTransaction.AmountPaidOff) > 0 {
				invoice.Status = constants.StatusPartiallyPaid
			} else {
				invoice.Status = constants.StatusNotPaid
			}

			err = invoiceRepositories.SaveInvoice(invoice, tx)
			if err != nil {
				return errors.New("Updating promissory note status error " + err.Error())
			}
		}

		_, err = invoiceRepositories.SaveTransaction(transaction, tx)
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
		return errors.New("Checking promissory note transactions error " + err.Error())
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
		IssueDate:     time.Now(),
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		invoiceTransactionRepository.UseTx(tx)

		for _, invoiceTransaction := range invoiceTransactions {

			paidAmount, _ := invoiceRepositories.GetInvoicePaidAmount(invoiceTransaction.InvoiceNumber, tx)

			err := invoiceTransactionRepository.DeleteTransaction(invoiceTransaction)
			if err != nil {
				return errors.New("Clearing promissory note transaction error " + err.Error())
			}

			invoice, _ := invoiceRepositories.FindInvoiceByNumber(invoiceTransaction.InvoiceNumber, tx)

			if (paidAmount - invoiceTransaction.AmountPaidOff) > 0 {
				invoice.Status = constants.StatusPartiallyPaid
			} else {
				invoice.Status = constants.StatusNotPaid
			}

			err = invoiceRepositories.SaveInvoice(invoice, tx)
			if err != nil {
				return errors.New("Updating promissory note status error " + err.Error())
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

func ComputeNextBalance(newAmount float64, statementType string, userNumber string, tx ...*gorm.DB) (float64, error) {

	if statementType != constants.StatementTypeIn && statementType != constants.StatementTypeOut {
		return 0, errors.New("Invalid statement type: '" + statementType + "'")
	}

	recentStatement, err := invoiceRepositories.GetRecentStatementByUserNumber(userNumber, tx...)
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

	dueDate, err := commons.ConvertStringToDate(updateInvoiceRequest.DueDate)
	if err != nil {
		return err
	}

	issueDate, err := commons.ConvertStringToDate(updateInvoiceRequest.IssueDate)
	if err != nil {
		return err
	}

	if dueDate.Before(issueDate) {
		return errors.New("due date should be after the date of issue")
	}

	invoice, err := invoiceRepositories.FindInvoiceByNumber(updateInvoiceRequest.InvoiceNumber)
	if err != nil {
		return errors.New("checking promissory note number error: " + err.Error())
	}

	if invoice.Number == "" {
		return errors.New("invalid promissory note number")
	}

	paidAmount, err := invoiceRepositories.GetInvoicePaidAmount(invoice.Number)
	if err != nil {
		return errors.New("checking paid amount error: " + err.Error())
	}

	if paidAmount > 0 {
		return errors.New("you cannot update a promissory note with payments")
	}

	items, err := invoiceRepositories.GetInvoiceItems(invoice.Number)
	if err != nil {
		return errors.New("getting current promissory note items failed : " + err.Error())
	}

	invoice.Description = updateInvoiceRequest.Description
	invoice.DueDate = dueDate

	tx := config.DB.Begin()
	err = invoiceRepositories.UpdateInvoice(invoice, tx)
	if err != nil {
		tx.Rollback()
		return errors.New("updating promissory note failed: " + err.Error())
	}

	for _, item := range items {
		err = invoiceRepositories.DeleteInvoiceItem(item, tx)
		if err != nil {
			tx.Rollback()
			return errors.New("removing promissory note item failed : " + err.Error())
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
		return errors.New("completing promissory note editing failed")
	}

	return nil
}

func UserStatement(userNumber string, startDate, endDate time.Time) ([]responses.StatementResponse, error) {

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
			IssueDate:    statement.IssueDate.Format(time.DateTime),
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

func GeneralStatement(userNumber string, minAmount, maxAmount float64) ([]responses.StatementResponse, error) {
	var statementResponse = make([]responses.StatementResponse, 0)

	allStatements, err := invoiceRepositories.FindAllGeneralStatements(userNumber, minAmount, maxAmount)
	if err != nil {
		return statementResponse, errors.New("Error getting statements " + err.Error())
	}

	for _, statement := range allStatements {

		receivedBy, _ := repository.GetUserByNumber(statement.AddedByNumber)
		learner, _ := repository.GetUserByNumber(statement.UserNumber)
		student, err := repositories.FindStudentByUserNumber(statement.UserNumber)
		if err != nil {
			return statementResponse, errors.New("Error getting student from user number " + statement.UserNumber)
		}

		statementResponse = append(statementResponse, responses.StatementResponse{
			Amount:            statement.Amount,
			Balance:           statement.Balance,
			Type:              statement.Type,
			Description:       statement.Description,
			DateReceived:      statement.CreatedAt.Format(time.DateTime),
			ReceivedBy:        receivedBy.FirstName,
			LearnerName:       fmt.Sprintf("%s %s %s", learner.FirstName, learner.SecondName, learner.LastName),
			LearnerUserNumber: statement.UserNumber,
			StudentNumber:     student.Number,
			IssueDate:         statement.IssueDate.Format(time.DateTime),
		})
	}

	return statementResponse, nil
}
