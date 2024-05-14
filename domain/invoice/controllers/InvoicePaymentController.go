package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
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

	var minInvoiceResponse responses.MinInvoiceResponse
	minInvoiceResponse, err = services.ReceiveInvoicePayment(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Failed saving payment with error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice payment received successfully", "data": minInvoiceResponse})
}

func GetAllPayments(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	refNumber := context.Query("ref_number")
	schoolNumber := context.Query("school_number")
	userNumber := context.Query("user_number")

	if refNumber == "" && schoolNumber == "" && userNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide at least the ref number, user number or school number"))
		return
	}

	payments, err := services.GetAllPayments(userNumber, refNumber, schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Getting all payments with error: "+err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": payments})
}
