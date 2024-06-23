package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/services"
	"github.com/gin-gonic/gin"
	"net/http"
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

	err = services.SaveAssessmentArea(payload, user.Number)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to save assessment area: %s", err.Error())})
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

func FetchAssessmentAreasByCampusClassNumber(c *gin.Context) {
	campusClassNumber := c.Param("campus_class_number")
	if campusClassNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus class number is required"})
		return
	}

	var campusClass campusClassModels.CampusClass
	if err := config.DB.Preload("CurriculumClass").Where("number = ?", campusClassNumber).First(&campusClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus class"})
		return
	}

	var assessmentAreas []models.AssessmentArea
	if err := config.DB.Where("curriculum_class_number = ?", campusClass.CurriculumClass.Number).Find(&assessmentAreas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentAreas})
}

//prepare Scope for assessment area
