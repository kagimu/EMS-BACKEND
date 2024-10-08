package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/requests"
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
	var classNumber = c.Query("class_number")
	var curriculumNumber = c.Query("curriculum_number")
	var levelNumber = c.Query("level_number")
	var curriculumLevelNumber = c.Query("curriculum_level_number")

	if curriculumClassNumber == "" && (classNumber == "" || curriculumNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose class form the curriculum"})
		return
	}

	/*var assessmentAreas []models.AssessmentArea

	if classNumber != "" && curriculumNumber != "" {

		if err := config.DB.
			Joins("JOIN curriculum_classes ON curriculum_classes.number = assessment_areas.curriculum_class_number").
			Joins("JOIN curriculum_levels ON curriculum_levels.number = curriculum_classes.curriculum_level_number").
			Where("curriculum_classes.class_number = ? AND curriculum_levels.curriculum_number = ?", classNumber, curriculumNumber).
			Distinct("assessment_areas.id", "assessment_areas.curriculum_class_number", "assessment_areas.number", "assessment_areas.name", "assessment_areas.user_number", "assessment_areas.subject_number").
			Find(&assessmentAreas).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"data": make([]models.AssessmentArea, 0)})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessment areas with class and curriculum" + err.Error()})
			return
		}

	} else {
		if err := config.DB.Where("curriculum_class_number = ?", curriculumClassNumber).Find(&assessmentAreas).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"data": make([]models.AssessmentArea, 0)})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessment areas"})
			return
		}
	}*/

	all, err := services.FindAllAssessmentAreas(curriculumClassNumber, classNumber, curriculumNumber, levelNumber, curriculumLevelNumber, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": all})
}

func EditAssessmentArea(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateAssessmentAreaRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	err = services.UpdateAssessmentArea(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to update assessment area: %s", err.Error())})
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

	/*var campusClass campusClassModels.CampusClass
	if err := config.DB.Preload("CurriculumClass").Where("number = ?", campusClassNumber).First(&campusClass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus class"})
		return
	}

	var assessmentAreas []models.AssessmentArea
	if err := config.DB.Where("curriculum_class_number = ?", campusClass.CurriculumClass.Number).Find(&assessmentAreas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
		return
	}*/

	all, err := services.FindAllAssessmentAreas(campusClassNumber, "", "", "", "", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": all})
}

//prepare Scope for assessment area
