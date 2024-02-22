package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateRole(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload models.Role
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}
	if payload.SchoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "School number is required"})
		return
	}

	saveRole := models.Role{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		UserNumber:   user.Number,
		SchoolNumber: payload.SchoolNumber,
	}

	if err := config.DB.Create(&saveRole).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveRole})
}

func FetchRoles(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var roles []models.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}
