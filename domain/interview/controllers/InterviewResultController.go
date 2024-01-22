package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateInterviewResults(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User does not exist"})
		return
	}

	var payload requests.InterviewResultRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process the input"})
		return
	}

	studentNumber := c.Param("student_number")
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	for i, questionNumber := range payload.QuestionNumber {
		score := payload.Score[i]
		interviewResultNumber := uuid.NewV4().String()

		saveInterviewResult := models.InterviewResult{
			Number:         interviewResultNumber,
			QuestionNumber: questionNumber,
			Score:          score,
			StudentNumber:  studentNumber,
			UserNumber:     user.Number,
		}

		if err := config.DB.Create(&saveInterviewResult).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview results"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Interview result created successfully"})
}

func FetchInterviewResultByStudentNumber(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User does not exist"})
		return
	}

	studentNumber := c.Param("student_number")
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	var results []models.InterviewResult
	query := fmt.Sprintf("SELECT * FROM interview_results WHERE student_number = '%s'", studentNumber)

	if err := config.DB.Raw(query).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to fetch interview results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
