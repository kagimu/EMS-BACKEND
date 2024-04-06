package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePreviousSchoolResults(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.PreviousSchoolResultsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input", "err": err.Error()})
		return
	}

	pdfFile, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PDF file is missing"})
		return
	}

	if pdfFile.Size > 20971520 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Pdf file must not exceed 20 mbs!"})
		return
	}

	if pdfFile.Size < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Pdf"})
		return
	}

	filename := pdfFile.Filename

	ext := filepath.Ext(filename)

	newFileName := uuid.NewV4().String() + ext

	path := filepath.Join("/fileRepo/uploads", newFileName)

	err = c.SaveUploadedFile(pdfFile, path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error() + "Failed to save File"})
		return
	}

	for i, score := range payload.Score {
		previousSchoolResultsNumber := uuid.NewV4().String()

		savepreviousSchoolResults := models.PreviousSchoolResults{
			Number:               previousSchoolResultsNumber,
			StudentNumber:        payload.StudentNumber,
			AssessmentAreaNumber: payload.AssessmentAreaNumber[i],
			Score:                score,
			Attachment:           newFileName,
			ApplicationNumber:    payload.ApplicationNumber,
			UserNumber:           user.Number,
		}

		if err := config.DB.Create(&savepreviousSchoolResults).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview grade"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Previous school results created successfully"})
}

func FetchPreviousSchoolResultsByStudent(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	var student studentModels.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch student"})
		return
	}

	var previousSchoolResults []models.PreviousSchoolResults
	if err := config.DB.Preload("AssessmentArea").Preload("Application").Find(&previousSchoolResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch previous school results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": previousSchoolResults})
}

func FetchAllPreviousSchoolResults(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var previousSchoolResults []models.PreviousSchoolResults

	if err := config.DB.Preload("CampusClass").Find(&previousSchoolResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch previous school results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": previousSchoolResults})
}
