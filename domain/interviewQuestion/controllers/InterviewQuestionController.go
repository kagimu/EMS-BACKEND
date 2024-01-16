package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interviewQuestion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interviewQuestion/models"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateInterviewQuestion(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	var payload requests.InterviewQuestionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	interviewNumber := payload.InterviewNumber
	var interview interviewModel.Interview
	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Interview not found"})
		return
	}

	for _, possibleAnswer := range payload.PossibleAnswers {
		saveInterviewQuestion := models.InterviewQuestion{
			Number:          uuid.NewV4().String(),
			InterviewNumber: interviewNumber,
			Question:        payload.Question,
			CorrectAnswer:   payload.CorrectAnswer,
			PossibleAnswers: possibleAnswer,
		}

		if err := config.DB.Create(&saveInterviewQuestion).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create interview question"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview question created"})
}

func FetchInterviewQuestion(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	interviewQuestionNumber := c.Param("number")
	var interviewQuestion models.InterviewQuestion
	if err := config.DB.Where("number = ?", interviewQuestionNumber).First(&interviewQuestion).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewQuestion})
}

func EditInterviewQuestion(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var payload requests.InterviewQuestionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to process input",
		})
		return
	}

	var interviewQuestion models.InterviewQuestion
	interviewQuestionNumber := c.Param("number")
	if err := config.DB.First(&interviewQuestion, "apn = ?", interviewQuestionNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No interview question with the specified number",
		})
		return
	}

	for _, possibleAnswer := range payload.PossibleAnswers {
		updateInterviewQuestion := models.InterviewQuestion{
			Number:          uuid.NewV4().String(),
			InterviewNumber: payload.InterviewNumber,
			Question:        payload.Question,
			CorrectAnswer:   payload.CorrectAnswer,
			PossibleAnswers: possibleAnswer,
		}

		if err := config.DB.Model(&interviewQuestion).Where("apn = ?", interviewQuestionNumber).Updates(updateInterviewQuestion).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to update interview question",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": interviewQuestion})
}
