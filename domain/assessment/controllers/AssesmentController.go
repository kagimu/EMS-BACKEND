package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
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
	if err := config.DB.Model(&studentModel.Student{}).Where("number = ?", payload.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	if student.APNStatus == "Inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student ADN is still active"})
		return
	}

	saveAssessment := models.Assessment{
		Number:                 uuid.NewV4().String(),
		AssessmentDate:         payload.AssessmentDate,
		UserNumber:             user.Number,
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

func FetchAssessment(c *gin.Context) {
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

	var assessments []models.Assessment
	if err := config.DB.Model(&models.Assessment{}).Preload("AssessmentMetric").Preload("Student").Preload("AssessmentArea").Preload("Level").Where("student_number = ?", studentNumber).Find(&assessments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessments})
}
