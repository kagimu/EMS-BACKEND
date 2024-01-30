package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
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

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	saveInterview := models.Interview{
		Number:        uuid.NewV4().String(),
		Message:       payload.Message,
		Venue:         payload.Venue,
		Name:          payload.Name,
		InterviewDate: payload.InterviewDate,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveInterview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview"})
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
		"data": interview,
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviews})
}

func FetchStudentsByInterviewNumber(c *gin.Context) {

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

	var students []models.InterviewStudent
	if err := config.DB.Where("interview_number = ?", interviewNumber).Find(&students).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to fetch students",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": students})

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

	c.JSON(http.StatusOK, gin.H{"message": "Student added to the interview"})
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
