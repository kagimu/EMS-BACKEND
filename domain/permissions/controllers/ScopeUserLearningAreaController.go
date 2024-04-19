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

func CreateScopeUserLearningArea(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserLearningAreaRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateLearningArea(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveCampusScope := models.ScopeUserLearningArea{
		Number:             uuid.NewV4().String(),
		AddedByNumber:      payload.AddedByNumber,
		LearningAreaNumber: payload.LearningAreaNumber,
		UserNumber:         user.Number,
	}

	if err := config.DB.Create(&saveCampusScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save learning area"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Learning area scope created successfully"})
}

func FetchScopeUserLearningArea(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var learningAreaScopes []models.ScopeUserLearningArea
	if err := config.DB.Preload("Campus").Find(&learningAreaScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch learning area scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": learningAreaScopes})
}

func DeleteLearningAreaUserScope(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	scopeNumber := c.Query("scope_number")
	if scopeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Scope number is required"})
		return
	}

	var learningarea models.ScopeUserLearningArea
	if err := config.DB.Where("number = ?", scopeNumber).First(&learningarea).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch learning area scope"})
		return
	}

	if err := config.DB.Delete(&learningarea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete learning area scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Learning area scope deleted successfully"})
}
