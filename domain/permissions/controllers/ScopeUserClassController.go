package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateScopeUserClass(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserClassRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateClass(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveClassScope := models.ScopeUserClass{
		Number:        uuid.NewV4().String(),
		AddedByNumber: payload.AddedByNumber,
		ClassNumber:   payload.ClassNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveClassScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save class scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class scope created successfully"})
}

func FetchScopeUserClass(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var classScopes []models.ScopeUserClass

	if err := config.DB.Preload("Class").Find(&classScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch class scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": classScopes})
}

func DeleteClassUserScope(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	scopeNumber := c.Query("scope_number")
	if scopeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Scope number is required"})
		return
	}

	var classScope models.ScopeUserClass
	if err := config.DB.Where("number = ?", scopeNumber).First(&classScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch clalss scope"})
		return
	}

	if err := config.DB.Delete(&classScope).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete clalss scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class scope deleted successfully"})
}
