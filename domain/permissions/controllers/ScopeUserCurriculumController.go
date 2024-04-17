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

func CreateScopeUserCurriculum(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeUserCurriculumRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if err := payload.ValidateCurriculum(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveCampusScope := models.ScopeUserCurriculum{
		Number:           uuid.NewV4().String(),
		AddedByNumber:    payload.AddedByNumber,
		CurriculumNumber: payload.CurriculumNumber,
		UserNumber:       user.Number,
	}

	if err := config.DB.Create(&saveCampusScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save curriculum scope"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curriculum scope created successfully"})
}

func FetchScopeUserCurriculum(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var curriculumScopes []models.ScopeUserCurriculum
	if err := config.DB.Preload("Campus").Find(&curriculumScopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch curriculum scopes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": curriculumScopes})
}
