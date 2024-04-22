package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateBank(c *gin.Context) {
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

	saveBank := models.Bank{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		CampusNumber: payload.CampusNumber,
		UserNumber:   user.Number,
	}

	var duplicateBank models.Bank
	if err := config.DB.Where("name = ?", payload.Name).First(&duplicateBank).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bank already exists"})
		return
	}

	if err := config.DB.Create(&saveBank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method created"})
}

func FetchBanks(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var banks []models.Bank
	if err := config.DB.Preload("Campus").Find(&banks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch banks"})
		return
	}

	type BankResponse struct {
		Name         string `json:"name"`
		CampusName   string `json:"campus_name"`
		CampusNumber string `json:"campus_number"`
	}

	var response []BankResponse
	for _, bank := range banks {
		response = append(response, BankResponse{
			Name:         bank.Name,
			CampusName:   bank.Campus.Name,
			CampusNumber: bank.Campus.Number,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func DeleteBank(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	bankNumber := c.Param("number")
	if bankNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bank number is required"})
		return
	}

	var bank models.Bank
	if err := config.DB.Where("number = ?", bankNumber).First(&bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bank not found"})
		return
	}

	if err := config.DB.Delete(&bank).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete bank"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank deleted"})
}
