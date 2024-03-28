package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentMetricModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

	var student studentModel.Student
	if err := config.DB.Model(&studentModel.Student{}).Preload("CampusClass").Where("number = ?", payload.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	var assessmentMetric assessmentMetricModel.AssessmentMetric
	if err := config.DB.First(&assessmentMetric, payload.AssessmentMetricNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assessment metric not found"})
		return
	}

	if student.CampusClassNumber != assessmentMetric.CampusClassNumber {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric and student are not in the same class"})
		return
	}

	saveAssessment := models.Assessment{
		Number:                 uuid.NewV4().String(),
		AssessmentDate:         payload.AssessmentDate,
		UserNumber:             user.Number,
		Comment:                payload.Comment,
		AssessmentMetricNumber: payload.AssessmentMetricNumber,
		StudentNumber:          payload.StudentNumber,
		ActualScore:            payload.ActualScore,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveAssessment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save assessment"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save assessment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": saveAssessment})
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

func FetchAssessments(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var assessments []models.Assessment
	if err := config.DB.Model(&models.Assessment{}).Preload("AssessmentMetric").Preload("Student").Find(&assessments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessments})
}
