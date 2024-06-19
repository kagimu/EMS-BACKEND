package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateScope(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ScopeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// validate if user has existing scope

	var scope models.Scope
	if err := config.DB.Where("user_number = ?", payload.UserNumber).First(&scope).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already has a scope please delete the existing scope first"})
		return

	}

	saveScope := models.Scope{
		Number:            uuid.NewV4().String(),
		CampusNumber:      payload.CampusNumber,
		SchoolNumber:      payload.SchoolNumber,
		CurriculumNumber:  payload.CurriculumNumber,
		CampusClassNumber: payload.ClassNumber,
		StreamNumber:      payload.StreamNumber,
		UserNumber:        payload.UserNumber,
		AddedBy:           user.Number,
	}

	if err := config.DB.Create(&saveScope).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save scope"})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Scope has been created"))
}

func FetchScopes(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var scopes []models.Scope
	if err := config.DB.Preload("User").Preload("Campus").Preload("Curriculum").Preload("School").Preload("Class").Preload("Stream").Find(&scopes).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch scopes"})
		return
	}

	var response []responses.ScopeResponse
	for _, scope := range scopes {
		response = append(response, responses.ScopeResponse{
			Number:           scope.Number,
			CampusNumber:     scope.CampusNumber,
			SchoolNumber:     scope.SchoolNumber,
			CurriculumNumber: scope.CurriculumNumber,
			ClassNumber:      scope.CampusClassNumber,
			StreamNumber:     scope.StreamNumber,
			UserNumber:       scope.UserNumber,
			UserName:         scope.User.FirstName + " " + scope.User.SecondName + scope.User.LastName,
			CampusName:       scope.Campus.Name,
			SchoolName:       scope.School.Name,
			ClassName:        scope.Class.ClassName,
			StreamName:       scope.Stream.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
