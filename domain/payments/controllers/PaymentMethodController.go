package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePaymentMethod(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	type Payload struct {
		Name         string `json:"name"`
		CampusNumber string `json:"campus_number"`
	}

	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	savePaymentMethod := models.PaymentMethod{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		CampusNumber: payload.CampusNumber,
		UserNumber:   user.Number,
	}

	var duplicatePaymentMethod models.PaymentMethod
	if err := config.DB.Where("name = ?", payload.Name).First(&duplicatePaymentMethod).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method already exists"})
		return
	}

	if err := config.DB.Create(&savePaymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method created"})
}

func FetchPaymentMethods(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var paymentMethods []models.PaymentMethod
	if err := config.DB.Preload("Campus").Find(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch payment methods"})
		return
	}

	type PaymentMethodResponse struct {
		Name         string `json:"name"`
		CampusName   string `json:"campus_name"`
		CampusNumber string `json:"campus_number"`
	}

	var response []PaymentMethodResponse
	for _, paymentMethod := range paymentMethods {
		response = append(response, PaymentMethodResponse{
			Name:         paymentMethod.Name,
			CampusName:   paymentMethod.Campus.Name,
			CampusNumber: paymentMethod.Campus.Number,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func DeletePaymentMethod(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	paymentMethodNumber := c.Param("number")
	if paymentMethodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method number is required"})
		return
	}

	var bank models.Bank
	if err := config.DB.Where("number = ?", paymentMethodNumber).First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method not found"})
		return
	}

	if err := config.DB.Delete(&bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted"})
}
