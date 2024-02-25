package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/responses"
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

	var payload requests.NewRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	roleNumber := uuid.NewV4().String()
	saveRole := models.Role{
		Number:      roleNumber,
		Name:        payload.Name,
		UserNumber:  user.Number,
		Description: payload.Description,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveRole).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save role"})
		tx.Rollback()
		return
	}

	// save the permissions
	for _, permission := range payload.Permissions {
		savePermission := models.RolePermission{
			RoleNumber: roleNumber,
			Permission: permission,
		}

		if err := tx.Create(&savePermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save role permissions"})
			tx.Rollback()
			return
		}

	}

	tx.Commit()

	c.JSON(http.StatusOK, commons.Feedback("Role has been created"))
}

func FetchRoles(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var rolesResponse []responses.RolesList

	var roles []models.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch roles"})
		return
	}

	// fetch the permissions
	for _, role := range roles {
		var permissions []models.RolePermission
		// select only the permission column in the role_permissions table
		if err := config.DB.Where("role_number = ?", role.Number).Find(&permissions).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch role details"})
			return
		}

		//prepare array of permissions
		var pms []string
		for _, permission := range permissions {
			pms = append(pms, permission.Permission)
		}

		rolesResponse = append(rolesResponse, responses.RolesList{
			Number:      role.Number,
			Name:        role.Name,
			Description: role.Description,
			Permissions: pms,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": rolesResponse})
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
