package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	levelModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
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
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var assessmentArea assessmentAreaModel.AssessmentArea
	if err := config.DB.Where("number = ?", payload.AssessmentAreaNumber).First(&assessmentArea).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment area not found"})
		return
	}

	var level levelModel.Level
	if err := config.DB.Where("number = ?", payload.LevelNumber).First(&level).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Level not found"})
		return
	}

	var examinationType models.ExaminationType
	if err := config.DB.Where("number = ?", payload.ExaminationTypeNumber).First(&examinationType).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Examination type not found"})
		return
	}

	var student studentModel.Student
	if err := config.DB.Where("number = ?", payload.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	saveAssessment := models.Assessment{
		Number:                uuid.NewV4().String(),
		Subject:               payload.Subject,
		Attendance:            payload.Attendance,
		Observations:          payload.Observations,
		GuidedExercise:        payload.GuidedExercise,
		UnguidedExercise:      payload.UnguidedExercise,
		ExplorativeAssignment: payload.ExplorativeAssignment,
		AssessmentAreaNumber:  payload.AssessmentAreaNumber,
		GroupWork:             payload.GroupWork,
		Project:               payload.Project,
		ExaminationTypeNumber: payload.ExaminationTypeNumber,
		TopicalText:           payload.TopicalText,
		StudentNumber:         payload.StudentNumber,
		LevelNumber:           payload.LevelNumber,
		UserNumber:            user.Number,
	}

	if err := config.DB.Create(&saveAssessment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create assessment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAssessment})
}

func FetchAssessmentByStudentNumber(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	var assessment assessmentAreaModel.AssessmentArea
	if err := config.DB.Where("student_number = ?", studentNumber).Find(&assessment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment fetching failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessment})
}
