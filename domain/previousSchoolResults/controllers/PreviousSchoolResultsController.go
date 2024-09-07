package controllers

import (
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/api/responses"
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

	studentNumber := c.PostForm("student_number")
	applicationNumber := c.PostForm("application_number")

	subjects := c.PostFormArray("subject")
	scores := c.PostFormArray("score")

	fileName, err := commons.HandleFileUpload(c, "attachment", "./fileRepo/uploads")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
		return
	}

	for i := range subjects {
		score, _ := strconv.ParseFloat(scores[i], 64)

		previousSchoolResultsNumber := uuid.NewV4().String()

		previousSchoolResults := models.PreviousSchoolResults{
			Number:            previousSchoolResultsNumber,
			StudentNumber:     studentNumber,
			Subject:           subjects[i],
			Score:             score,
			Attachment:        fileName,
			ApplicationNumber: applicationNumber,
			UserNumber:        user.Number,
		}

		if err := config.DB.Create(&previousSchoolResults).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create previous school results"})
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
	if err := config.DB.Preload("AssessmentArea").Preload("Application").Preload("Student").Find(&previousSchoolResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch previous school results"})
		return
	}

	var response []responses.PreviousSchoolResultResponses
	for _, result := range previousSchoolResults {
		response = append(response, responses.PreviousSchoolResultResponses{
			Number:             result.Number,
			StudentName:        result.Application.Student.SecondName + " " + result.Student.OtherName + " " + result.Student.Surname,
			AssessmentAreaName: result.Subject,
			Score:              result.Score,
			Attachment:         result.Attachment,
			ApplicationNumber:  result.ApplicationNumber,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func FetchAllPreviousSchoolResults(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var previousSchoolResults []models.PreviousSchoolResults

	if err := config.DB.Preload("Application").Preload("AssessmentArea").Preload("Student").Find(&previousSchoolResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch previous school results"})
		return
	}

	var response []responses.PreviousSchoolResultResponses
	for _, result := range previousSchoolResults {
		response = append(response, responses.PreviousSchoolResultResponses{
			Number:             result.Number,
			StudentName:        result.Application.Student.SecondName + " " + result.Student.OtherName + " " + result.Student.Surname,
			AssessmentAreaName: result.Subject,
			Score:              result.Score,
			Attachment:         result.Attachment,
			ApplicationNumber:  result.ApplicationNumber,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
