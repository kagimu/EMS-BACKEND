package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateChartOfAccount(c *gin.Context) {
	var input struct {
		Name          string `json:"name" binding:"required"`
		AccountTypeID string `json:"account_type_id" binding:"required"`
		Description   string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chartOfAccount := models.QBChartOfAccount{
		Number:        uuid.NewV4().String(),
		Name:          input.Name,
		AccountTypeID: input.AccountTypeID,
		Description:   input.Description,
		IsActive:      true,
	}

	if err := config.DB.Create(&chartOfAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chart of account"})
		return
	}

	c.JSON(http.StatusCreated, chartOfAccount)
}

func GetChartOfAccounts(c *gin.Context) {
	var chartOfAccounts []models.QBChartOfAccount
	if err := config.DB.Find(&chartOfAccounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chart of accounts"})
		return
	}
	c.JSON(http.StatusOK, chartOfAccounts)
}

func UpdateChartOfAccount(c *gin.Context) {
	id := c.Param("id")
	var chartOfAccount models.QBChartOfAccount

	if err := config.DB.Where("number = ?", id).First(&chartOfAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chart of account not found"})
		return
	}

	var input struct {
		Name          string `json:"name" binding:"required"`
		AccountTypeID string `json:"account_type_id" binding:"required"`
		Description   string `json:"description"`
		IsActive      bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chartOfAccount.Name = input.Name
	chartOfAccount.AccountTypeID = input.AccountTypeID
	chartOfAccount.Description = input.Description
	chartOfAccount.IsActive = input.IsActive

	if err := config.DB.Save(&chartOfAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chart of account"})
		return
	}

	c.JSON(http.StatusOK, chartOfAccount)
}

func DeleteChartOfAccount(c *gin.Context) {
	id := c.Param("id")
	var chartOfAccount models.QBChartOfAccount

	if err := config.DB.Where("number = ?", id).First(&chartOfAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chart of account not found"})
		return
	}

	if err := config.DB.Delete(&chartOfAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chart of account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chart of account deleted successfully"})
}
