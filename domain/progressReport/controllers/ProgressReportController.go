package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/progressReport/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/progressReport/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateProgressReport(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ProgressReportRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
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

	saveReport := models.ProgressReport{
		Number:            uuid.NewV4().String(),
		StudentNumber:     studentNumber,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		Comment:           payload.Comment,
		ClassTeacher:      payload.ClassTeacher,
		SubjectTeacher:    payload.SubjectTeacher,
		UserNumber:        user.Number,
	}

	if err := config.DB.Create(&saveReport).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create student report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report card created successfully"})
}

func FetchAllProgressReports(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var reports []models.ProgressReport
	if err := config.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report cards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reports})
}
