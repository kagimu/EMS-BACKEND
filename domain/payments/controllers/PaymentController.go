package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePayment(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	var payload models.Payment

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to process input",
		})
		return
	}

	applicationNumber := payload.ApplicationNumber
	var application applicationModel.Application
	if err := config.DB.Where("number = ?", applicationNumber).First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Student not found"})
		return
	}

	savePayment := models.Payment{
		Number:            uuid.NewV4().String(),
		ApplicationNumber: applicationNumber,
		Amount:            payload.Amount,
		Status:            payload.Status,
		SlipSerialNumber:  payload.SlipSerialNumber,
		ApprovedBy:        payload.ApprovedBy,
		Details:           payload.Details,
		Date:              payload.Date,
		Method:            payload.Method,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&savePayment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to save payment information",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save payment information"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": savePayment,
	})
}
