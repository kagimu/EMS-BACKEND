package controllers

import (
	requests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAssessmentGrade(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentGradeRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.SaveAssessmentGrade(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Assessment grade created successfully"})
}

func FetchAssessmentGrades(c *gin.Context) {

	classNumber := c.Query("class_number")
	status := c.Query("status")

	assessmentGradeResponse, err := services.GetAssessmentGrades(classNumber, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assessmentGradeResponse})
}

func EditAssessmentGrade(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("coming soon"))
}

func DeleteAssessmentGrade(c *gin.Context) {

	assessmentGradeNumber := c.Param("assessment_grade_number")
	if assessmentGradeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "provide assessment grade number"})
		return
	}

	err := services.RemoveAssessmentGrade(assessmentGradeNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment grade deleted successfully"})
}
