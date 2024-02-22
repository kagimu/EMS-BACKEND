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

func AssignRoles(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if user.Role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not an admin"})
		return
	}

	type Input struct {
		Name       string `json:"name"`
		UserNumber string `json:"user_number"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	if payload.UserNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User number is required"})
		return
	}

	updateUser := models.User{
		Role: payload.Name,
	}

	if err := config.DB.Model(&user).Where("number = ?", payload.UserNumber).Updates(updateUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign a role to the specified user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been specified"})
}
