package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type SettingsRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

func GetSettings(c *gin.Context) {
	var settings models.QBSettings
	if err := config.DB.First(&settings).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"username": "", "has_password": false})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"username": settings.Username,
		"has_password": settings.PasswordHash != "",
	})
}

func UpdateSettings(c *gin.Context) {
	var req SettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var settings models.QBSettings
	if err := config.DB.First(&settings).Error; err != nil {
		// Create new
		settings = models.QBSettings{
			Number:       uuid.NewV4().String(),
			Username:     req.Username,
			PasswordHash: req.Password,
			IsActive:     true,
			SchoolNumber: "default", // Need auth token context for this in production
		}
		if err := config.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
			return
		}
	} else {
		// Update existing
		settings.Username = req.Username
		if req.Password != "" {
			settings.PasswordHash = req.Password
		}
		if err := config.DB.Save(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "QuickBooks settings updated successfully"})
}

func GetQueueStatus(c *gin.Context) {
	var count int64
	config.DB.Model(&models.QBSyncQueue{}).Where("status = ?", "QUEUED").Count(&count)
	c.JSON(http.StatusOK, gin.H{"queued_items": count})
}
