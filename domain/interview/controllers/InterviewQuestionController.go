package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateInterviewQuestion(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
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
	var interview models.Interview
	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Interview not found"})
		return
	}

	tx := config.DB.Begin()

	interviewQuestionNumber := uuid.NewV4().String()
	saveInterviewQuestion := models.InterviewQuestion{
		Number:          interviewQuestionNumber,
		InterviewNumber: interviewNumber,
		Question:        payload.Question,
		Score:           payload.Score,
		CorrectAnswer:   payload.CorrectAnswer,
		UserNumber:      user.Number,
	}

	if err := tx.Create(&saveInterviewQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create interview question"})
		return
	}

	for _, possibleAnswer := range payload.PossibleAnswers {
		saveAnswers := models.Answer{
			Number:         uuid.NewV4().String(),
			QuestionNumber: interviewQuestionNumber,
			Body:           possibleAnswer,
		}

		if err := tx.Create(&saveAnswers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create possible answers"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveInterviewQuestion})
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Interview question not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch interview question", "err": err.Error()})
			return
		}
	}

	var answers []models.Answer
	if err := config.DB.Where("question_number = ?", interviewQuestion.Number).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch answers"})
		return
	}

	var results []responses.InterviewQuestionResponse

	for _, answer := range answers {
		result := responses.InterviewQuestionResponse{
			InterviewNumber: answer.Number,
			Question:        interviewQuestion.InterviewNumber,
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func EditInterviewQuestion(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var payload requests.InterviewQuestionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var interviewQuestion models.InterviewQuestion
	interviewQuestionNumber := c.Param("number")
	if err := config.DB.First(&interviewQuestion, "number = ?", interviewQuestionNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No interview question with the specified number",
		})
		return
	}

	updateInterviewQuestion := models.InterviewQuestion{
		Number:          uuid.NewV4().String(),
		InterviewNumber: payload.InterviewNumber,
		Question:        payload.Question,
		CorrectAnswer:   payload.CorrectAnswer,
	}

	if err := config.DB.Model(&interviewQuestion).Where("number = ?", interviewQuestionNumber).Updates(updateInterviewQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update interview question",
		})
		return
	}

	for _, possibleAnswer := range payload.PossibleAnswers {
		saveAnswers := models.Answer{
			QuestionNumber: interviewQuestionNumber,
			Body:           possibleAnswer,
		}

		if err := config.DB.Create(&saveAnswers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create possible answers"})
			return
		}

		if err := config.DB.Model(&interviewQuestion).Where("number = ?", interviewQuestionNumber).Updates(updateInterviewQuestion).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to update interview question",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": interviewQuestion})
}

func FetchAllInterviewQuestions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var interviewQuestions []models.InterviewQuestion

	if err := config.DB.Find(&interviewQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch interview questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": interviewQuestions})
}
