package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
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

	var payload requests.ExaminationTypeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	saveExaminationType := models.ExaminationType{
		Number:      uuid.NewV4().String(),
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := config.DB.Create(&saveExaminationType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create examination"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveExaminationType})
}

func FetchExaminationTypes(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var examinationTypes []models.ExaminationType

	if err := config.DB.Find(&examinationTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch examination types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": examinationTypes})
}
