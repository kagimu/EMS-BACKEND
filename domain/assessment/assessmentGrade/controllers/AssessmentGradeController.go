package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	requests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAssessmentGrade(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentGradeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var level levelModels.Level
	if err := config.DB.First(&level, "number = ?", payload.LevelNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No level found with the specified number"})
		return
	}

	saveAssessmentGrade := models.AssessmentGrade{
		Number:         uuid.NewV4().String(),
		Grade:          payload.Grade,
		Description:    payload.Description,
		AggregateRange: payload.AggregateRange,
		Abbreviation:   payload.Abbreviation,
		Condition:      payload.Condition,
		Remarks:        payload.Remarks,
		LevelNumber:    payload.LevelNumber,
		UserNumber:     user.Number,
	}

	if err := config.DB.Create(&saveAssessmentGrade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create assessment grade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Assessment grade created successfully"})
}

func FetchAssessmentGrades(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var assessmentGrades []models.AssessmentGrade
	if err := config.DB.Find(&assessmentGrades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment grades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentGrades})
}

func EditAssessmentGrade(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentGradeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	assessmentGradeNumber := c.Query("assessment_grade_number")
	assessmentGrade := models.AssessmentGrade{}
	if err := config.DB.First(&assessmentGrade, "number = ?", assessmentGradeNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No assessment grade found with the specified number"})
		return
	}

	updateAssessmentGrade := models.AssessmentGrade{
		Grade:          payload.Grade,
		Description:    payload.Description,
		AggregateRange: payload.AggregateRange,
		Abbreviation:   payload.Abbreviation,
		Condition:      payload.Condition,
		Remarks:        payload.Remarks,
	}

	if err := config.DB.Model(&assessmentGrade).Where("number  = ?", assessmentGradeNumber).Updates(updateAssessmentGrade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update assessment grade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentGrade})
}

func DeleteAssessmentGrade(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	assessmentGradeNumber := c.Query("assessment_grade_number")
	var assessmentGrade models.AssessmentGrade

	if err := config.DB.First(&assessmentGrade, assessmentGradeNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assessment grade not found"})
		return
	}

	if err := config.DB.Delete(&assessmentGrade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment grade deleted successfully"})
}
