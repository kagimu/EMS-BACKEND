package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func CreateDepositAccount(c *gin.Context) {
	var input struct {
		Name                 string `json:"name" binding:"required"`
		AccountTypeID        string `json:"account_type_id"`
		DepositAccountTypeID string `json:"deposit_account_type_id"`
		Details              string `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountTypeID := input.AccountTypeID
	if accountTypeID == "" {
		accountTypeID = input.DepositAccountTypeID
	}

	depositAccount := models.QBDepositAccount{
		Number:               uuid.NewV4().String(),
		Name:                 input.Name,
		AccountTypeID:        accountTypeID,
		DepositAccountTypeID: input.DepositAccountTypeID,
		Details:              input.Details,
		IsActive:             true,
	}

	if err := config.DB.Create(&depositAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deposit account"})
		return
	}

	config.DB.Preload("AccountType").Preload("DepositAccountType").First(&depositAccount, "number = ?", depositAccount.Number)

	c.JSON(http.StatusCreated, depositAccount)
}

func GetDepositAccounts(c *gin.Context) {
	var depositAccounts []models.QBDepositAccount
	if err := config.DB.Preload("AccountType").Preload("DepositAccountType").Find(&depositAccounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deposit accounts"})
		return
	}
	c.JSON(http.StatusOK, depositAccounts)
}

func UpdateDepositAccount(c *gin.Context) {
	id := c.Param("id")
	var depositAccount models.QBDepositAccount

	if err := config.DB.Where("number = ?", id).First(&depositAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit account not found"})
		return
	}

	var input struct {
		Name                 string `json:"name" binding:"required"`
		AccountTypeID        string `json:"account_type_id"`
		DepositAccountTypeID string `json:"deposit_account_type_id"`
		Details              string `json:"details"`
		IsActive             bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountTypeID := input.AccountTypeID
	if accountTypeID == "" {
		accountTypeID = input.DepositAccountTypeID
	}

	depositAccount.Name = input.Name
	depositAccount.AccountTypeID = accountTypeID
	if input.DepositAccountTypeID != "" {
		depositAccount.DepositAccountTypeID = input.DepositAccountTypeID
	}
	depositAccount.Details = input.Details
	depositAccount.IsActive = input.IsActive

	if err := config.DB.Save(&depositAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deposit account"})
		return
	}

	config.DB.Preload("AccountType").Preload("DepositAccountType").First(&depositAccount, "number = ?", depositAccount.Number)

	c.JSON(http.StatusOK, depositAccount)
}

func DeleteDepositAccount(c *gin.Context) {
	id := c.Param("id")
	var depositAccount models.QBDepositAccount

	if err := config.DB.Where("number = ?", id).First(&depositAccount).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deposit account not found"})
		return
	}

	if err := config.DB.Delete(&depositAccount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deposit account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit account deleted successfully"})
}
