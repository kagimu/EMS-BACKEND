package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAssessmentMetric(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	formatDate, err := time.Parse("2006-01-02", payload.ScheduledOn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse time"})
		return
	}

	saveAssessmentMetric := models.AssessmentMetric{
		Number:                uuid.NewV4().String(),
		GradingStandardNumber: payload.GradingStandardNumber,
		StudyPeriodNumber:     payload.StudyPeriodNumber,
		ScheduledOn:           formatDate,
		CampusClassNumber:     payload.CampusClassNumber,
		MetricNumber:          payload.MetricNumber,
		AssessmentAreaNumber:  payload.AssessmentAreaNumber,
		Weight:                payload.Weight,
		Description:           payload.Description,
		MaxScore:              payload.MaxScore,
	}

	if err := config.DB.Create(&saveAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAssessmentMetric})
}

func FetchAssessmentMetrics(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var assessmentAreaNumber = c.Query("assessment_area_number")
	var studyPeriodNumber = c.Query("study_period_number")

	var assessmentMetrics []models.AssessmentMetric
	if err := config.DB.Model(models.AssessmentMetric{}).Scopes(
		filters.ByStudyPeriodNumber(studyPeriodNumber),
		filters.ByAssessmentAreaNumber(assessmentAreaNumber),
	).Preload("AssessmentArea").Preload("Metric").Preload("CampusClass").Find(&assessmentMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessment  metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentMetrics})
}
