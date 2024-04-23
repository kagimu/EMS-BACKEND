package controllers

import (
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/api/requests"
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

	var payload requests.PreviousSchoolResultsRequest
	payload.StudentNumber = c.PostForm("student_number")
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.ApplicationNumber = c.PostForm("application_number")

	var assessmentAreasAndScores []requests.AssessmentAndScore
	areas := c.PostFormArray("assessment_area_number")
	scores := c.PostFormArray("score")

	for i := range areas {
		score, _ := strconv.ParseFloat(scores[i], 64)
		assessmentAreasAndScores = append(assessmentAreasAndScores, requests.AssessmentAndScore{
			AssessmentAreaNumber: areas[i],
			Score:                score,
		})
	}

	payload.AssessmentAndScore = assessmentAreasAndScores

	fileName, err := commons.HandleFileUpload(c, "attachment", "./fileRepo/uploads")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
		return
	}

	for _, assessment := range payload.AssessmentAndScore {
		previousSchoolResultsNumber := uuid.NewV4().String()

		savepreviousSchoolResults := models.PreviousSchoolResults{
			Number:               previousSchoolResultsNumber,
			StudentNumber:        payload.StudentNumber,
			AssessmentAreaNumber: assessment.AssessmentAreaNumber,
			Score:                assessment.Score,
			Attachment:           fileName,
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

	var response []responses.PreviousSchoolResultResponses
	for _, result := range previousSchoolResults {
		response = append(response, responses.PreviousSchoolResultResponses{
			Number:             result.Number,
			StudentName:        result.Application.Student.SecondName + " " + result.Student.OtherName + " " + result.Student.Surname,
			AssessmentAreaName: result.AssessmentArea.Name,
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

	if err := config.DB.Preload("CampusClass").Find(&previousSchoolResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch previous school results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": previousSchoolResults})
}
