package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/service"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateInterview(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.InterviewRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	err = service.CreateInterview(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Interview created successfully"))
}

func EditInterview(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.EditInterviewRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	err = service.UpdateInterview(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Interview updated successfully"))
}

func FetchAnInterview(c *gin.Context) {

	interviewNumber := c.Query("interview_number")
	if interviewNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide interview number"})
		return
	}

	interviewDetails, err := service.GetInterviewDetails(interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": interviewDetails,
	})
}

func FetchAllInterviews(c *gin.Context) {

	campusClassNumber := c.Query("campus_class_number")
	if campusClassNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose campus class number"})
		return
	}

	allInterviews, err := service.AllInterviews(campusClassNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": allInterviews,
	})
}

func FetchStudentsByInterviewNumber(c *gin.Context) {
	interviewNumber := c.Query("interview_number")
	assessmentMetricNumber := c.Query("assessment_metric_number")

	students, err := service.InterviewStudents(assessmentMetricNumber, interviewNumber)

	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responses.AssessmentAndStudentResponse{
		Headers:  []responses.SingleAssessment{},
		Students: students,
	}})

}

func AddStudentToInterview(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var request requests.InterviewStudentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	err := request.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var InterviewStudent models.InterviewStudent
	if err := config.DB.Where("interview_number = ? AND student_number = ?", request.InterviewNumber, request.StudentNumber).First(&InterviewStudent).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student already added to interview"})
		return
	}

	var interview models.Interview
	if err := config.DB.Where("number = ?", request.InterviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview not found"})
		return
	}

	var student studentModel.Student
	if err := config.DB.Where("number = ?", request.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	saveInterviewStudent := models.InterviewStudent{
		Number:          uuid.NewV4().String(),
		InterviewNumber: request.InterviewNumber,
		StudentNumber:   request.StudentNumber,
	}

	if err := config.DB.Create(&saveInterviewStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add student to interview"})
		return
	}

	// Update the application status
	applicationNumber := student.ApplicationNumber
	application, err := applicationController.FetchApplicationByNumber(c, applicationNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred while fetching application"})
		return
	}

	updateApplicationStatus := applicationModel.Application{
		Status: constants.StatusInterviewScheduled,
	}

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplicationStatus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student added to the interview"})
}

func RemoveStudentFromInterview(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	interviewNumber := c.Query("interview_number")
	studentNumber := c.Query("student_number")

	var InterviewStudent models.InterviewStudent
	if err := config.DB.Where("interview_number = ? AND student_number = ?", interviewNumber, studentNumber).First(&InterviewStudent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found in interview"})
		return
	}

	if err := config.DB.Delete(&InterviewStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove student from interview"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student removed from interview"})
}

func FetchQuestionsByInterviewNumber(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	interviewNumber := c.Param("interview_number")

	var interview models.Interview
	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Interview not found",
		})
		return
	}

	var questions []models.InterviewQuestion
	if err := config.DB.Model(&models.InterviewQuestion{}).Where("interview_number = ?", interviewNumber).Find(&questions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to fetch questions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": questions})
}
