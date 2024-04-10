package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAssessmentArea(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentAreaRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	//validate for composite unique curriculum_class_number and subject number

	var assessmentArea models.AssessmentArea
	if err := config.DB.Where("curriculum_class_number = ? AND subject_number = ?", payload.CurriculumClassNumber, payload.SubjectNumber).First(&assessmentArea).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment area already exists for the specified curriculum class and subject"})
		return

	}

	saveAssessmentArea := models.AssessmentArea{
		Number:                uuid.NewV4().String(),
		CurriculumClassNumber: payload.CurriculumClassNumber,
		Name:                  payload.Name,
		UserNumber:            user.Number,
		SubjectNumber:         payload.SubjectNumber,
	}

	if err := config.DB.Create(&saveAssessmentArea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create assessment area"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment area created successfully"})
}

func FetchAssessmentAreas(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var curriculumClassNumber = c.Query("curriculum_class_number")
	var assessmentAreas []models.AssessmentArea

	if err := config.DB.Where("curriculum_class_number = ?", curriculumClassNumber).Find(&assessmentAreas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentAreas})
}

func EditAssessmentArea(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var assessmentArea models.AssessmentArea
	assessmentAreaNumber := c.Param("assessment_area_number")
	if assessmentAreaNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment area number is required"})
		return
	}

	if err := config.DB.First(&assessmentArea, "number = ?", assessmentAreaNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No assessment area found with the specified number"})
		return
	}

	var payload requests.AssessmentAreaRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	updateAssessmentArea := models.AssessmentArea{
		Name: payload.Name,
	}

	if err := config.DB.Model(&assessmentArea).Where("number = ?").Updates(updateAssessmentArea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update assessment area"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment area updated successfully"})
}

//prepare Scope for assessment area
