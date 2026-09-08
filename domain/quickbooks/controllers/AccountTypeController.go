package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func CreateAccountType(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountType := models.QBAccountType{
		Number:      uuid.NewV4().String(),
		Name:        input.Name,
		Description: input.Description,
		IsActive:    true,
	}

	if err := config.DB.Create(&accountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account type"})
		return
	}

	c.JSON(http.StatusCreated, accountType)
}

func GetAccountTypes(c *gin.Context) {
	var accountTypes []models.QBAccountType
	if err := config.DB.Find(&accountTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account types"})
		return
	}
	c.JSON(http.StatusOK, accountTypes)
}

func UpdateAccountType(c *gin.Context) {
	id := c.Param("id")
	var accountType models.QBAccountType

	if err := config.DB.Where("number = ?", id).First(&accountType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account type not found"})
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

	accountType.Name = input.Name
	accountType.Description = input.Description
	accountType.IsActive = input.IsActive

	if err := config.DB.Save(&accountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account type"})
		return
	}

	c.JSON(http.StatusOK, accountType)
}

func DeleteAccountType(c *gin.Context) {
	id := c.Param("id")
	var accountType models.QBAccountType

	if err := config.DB.Where("number = ?", id).First(&accountType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account type not found"})
		return
	}

	if err := config.DB.Delete(&accountType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account type deleted successfully"})
}
