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

func CreateScopeUserCampus(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserCampusRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateCampus(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveCampusScope := models.ScopeUserCampus{
		Number:        uuid.NewV4().String(),
		AddedByNumber: payload.AddedByNumber,
		CampusNumber:  payload.CampusNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveCampusScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save campus scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campus scope created successfully"})
}

func FetchScopeUserCampus(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var campusScopes []models.ScopeUserCampus
	if err := config.DB.Preload("Campus").Find(&campusScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch campus scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": campusScopes})
}
