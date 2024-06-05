package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
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

	if err = services.CreateNewInvoice(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice created successfully"})
}

func UpdateInvoice(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	type Input struct {
		Status string `json:"status"`
		Type   string `json:"type"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	studentNumber := c.Param("student_number")
	var student studentModel.Student
	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No student found with the specified number"})
		return
	}

	paymentReference := c.Query("payment_reference")
	var invoice models.Invoice
	if err := config.DB.First(&invoice, "payment_reference  = ?", paymentReference).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No invoice found with the specified number"})
		return
	}

	updateInvoice := models.Invoice{
		Status: payload.Status,
		Type:   payload.Type,
	}

	if err := config.DB.Model(&invoice).Where("number = ?", invoice.Number).Updates(updateInvoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update invoice", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully", "data": invoice})
}

func GetAllInvoices(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
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

	currentDateTime := time.Now()

	if _startDate == "" {
		startOfMonth := time.Date(currentDateTime.Year(), currentDateTime.Month(), 1, 0, 0, 0, 0, currentDateTime.Location())
		_startDate = startOfMonth.Format(time.DateOnly)
	}

	var startDate time.Time
	if startDate, err = time.Parse(time.DateOnly, _startDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
		return
	}

	if _endDate == "" {
		_endDate = currentDateTime.Format(time.DateOnly)
	}

	var endDate time.Time
	if endDate, err = time.Parse(time.DateOnly, _endDate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
		return
	}

	invoices, err := services.AllInvoices(status, userNumber, campusNumber, schoolNumber, itemNumbers, startDate, endDate)

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
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide the invoice ref number"))
		return
	}

	details, err := services.InvoiceDetails(refNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Getting details failed, try again: "+err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Invoice details", "data": details})
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

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Getting financial summary failed, try again: "+err.Error()))
	}

	context.JSON(http.StatusOK, gin.H{"message": "Financial summary", "data": summary})

}

func DeleteInvoice(context *gin.Context) {
	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	invoiceNumber := context.Query("invoice_number")
	if invoiceNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide invoice number"))
		return
	}

	err := services.DeleteInvoice(invoiceNumber)

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	}

	context.JSON(http.StatusOK, commons.Feedback("Invoice removed successfully"))
}
