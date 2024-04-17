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

func CreateScopeUserLevel(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserLevelRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateLevel(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveCampusScope := models.ScopeUserLevel{
		Number:        uuid.NewV4().String(),
		AddedByNumber: payload.AddedByNumber,
		LevelNumber:   payload.LevelNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveCampusScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save level scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Level scope created successfully"})
}

func FetchScopeUserLevel(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var levelScopes []models.ScopeUserLevel
	if err := config.DB.Preload("Campus").Find(&levelScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch level scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": levelScopes})
}
