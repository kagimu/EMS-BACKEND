package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	saveRole := models.Role{
		Number:      uuid.NewV4().String(),
		Name:        payload.Name,
		UserNumber:  user.Number,
		Description: payload.Description,
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

func AssignRolesToAUser(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	type Input struct {
		UserNumber string   `json:"user_number"`
		RoleNumber []string `json:"role_number"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if len(payload.RoleNumber) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one role is required"})
		return
	}

	if payload.UserNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User number is required"})
		return
	}

	for _, roleNumber := range payload.RoleNumber {
		saveRole := models.UserRole{
			RoleNumber: roleNumber,
			UserNumber: payload.UserNumber,
			AddedBy:    user.Number,
		}

		if err := config.DB.Save(&saveRole); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign a role to the specified user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
}

func FetchRoleDetails(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	roleNumber := c.Param("role_number")

	var role models.Role
	if err := config.DB.Where("number = ?", roleNumber).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch role details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func RemoveRoleAssignment(c *gin.Context) {

	type Input struct {
		UserNumber string   `json:"user_number"`
		RoleNumber []string `json:"role_number"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if len(payload.RoleNumber) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one role is required"})
		return
	}

	// delete the role assignment

	if err := config.DB.Where("user_number = ?", payload.UserNumber).Where("role_number = ?", payload.RoleNumber).Delete(&models.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove role assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assignment has been updated"})
}
