package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePermission(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload models.Permission
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	savePermission := models.Permission{
		Number:     uuid.NewV4().String(),
		Name:       payload.Name,
		UserNumber: user.Number,
	}

	if err := config.DB.Create(&savePermission).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": savePermission})
}

func FetchPermissions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

func FetchPermissionGroups(c *gin.Context) {
	permissions := commons.PermissionGroups()
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch permissions"})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"data": permissions})

}
