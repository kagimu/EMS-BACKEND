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
		//ExamType:         payload.ExamType,
		UserNumber: user.Number,
	}

	var duplicateMetric models.Metric
	if duplicate := config.DB.Where("name = ? AND curriculum_number = ? ", saveMetric.Name, saveMetric.CurriculumNumber).First(&duplicateMetric).Error; duplicate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name already exists on this curriculum"})
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

	curricullumNumber := c.Query("curriculum_number")

	if curricullumNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum number is required"})
		return

	}
	var metrics []models.Metric
	if err := config.DB.Where("curriculum_number = ? ", curricullumNumber).Find(&metrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": metrics})
}

func EditMetric(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	metricNumber := c.Param("metric_number")
	var payload requests.MetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var metric models.Metric
	if err := config.DB.Where("number = ?", metricNumber).First(&metric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Metric not found"})
		return
	}

	updateMetric := models.Metric{
		Name:             payload.Name,
		CurriculumNumber: payload.CurriculumNumber,
		//ExamType:         payload.ExamType,
	}

	if err := config.DB.Model(&metric).Where("number = ?", metricNumber).Updates(updateMetric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": metric})
}

func DeleteMetric(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	metricNumber := c.Param("metric_number")
	var metric models.Metric

	err := config.DB.Where("number = ?", metricNumber).First(&metric).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Metric not found"})
		return
	}

	if err := config.DB.Delete(&metric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metric deleted successfully"})
}
