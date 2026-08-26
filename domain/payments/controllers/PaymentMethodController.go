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
		Account      string `json:"account"`
		OtherDetails string `json:"other_details"`
		CampusNumber string `json:"campus_number"`
	}

	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	campusNumber := payload.CampusNumber
	if campusNumber == "" {
		campusNumber = c.Request.Header.Get("Campus-Number")
	}
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus is required"})
		return
	}

	savePaymentMethod := models.PaymentMethod{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		Account:      payload.Account,
		OtherDetails: payload.OtherDetails,
		CampusNumber: campusNumber,
		UserNumber:   user.Number,
	}

	var duplicatePaymentMethod models.PaymentMethod
	if err := config.DB.Where("name = ? AND campus_number = ?", payload.Name, campusNumber).First(&duplicatePaymentMethod).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method already exists in this campus"})
		return
	}

	if err := config.DB.Create(&savePaymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method created"})
}

func EditPaymentMethod(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	type Payload struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Account      string `json:"account"`
		OtherDetails string `json:"other_details"`
		CampusNumber string `json:"campus_number"`
	}

	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	var method models.PaymentMethod
	if err := config.DB.Where("number = ?", payload.ID).First(&method).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method not found"})
		return
	}

	method.Name = payload.Name
	method.Account = payload.Account
	method.OtherDetails = payload.OtherDetails
	if payload.CampusNumber != "" {
		method.CampusNumber = payload.CampusNumber
	}

	if err := config.DB.Save(&method).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method updated"})
}

func FetchPaymentMethods(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := c.Request.Header.Get("Campus-Number")

	var paymentMethods []models.PaymentMethod
	query := config.DB.Preload("Campus")
	if campusNumber != "" {
		query = query.Where("campus_number = ?", campusNumber)
	}

	if err := query.Find(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch payment methods"})
		return
	}

	type PaymentMethodResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Account      string `json:"account"`
		OtherDetails string `json:"other_details"`
		CampusName   string `json:"campus_name"`
		CampusNumber string `json:"campus_number"`
	}

	var response []PaymentMethodResponse
	for _, paymentMethod := range paymentMethods {
		response = append(response, PaymentMethodResponse{
			ID:           paymentMethod.Number,
			Name:         paymentMethod.Name,
			Account:      paymentMethod.Account,
			OtherDetails: paymentMethod.OtherDetails,
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

	paymentMethodNumber := c.Query("id")
	if paymentMethodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method id is required"})
		return
	}

	var method models.PaymentMethod
	if err := config.DB.Where("number = ?", paymentMethodNumber).First(&method).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment method not found"})
		return
	}

	if err := config.DB.Delete(&method).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted"})
}
