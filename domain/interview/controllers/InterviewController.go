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
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repository"
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

	tx := config.DB.Begin()
	var repo = repository.NewInterviewRepoImpl(tx)

	interviewNumber := uuid.NewV4().String()
	saveInterview := models.Interview{
		Number:            interviewNumber,
		Message:           payload.Message,
		Venue:             payload.Venue,
		Name:              payload.Name,
		InterviewDate:     payload.InterviewDate,
		UserNumber:        user.Number,
		CampusClassNumber: payload.CampusClassNumber,
	}

	if err := repo.CreateInterview(&saveInterview); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create interview: " + err.Error()})
		return
	}

	var interview models.Interview
	if err := tx.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview schedule not found"})
		return
	}

	for _, assessmentMetricNumber := range payload.AssessmentMetrics {

		var assessmentMetric assessmentModel.AssessmentMetric
		if err := tx.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
			continue
		}

		var interviewAssessmentMetric models.InterviewAssessmentMetric
		if err := tx.Where("interview_number = ? AND assessment_metric_number = ?", interviewNumber, assessmentMetricNumber).First(&interviewAssessmentMetric).Error; err == nil {
			continue
		}

		interviewAssessmentMetric = models.InterviewAssessmentMetric{
			Number:                 uuid.NewV4().String(),
			InterviewNumber:        interviewNumber,
			AssessmentMetricNumber: assessmentMetricNumber,
			UserNumber:             user.Number,
		}

		if err := tx.Create(&interviewAssessmentMetric).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview assessment metric"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit transaction"})
		return

	}
	//sendEmail := emailSender.Info{Email: application.Parent.Email, Token: "", MailType: "Wholesome: Scheduled Interview", UserName: application.Parent.SecondName, CurrentYear: time.Now().Year(), Type: "scheduled-interview", Data: saveInterview.InterviewDate}
	//if err := sendEmail.SendEmail(); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"message": "failed to send email to the parent"})
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Interview created successfully", "data": saveInterview})
}

func EditInterview(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var repo = repository.NewInterviewRepoImpl(config.DB)

	type InterviewPayload struct {
		Message           string `json:"message"`
		Venue             string `json:"venue"`
		InterviewDate     string `json:"interview_date"`
		CampusClassNumber string `json:"campus_class_number"`
	}

	var payload InterviewPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	interviewNumber := c.Param("interview_number")

	interview, err := repo.FindInterviewByNumber(interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview not found"})
		return
	}

	interview.Message = payload.Message
	interview.Venue = payload.Venue
	interview.CampusClassNumber = payload.CampusClassNumber
	interview.InterviewDate = payload.InterviewDate
	interview.UserNumber = user.Number

	if err := repo.UpdateInterview(interview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update interview: " + err.Error()})
		return

	}
	c.JSON(http.StatusOK, commons.Feedback("Interview updated successfully"))
}

func FetchAnInterview(c *gin.Context) {
	var repo = repository.NewInterviewRepoImpl(config.DB)
	interviewNumber := c.Param("interview_number")

	interview, err := repo.FindInterviewByNumber(interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": interview,
	})
}

func FetchAllInterviews(c *gin.Context) {

	var interviews []models.Interview
	if err := config.DB.Find(&interviews).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interviews: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviews})
}

func FetchStudentsByInterviewNumber(c *gin.Context) {
	var repo = repository.NewInterviewRepoImpl(config.DB)
	interviewNumber := c.Param("interview_number")

	_, err := repo.FindInterviewByNumber(interviewNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Interview not found",
		})
		return
	}

	// invoke the repository to findInterview by number

	interviewStudents, err := assessmentMetricRepositories.FindInterviewStudentsByInterViewNumber(interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to fetch students",
		})
		return
	}

	assessmentMetrics, err := assessmentMetricRepositories.FindAssessmentMetricsByInterviewNumber(interviewNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("assessment metric not found"+err.Error()))
		return
	}

	type singleAssessment struct {
		AssessmentScheduleNumber string `json:"assessment_schedule_number"`
		AssessmentScheduleName   string `json:"assessment_schedule_name"`
	}

	type singleStudent struct {
		enrollmentResponse.EnrolledStudent
		Assessments []assessmentModel.Assessment
	}
	type Response struct {
		Headers  []singleAssessment `json:"headers"`
		Students []singleStudent    `json:"students"`
	}

	var studentList = make([]singleStudent, 0)

	var assesssment_metric_number = c.Query("assessment_metric_number")

	for _, interviewStudent := range interviewStudents {
		student := interviewStudent.Student
		var assessments []assessmentModel.Assessment

		for _, assessmentMetric := range assessmentMetrics {

			assmts, err := assessmentMetricRepositories.FindAssessmentByStudentNumberAndAssessmentMetricNumber(interviewStudent.StudentNumber, assessmentMetric.Number)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessments " + err.Error()})
					return
				}
			}

			assessments = append(assessments, assmts...)

		}

		var assessment assessmentModel.Assessment
		if assesssment_metric_number != "" {
			assessment, err = assessmentMetricRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, assesssment_metric_number)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessments " + err.Error()})
					return
				}
			}
		}
		studentList = append(studentList, singleStudent{
			EnrolledStudent: enrollmentResponse.EnrolledStudent{
				Number:            student.Number,
				StudentOtherName:  student.OtherName,
				StudentSurname:    student.Surname,
				StudentSecondName: student.SecondName,
				Score:             assessment.Score,
				ActualScore:       assessment.ActualScore,
				Comment:           assessment.Comment,
			},
			Assessments: assessments,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": Response{
		Headers:  []singleAssessment{},
		Students: studentList,
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
