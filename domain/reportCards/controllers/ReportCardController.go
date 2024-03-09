package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportCards/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportCards/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateReportCard(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ReportCardRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	studentNumber := payload.StudentNumber
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	assessmentNumber := payload.AssesmentNumber
	var assessment assessmentModel.Assessment
	if err := config.DB.Where("number = ?", assessmentNumber).First(&assessment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment not found"})
		return
	}

	saveReport := models.ReportCard{
		Number:           uuid.NewV4().String(),
		StudentNumber:    studentNumber,
		AssessmentNumber: assessmentNumber,
		UserNumber:       user.Number,
	}

	if err := config.DB.Create(&saveReport).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create student report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report card created successfully"})
}

func FetchAllPromotions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var reports []models.ReportCard
	if err := config.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}
