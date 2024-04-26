package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

	studentNumber := payload.StudentNumber
	var student studentModel.Student
	if err := config.DB.First(&student, "number  = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No student found with the specified number"})
		return
	}

	invoiceNumber, err := InvoiceNumberGenerator()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate invoice number", "err": err.Error()})
		return
	}

	saveInvoice := models.Invoice{
		Number:              uuid.NewV4().String(),
		StudentNumber:       payload.StudentNumber,
		InvoiceUniqueNumber: invoiceNumber,
		DueDate:             payload.DueDate,
		Type:                payload.Type,
		Description:         payload.Description,
		UserNumber:          user.Number,
	}

	if err := config.DB.Create(&saveInvoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create invoice"})
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
