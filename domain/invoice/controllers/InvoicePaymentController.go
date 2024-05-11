package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/gin-gonic/gin"
)

func CreateInvoicePayment(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	var payload requests.InvoicePaymentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide a valid payment request"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate payment request: %s", err.Error())})
		return
	}

	if err = services.ReceiveInvoicePayment(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Failed saving payment, try again "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice payment received successfully"})
}
