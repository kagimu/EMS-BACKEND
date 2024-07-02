package controllers

import (
	"errors"
	"fmt"
	assessmentMetricRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	enrollmentResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	"gorm.io/gorm"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
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

	var interviews []models.Interview
	if err := config.DB.Find(&interviews).Error; err != nil {
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
	assessmentMetricNumber := c.Query("assessment_metric_number")

	if assessmentMetricNumber == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("Choose the assessment schedule"))
		return
	}

	assessmentMetric, err := assessmentMetricRepositories.FindAssessmentMetricByNumber(assessmentMetricNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("assessment metric not found"+err.Error()))
		return
	}

	if assessmentMetric.Number == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("assessment metric not found"))
		return
	}

	var interview models.Interview
	if err := config.DB.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Interview not found",
		})
		return
	}

	var interviewStudents []models.InterviewStudent
	if err := config.DB.Preload("Student").Where("interview_number = ?", interviewNumber).Find(&interviewStudents).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to fetch students",
		})
		return
	}

	//assessmentMetricNumber := c.Query("assessment_metric_number")
	var stds []enrollmentResponse.EnrolledStudent

	for _, interviewStudent := range interviewStudents {
		student := interviewStudent.Student

		if student.CampusClassNumber != assessmentMetric.CampusClassNumber {
			continue
		}

		if student.Status != constants.StatusInterviewScheduled {
			continue
		}

		var assessment assessmentModel.Assessment

		if err := config.DB.Where("student_number = ? AND assessment_metric_number = ?", interviewStudent.StudentNumber, assessmentMetricNumber).First(&assessment).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessments " + err.Error()})
				return
			}
		}

		//var campusClass classModels.CampusClass
		//if err := config.DB.Where("number = ?", enrollment.CampusClassNumber).First(&campusClass).Error; err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"message": "Campus class not found"})
		//	return
		//}

		stds = append(stds, enrollmentResponse.EnrolledStudent{
			Number:            student.Number,
			StudentSurname:    student.Surname,
			StudentOtherName:  student.OtherName,
			StudentSecondName: student.SecondName,
			Comment:           assessment.Comment,
			Score:             assessment.Score,
			ActualScore:       assessment.ActualScore,
			Status:            student.Status,
		})
	}

	if len(stds) == 0 {
		c.JSON(http.StatusBadRequest, commons.Feedback("No student on this interview for the selected assessment in class "+assessmentMetric.CampusClassNumber))
	} else {
		c.JSON(http.StatusOK, gin.H{"data": stds})
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
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
