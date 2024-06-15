package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	var student studentModel.Student
	if err := config.DB.Model(&studentModel.Student{}).Preload("CampusClass").Where("number = ?", payload.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	var assessmentMetric models.AssessmentMetric
	if err := config.DB.Where("number = ?", payload.AssessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assessment metric not found"})
		return
	}

	if payload.ActualScore < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Actual score should not be less than 0"})
		return
	}

	if payload.ActualScore > assessmentMetric.MaxScore {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Actual score should not be more than "})
		return

	}

	//check if there is existing enrollment for the student

	var enrollment enrollmentModel.Enrollment
	if err := config.DB.Where("student_number = ?", student.Number).First(&enrollment).Error; err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if student.CampusClassNumber != assessmentMetric.CampusClassNumber {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric and student are not in the same class"})
				return
			}
		}

	}

	//find if assessment with student number and assessment metric  already exists

	var existingAssessment models.Assessment
	if err := config.DB.Where("student_number = ? AND assessment_metric_number = ?", payload.StudentNumber, payload.AssessmentMetricNumber).First(&existingAssessment).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment already exists"})
			return
		}
	}

	tx := config.DB.Begin()

	if existingAssessment.Number != "" {
		existingAssessment.Comment = payload.Comment
		existingAssessment.ActualScore = payload.ActualScore

		if err := tx.Save(&existingAssessment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save assessment"})
			return

		}

	} else {

		saveAssessment := models.Assessment{
			Number:                 uuid.NewV4().String(),
			UserNumber:             user.Number,
			Comment:                payload.Comment,
			AssessmentMetricNumber: payload.AssessmentMetricNumber,
			StudentNumber:          payload.StudentNumber,
			ActualScore:            payload.ActualScore,
		}

		if err := tx.Create(&saveAssessment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save assessment"})
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save assessment"})
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
