package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func CreateDepositAccountType(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depositAccountType := models.QBDepositAccountType{
		Number:      uuid.NewV4().String(),
		Name:        input.Name,
		Description: input.Description,
		IsActive:    true,
	}

	if err := config.DB.Create(&depositAccountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deposit account type"})
		return
	}

	c.JSON(http.StatusCreated, depositAccountType)
}

func GetDepositAccountTypes(c *gin.Context) {
	var depositAccountTypes []models.QBDepositAccountType
	if err := config.DB.Find(&depositAccountTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deposit account types"})
		return
	}
	c.JSON(http.StatusOK, depositAccountTypes)
}

func UpdateDepositAccountType(c *gin.Context) {
	id := c.Param("id")
	var depositAccountType models.QBDepositAccountType

	if err := config.DB.Where("number = ?", id).First(&depositAccountType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit account type not found"})
		return
	}

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depositAccountType.Name = input.Name
	depositAccountType.Description = input.Description
	depositAccountType.IsActive = input.IsActive

	if err := config.DB.Save(&depositAccountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deposit account type"})
		return
	}

	c.JSON(http.StatusOK, depositAccountType)
}

func DeleteDepositAccountType(c *gin.Context) {
	id := c.Param("id")
	var depositAccountType models.QBDepositAccountType

	if err := config.DB.Where("number = ?", id).First(&depositAccountType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit account type not found"})
		return
	}

	if err := config.DB.Delete(&depositAccountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deposit account type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit account type deleted successfully"})
}
