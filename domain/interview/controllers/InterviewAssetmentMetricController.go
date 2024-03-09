package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateInterviewAssessmentMetric(c *gin.Context) {
	var payload requests.InterviewAssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	err := payload.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// validate if interview number  is valid
	var interview models.Interview
	if err := config.DB.Where("number = ?", payload.InterviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview schedule not found"})
		return
	}

	var assessmentMetric assessmentModel.AssessmentMetric
	if err := config.DB.Where("number = ?", payload.AssessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric not found"})
		return
	}

	var interviewAssessmentMetric models.InterviewAssessmentMetric
	if err := config.DB.Where("interview_number = ? AND assessment_metric_number = ?", payload.InterviewNumber, payload.AssessmentMetricNumber).First(&interviewAssessmentMetric).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric already exists"})
		return
	}

	interviewAssessmentMetric = models.InterviewAssessmentMetric{
		InterviewNumber:        payload.InterviewNumber,
		AssessmentMetricNumber: payload.AssessmentMetricNumber,
		UserNumber:             c.GetHeader("User-Number"),
	}

	if err := config.DB.Create(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetric, "message": "Interview assessment metric created successfully"})
}

func FetchInterviewAssessmentMetrics(c *gin.Context) {
	var interviewAssessmentMetrics []models.InterviewAssessmentMetric
	if err := config.DB.Preload("AssessmentMetric").Preload("Interview").Find(&interviewAssessmentMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interview assessment metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetrics})
}
