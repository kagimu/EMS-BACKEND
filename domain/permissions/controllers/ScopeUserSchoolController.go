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

func CreateScopeUserSchool(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserSchoolRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateSchool(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveCampusScope := models.ScopeUserSchool{
		Number:        uuid.NewV4().String(),
		AddedByNumber: payload.AddedByNumber,
		SchoolNumber:  payload.SchoolNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveCampusScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "School scope created successfully"})
}

func FetchScopeUserSchool(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var schoolScopes []models.ScopeUserSchool
	if err := config.DB.Preload("Campus").Find(&schoolScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch school scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": schoolScopes})
}

func DeleteSchoolUserScope(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	scopeNumber := c.Query("scope_number")
	if scopeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Scope number is required"})
		return
	}

	var schoolScope models.ScopeUserSchool
	if err := config.DB.Where("number = ?", scopeNumber).First(&schoolScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch school scope"})
		return
	}

	if err := config.DB.Delete(&schoolScope).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete school scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "School scope deleted successfully"})
}
