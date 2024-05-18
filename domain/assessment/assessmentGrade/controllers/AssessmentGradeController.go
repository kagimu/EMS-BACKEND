package controllers

import (
	"net/http"
	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"

	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	requests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
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

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var fileName string
	_, err = c.FormFile("illustration")
	if err == nil {
		file, err := commons.HandleFileUpload(c, "illustration", "/fileRepo/illustrations")
		fileName = file
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
			return
		}
	}

	var campusClass classModel.CampusClass
	if err := config.DB.First(&campusClass, "number = ?", payload.CampusClassNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No campus class found with the specified number"})
		return
	}

	var duplicate models.AssessmentGrade
	if err := config.DB.Where("name = ? AND campus_class_number = ? AND grading_standard_number = ? ", payload.Name, payload.CampusClassNumber).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate aggregate scale name on this class"})
		return
	}

	if err := config.DB.Where("code = ? AND campus_class_number = ? ", payload.Code, payload.CampusClassNumber).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This code has already been used on this class"})
		return
	}

	saveAssessmentGrade := models.AssessmentGrade{
		Code:              payload.Code,
		CampusNumber:      payload.CampusNumber,
		Number:            uuid.NewV4().String(),
		Name:              payload.Name,
		Description:       payload.Description,
		MinimumPoints:     payload.MinimumPoints,
		MaximumPoints:     payload.MaximumPoints,
		Abbreviation:      payload.Abbreviation,
		Remarks:           payload.Remarks,
		Illustration:      fileName,
		CampusClassNumber: payload.CampusClassNumber,
		UserNumber:        user.Number,
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

	campusClassNumber := c.Query("campus_class_number")
	campusNumber := c.Query("campus_number")

	var assessmentGrades []models.AssessmentGrade
	if err := config.DB.Scopes(filters.ByCampusClassNumber(campusClassNumber), filters.ByCampusNumber(campusNumber)).Preload("CampusClass").Find(&assessmentGrades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment grades"})
		return
	}

	baseURL := os.Getenv("BASE_URL")

	var response []responses.AssessmentGradeResponse
	for _, assessmentGrade := range assessmentGrades {
		response = append(response, responses.AssessmentGradeResponse{
			Code:             assessmentGrade.Code,
			Number:           assessmentGrade.Number,
			Name:             assessmentGrade.Name,
			Description:      assessmentGrade.Description,
			MinimumPoints:    assessmentGrade.MinimumPoints,
			MaximumPoints:    assessmentGrade.MaximumPoints,
			Abbreviation:     assessmentGrade.Abbreviation,
			Remarks:          assessmentGrade.Remarks,
			CampusClassName:  assessmentGrade.CampusClass.ClassName,
			IllustrationPath: baseURL + assessmentGrade.Illustration,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
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

	assessmentGradeNumber := c.Param("assessment_grade_number")
	var assessmentGrade models.AssessmentGrade
	if err := config.DB.First(&assessmentGrade).Where("number = ?", assessmentGradeNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No assessment grade found with the specified number"})
		return
	}

	var fileName string
	_, err = c.FormFile("illustration")
	if err == nil {
		file, err := commons.HandleFileUpload(c, "illustration", "/fileRepo/illustrations")
		fileName = file
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
			return
		}
	}

	updateAssessmentGrade := models.AssessmentGrade{
		Code:              payload.Code,
		Description:       payload.Description,
		Name:              payload.Name,
		MaximumPoints:     payload.MaximumPoints,
		CampusClassNumber: payload.CampusClassNumber,
		CampusNumber:      payload.CampusNumber,
		Illustration:      fileName,
		MinimumPoints:     payload.MinimumPoints,
		Abbreviation:      payload.Abbreviation,
		Remarks:           payload.Remarks,
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
