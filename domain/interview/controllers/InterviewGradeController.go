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

func CreateInterviewGrade(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.InterviewGradeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	studentNumber := c.Param("number")
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Interview not found"})
		return
	}

	userNumber := user.Number
	for i, questionNumber := range payload.QuestionNumbers {
		score := payload.Scores[i]
		interviewGradeNumber := uuid.NewV4().String()

		saveInterviewGrade := models.InterviewGrade{
			Number:          interviewGradeNumber,
			QuestionNumber:  questionNumber,
			InterviewNumber: payload.InterviewNumber,
			StudentNumber:   studentNumber,
			Score:           score,
			UserNumber:      userNumber,
		}

		if err := config.DB.Create(&saveInterviewGrade).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview grade"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview grade created successfully"})
}

func GetStudentGrades(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	studentNumber := c.Param("number")
	var student studentModel.Student

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No student found with the specified number"})
		return
	}

	var results []models.InterviewGrade
	query := fmt.Sprintf("SELECT * FROM interview_grades WHERE student_number = '%s'", studentNumber)

	if err := config.DB.Raw(query).Preload("InterviewQuestion").Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch interview results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
