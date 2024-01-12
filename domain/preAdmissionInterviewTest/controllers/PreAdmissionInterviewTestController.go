package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	gradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
)

func CreatePreAdmissionInterviewTest(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not exist"})
		return
	}

	var payload requests.PreAdmissionInterviewTestRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to validate input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("failed to validate input %s:", err.Error()),
		})
		return
	}

	var student studentModel.Student
	studentNumber := c.Query("student_number")

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Student not found",
		})
		return
	}

	var grade gradeModel.Grade
	gradeNumber := c.Query("grade_number")

	if err := config.DB.First(&grade, "number = ?", gradeNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Grade not found",
		})
		return
	}

	var class classModel.Class
	classNumber := c.Query("class_number")

	if err := config.DB.First(&class, "number = ?", classNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Class not found",
		})
		return
	}

	var subject gradeModel.Grade
	subjectNumber := c.Query("subject_number")

	if err := config.DB.First(&subject, "number = ?", subjectNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Subject not found",
		})
		return
	}

	savePreAdmissionInterviewTest := models.PreAdmissionInterviewTest{
		UserNumber:    user.Number,
		SubjectNumber: payload.SubjectNumber,
		Marks:         payload.Marks,
		GradeNumber:   payload.GradeNumber,
		ClassNumber:   classNumber,
		Location:      payload.Location,
		InterviewDate: payload.InterviewDate,
	}

	if err := config.DB.Create(&savePreAdmissionInterviewTest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create Pre Admission Interview Test",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pre Admission Interview Test created successfully"})
}
