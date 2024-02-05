package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoicePayment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoicePayment/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateInvoicePayment(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	var payload requests.InvoicePaymentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	saveInvoicePayment := models.InvoicePayment{
		Number:           uuid.NewV4().String(),
		PaymentReference: payload.PaymentReference,
		Amount:           payload.Amount,
		Balance:          payload.Balance,
		Payer:            payload.Payer,
		Date:             payload.Date,
	}
	if err := config.DB.Create(&saveInvoicePayment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create invoice payment", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice payment created successfully"})

}
