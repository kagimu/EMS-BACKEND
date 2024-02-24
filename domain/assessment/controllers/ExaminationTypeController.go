package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateExaminationType(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload models.ExaminationType
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required"})
		return
	}

	saveExaminationType := models.ExaminationType{
		Number: uuid.NewV4().String(),
		Name:   payload.Name,
	}

	if err := config.DB.Create(&saveExaminationType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create examination type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveExaminationType})
}
