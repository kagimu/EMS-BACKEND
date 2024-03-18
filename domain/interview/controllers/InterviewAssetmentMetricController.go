package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func FetchInterviewAssessmentMetricByStudent(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	var student studentModel.Student

	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student number is missing"})
		return
	}

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	var assessment []assessmentModel.Assessment
	if err := config.DB.Where("student_number = ?", studentNumber).Preload("AssessmentMetric").Preload("Student").First(&assessment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assessment not found"})
		return
	}

	interviewNumber := c.Query("interview_number")
	var interview models.Interview

	if interviewNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview number is missing"})
		return
	}

	if err := config.DB.First(&interview, "number = ?", interviewNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Interview not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	var interviewAssessmentMetric models.InterviewAssessmentMetric
	if err := config.DB.Where("interview_number = ?", interviewNumber).Preload("Interview").First(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Interview not not found"})
		return
	}

	var results []interface{}

	results = append(results, assessment)
	results = append(results, interviewAssessmentMetric)

	c.JSON(http.StatusOK, gin.H{"data": results})
}
