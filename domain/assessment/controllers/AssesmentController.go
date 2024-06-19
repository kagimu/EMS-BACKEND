package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAssessment(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	err = services.SaveAssessment(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Assessment saved successfully"})
}

func FetchAssessmentByStudent(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	var student studentModel.Student
	if err := config.DB.Model(&studentModel.Student{}).Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	assessmentMetricNumber := c.Query("assessment_metric_number")
	var assessmentMetric models.AssessmentMetric
	if err := config.DB.Model(&models.AssessmentMetric{}).Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric not found"})
		return
	}

	var assessments []models.Assessment
	if err := config.DB.Model(&models.Assessment{}).Preload("AssessmentMetric").Preload("Student").Where("student_number = ? AND assessment_metric_number = ?", studentNumber, assessmentMetricNumber).Find(&assessments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessments})
}
