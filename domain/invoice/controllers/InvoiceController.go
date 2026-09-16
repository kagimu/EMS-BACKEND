package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"net/http"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/gin-gonic/gin"
)

func CreateInvoice(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.InvoiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	invoiceUUID, err := services.CreateInvoice(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promissory note created successfully", "number": invoiceUUID})
}

func UpdateInvoice(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = services.UpdateInvoice(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promissory note updated successfully"})
}

func GetAllInvoices(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := context.Query("school_number")
	status := context.Query("status")
	campusNumber := context.Query("campus_number")
	userNumber := context.Query("user_number")
	itemNumbers := context.Query("items")

	_startDate := context.Query("start_date")
	_endDate := context.Query("end_date")

	var startDate time.Time
	if _startDate != "" {
		if startDate, err = time.Parse(time.DateOnly, _startDate); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
			return
		}
	}

	var endDate time.Time
	if _endDate != "" {
		if endDate, err = time.Parse(time.DateOnly, _endDate); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
			return
		}
	}

	invoices, err := services.AllInvoices(status, userNumber, campusNumber, schoolNumber, itemNumbers, startDate, endDate, user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get invoices, try again", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, invoices)
}

func GetInvoiceDetails(context *gin.Context) {

	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	refNumber := context.Query("ref_number")

	if refNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide the promissory note ref number"))
		return
	}

	details, err := services.InvoiceDetails(refNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Getting details failed, try again: "+err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Promissory note details", "data": details})
}

func FinancialSummary(context *gin.Context) {

	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := context.Query("user_number")
	if userNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide user number to view financial summary"))
		return
	}

	summary, err := services.GetFinancialSummary(userNumber)

	var errorDesc string
	if err != nil {
		errorDesc = err.Error()
	}

	context.JSON(http.StatusOK, gin.H{"message": "Financial summary", "data": summary, "error": errorDesc})
}

func DeleteInvoice(context *gin.Context) {

	var user *models.User

	user, err := commons.GetUserFromHeader(context)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	invoiceNumber := context.Query("invoice_number")
	if invoiceNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide promissory note number to delete"))
		return
	}

	err = services.DeleteInvoice(invoiceNumber, user.Number)

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	}

	context.JSON(http.StatusOK, commons.Feedback("Promissory note removed successfully"))
}

func CancelInvoice(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	invoiceNumber := context.Query("invoice_number")
	if invoiceNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide promissory note number to cancel"))
		return
	}

	err = services.CancelInvoice(invoiceNumber, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Promissory note cancelled successfully"))
}

func GetInvoiceItemTransactionsHistory(context *gin.Context) {
	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	invoiceNumber := context.Query("invoice_number")
	if invoiceNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide promissory note number"))
		return
	}

	history, err := services.GetInvoiceItemTransactionHistory(invoiceNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item transaction history", "data": history})
}

func GetInvoiceItemTransactionsHistoryByTransaction(context *gin.Context) {
	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	transactionNumber := context.Query("transaction_number")
	if transactionNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide transaction number"))
		return
	}

	history, err := services.GetInvoiceItemTransactionHistoryByTransaction(transactionNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Item transaction history", "data": history})
}

func GetAllItemTransactionsHistory(context *gin.Context) {
	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := context.Query("campus_number")
	search := context.Query("search")
	startDate := context.Query("start_date")
	endDate := context.Query("end_date")
	
	pageStr := context.Query("page")
	perPageStr := context.Query("per_page")
	
	page := 1
	perPage := 20
	
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
		perPage = pp
	}

	history, err := services.GetAllItemTransactionHistory(campusNumber, search, startDate, endDate, page, perPage)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "All item transactions history", "data": history})
}
