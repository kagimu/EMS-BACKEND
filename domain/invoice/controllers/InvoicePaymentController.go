package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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

	//var minInvoiceResponse responses.MinInvoiceResponse
	minInvoiceResponse, err := services.ReceiveInvoicePayment(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Failed saving payment with error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promissory note payment received successfully", "data": minInvoiceResponse})
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

	payments, err := services.GetAllPayments(userNumber, refNumber, schoolNumber, startDate, endDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Getting all payments with error: "+err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": payments})
}

func CancelPayment(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	paymentNumber := context.Query("payment_number")
	if paymentNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("provide a payment number to cancel"))
		return
	}

	err = services.CancelPayment(paymentNumber, user.Number)

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Payment cancelled successfully"))

}

func DeletePayment(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
		return
	}

	paymentNumber := context.Query("payment_number")
	if paymentNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("provide a payment number to delete"))
		return
	}

	err = services.DeletePayment(paymentNumber, user.Number)

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	}

	context.JSON(http.StatusOK, commons.Feedback("Payment deleted successfully"))
}

func AllStatements(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := context.Query("user_number")
	if userNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide user number "})
		return
	}

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

	statementResponses, err := services.UserStatement(userNumber, startDate, endDate)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get statements, try again", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": statementResponses})
}

func TransactionCredits(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := context.Query("user_number")
	if userNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide user number "})
		return
	}

	statementResponses, err := services.UserCredits(userNumber)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get statements, try again", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": statementResponses})
}

func GeneralStatement(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := context.Query("user_number")
	minAmount_ := context.Query("min_amount")
	maxAmount_ := context.Query("max_amount")

	minAmount, err := parseFloat64Ptr(minAmount_)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing min_amount to float: " + err.Error()})
		return
	}

	maxAmount, err := parseFloat64Ptr(maxAmount_)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing max_amount to float: " + err.Error()})
		return
	}

	statementResponses, err := services.GeneralStatement(userNumber, minAmount, maxAmount)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get general statements, try again", "err": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": statementResponses})
}

func parseFloat64Ptr(val string) (float64, error) {
	if val == "" {
		return -1, nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return -1, err
	}
	return f, nil
}
