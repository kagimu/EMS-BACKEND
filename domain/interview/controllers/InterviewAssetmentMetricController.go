package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/service"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateInterviewAssessmentMetric(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var payload requests.InterviewAssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

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
		Number:                 uuid.NewV4().String(),
		InterviewNumber:        payload.InterviewNumber,
		AssessmentMetricNumber: payload.AssessmentMetricNumber,
	}

	if err := config.DB.Create(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetric, "message": "Interview assessment metric created successfully"})
}

func FetchInterviewAssessmentMetrics(c *gin.Context) {
	var interviewAssessmentMetrics []models.InterviewAssessmentMetric
	interviewNumber := c.Param("interview_number")
	if err := config.DB.Scopes(filters.ByInterviewNumber(interviewNumber)).Preload("AssessmentMetric").Preload("Interview").Find(&interviewAssessmentMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interview assessment metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetrics})
}

func FetchInterviewAssessmentMetricByStudent(c *gin.Context) {

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student number is missing"})
		return
	}

	interviewNumber := c.Query("interview_number")
	if interviewNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview number is missing"})
		return
	}

	results, err := service.FetchInterviewAssessmentMetricByStudent(studentNumber, interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func RemoveInterviewAssessmentMetric(c *gin.Context) {

	interviewAssessmentMetricNumber := c.Query("interview_assessment_metric_number")
	if interviewAssessmentMetricNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric number is missing"})
		return
	}

	var interviewAssessmentMetric models.InterviewAssessmentMetric
	if err := config.DB.Where("number = ?", interviewAssessmentMetricNumber).First(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric not found"})
		return
	}

	if err := config.DB.Delete(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete interview assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview assessment metric deleted successfully"})
}
