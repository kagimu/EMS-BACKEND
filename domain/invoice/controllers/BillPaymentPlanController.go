package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"github.com/gin-gonic/gin"
)

func CreateBillPaymentPlan(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.BillPaymentPlanRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := services.CreateBillPaymentPlan(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Bill payment plan created successfully"))
}

func FetchBillPaymentPlans(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	planNumber := c.Query("bill_payment_plan_number")
	promissoryNoteNumber := c.Query("promissory_note_number")
	studentNumber := c.Query("student_number")
	status := c.Query("status")

	data, err := services.GetBillPaymentPlans(planNumber, promissoryNoteNumber, studentNumber, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func UpdateBillPaymentPlan(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	planNumber := c.Param("bill_payment_plan_number")
	if planNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "provide bill_payment_plan_number"})
		return
	}

	var payload requests.UpdateBillPaymentPlanRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := services.UpdateBillPaymentPlan(planNumber, payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Bill payment plan updated successfully"))
}

func PayBillPaymentPlanInstallment(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.PayInstallmentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := services.PayBillPaymentPlanInstallment(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Payment processed successfully"))
}
