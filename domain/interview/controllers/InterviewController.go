package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
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

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	var student studentModel.Student
	studentNumber := payload.StudentNumber

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Student not found",
		})
		return
	}

	saveInterview := models.Interview{
		Number:        uuid.NewV4().String(),
		Message:       payload.Message,
		Venue:         payload.Venue,
		InterviewDate: payload.InterviewDate,
		StudentNumber: studentNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveInterview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create interview",
		})
		return
	}

	applicationNumber := student.ApplicationNumber
	application, err := applicationController.FetchApplicationByNumber(c, applicationNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred while fetching application"})
		return
	}

	updateApplicationStatus := applicationModel.Application{
		Status: "INTERVIEW_SCHEDULED",
	}

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplicationStatus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "application updating failed",
		})
		return
	}

	//sendEmail := emailSender.Info{Email: application.Parent.Email, Token: "", MailType: "Wholesome: Scheduled Interview", UserName: application.Parent.SecondName, CurrentYear: time.Now().Year(), Type: "scheduled-interview", Data: saveInterview.InterviewDate}
	//if err := sendEmail.SendEmail(); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "failed to send email to the parent"})
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Interview created successfully", "data": saveInterview})
}

func EditInterview(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	type InterviewPayload struct {
		Message       string `json:"message"`
		Venue         string `json:"venue"`
		InterviewDate string `json:"interview_date"`
	}

	var payload InterviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var interview models.Interview
	interviewNumber := c.Param("interview_number")

	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview not found"})
		return
	}

	saveInterview := models.Interview{
		Message:       payload.Message,
		Venue:         payload.Venue,
		InterviewDate: payload.InterviewDate,
	}

	if err := config.DB.Model(&interview).Where("number = ?", interviewNumber).Updates(saveInterview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update interview",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": interview})
}

func FetchAnInterview(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var interview models.Interview
	interviewNumber := c.Param("interview_number")

	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Interview not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    interview,
		"message": "Interview fetched successfully",
	})
}

func FetchAllInterviews(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	paginationParams := helpers.ParsePaginationParams(c)

	var interviews []models.Interview
	query := config.DB.Model(&models.Interview{})

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
		})
		return
	}

	if err := query.Find(&interviews).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "interview fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    interviews,
		"message": "interviews fetched successfully",
	})
}
