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

func CreateMetric(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveMetric := models.Metric{
		Number:           uuid.NewV4().String(),
		Name:             payload.Name,
		CurriculumNumber: payload.CurriculumNumber,
		UserNumber:       user.Number,
	}

	var duplicateMetric models.Metric
	if duplicate := config.DB.Where("name = ?", saveMetric.Name).First(&duplicateMetric).Error; duplicate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name already exists"})
		return
	}

	if err := config.DB.Create(&saveMetric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveMetric})
}

func FetchMetrics(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var metrics []models.Metric
	if err := config.DB.Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": metrics})
}
