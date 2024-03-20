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

func CreateAssessmentSubMetric(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentSubMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	saveAssessmentSubMetric := models.AssessmentSubMetric{
		Number: uuid.NewV4().String(),
		Name:   payload.Name,
	}

	if err := config.DB.Create(&saveAssessmentSubMetric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create assessment sub metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAssessmentSubMetric})
}

func FetchAssessmentSubMetrics(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var assessmentSubMetrics []models.AssessmentSubMetric
	if err := config.DB.Model(&models.AssessmentSubMetric{}).Find(&assessmentSubMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessment sub metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentSubMetrics})
}
