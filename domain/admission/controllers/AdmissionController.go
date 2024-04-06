package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"net/http"

	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAdmission(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	var payload requests.AdmissionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	for _, studentPayload := range payload.StudentNumber {
		var student studentModels.Student
		if err := config.DB.First(&student, "number = ?", studentPayload).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
			return
		}

		var campusClass campusClassModels.CampusClass
		if err := config.DB.First(&campusClass, "number = ?", payload.CampusClassNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus class not found"})
			return
		}

		var campus campusModels.Campus
		if err := config.DB.First(&campus, "number = ?", payload.CampusNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Campus not found"})
			return
		}

		var stream campusClassModels.Stream
		if err := config.DB.First(&stream, "number = ?", payload.StreamNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Stream not found"})
			return
		}

		var studyPeriod studyPeriodModel.StudyPeriod
		if err := config.DB.First(&studyPeriod, "number = ?", payload.StudyPeriodNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Study period not found"})
			return
		}

		var interview interviewModel.Interview
		if err := config.DB.First(&interview, "number = ?", payload.InterviewNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Interview not found"})
			return
		}

		var academicYear academicYearModel.AcademicYear
		if err := config.DB.First(&academicYear, "number = ?", payload.AcademicYearNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Academic year not found"})
			return
		}

		generateADN, err := studentControllers.CreateADN(student.APNStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ADN"})
			return
		}

		updateStudent := studentModels.Student{
			APNStatus: "Archive",
			ADN:       generateADN,
		}

		if err := config.DB.Model(&student).Where("number= ?", studentPayload).Updates(updateStudent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update student"})
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
			Status: "ADMITTED",
		}

		//// Checks if the campus class applied for is the same as the campus class being admitted for.
		//if student.CampusClassNumber != payload.CampusClassNumber {
		//	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Campus class applied for and selected are different for student: %s %s %s", student.Surname, student.OtherName, student.SecondName)})
		//	return
		//}

		if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplicationStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
			return
		}

		saveAdmission := models.Admission{
			Number:             uuid.NewV4().String(),
			StudyPeriodNumber:  payload.StudyPeriodNumber,
			InterviewNumber:    payload.InterviewNumber,
			StudentNumber:      studentPayload,
			CampusClassNumber:  payload.CampusClassNumber,
			CampusNumber:       payload.CampusNumber,
			StreamNumber:       payload.StreamNumber,
			UserNumber:         user.Number,
			SchoolNumber:       payload.SchoolNumber,
			AcademicYearNumber: payload.AcademicYearNumber,
		}

		if err := config.DB.Create(&saveAdmission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission created successfully"})

}

func FetchAdmissions(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	campusNumber := c.Query("campus_number")
	studyPeriodNumber := c.Query("study_period_number")
	streamNumber := c.Query("stream_number")
	interviewNumber := c.Query("interview_number")

	var admissions []models.Admission

	if err := config.DB.Scopes(
		filters.ByStudyPeriodNumber(studyPeriodNumber),
		filters.ByCampusNumber(campusNumber),
		filters.ByStreamNumber(streamNumber),
		filters.ByInterviewNumber(interviewNumber),
	).Preload("Student").Preload("Campus").Preload("CampusClass").Preload("Stream").Preload("StudyPeriod").Preload("AcademicYear").Find(&admissions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch admissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admissions})
}

func FetchAdmissionById(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	admissionNumber := c.Param("admission_number")
	var admission models.Admission

	if err := config.DB.Preload("Student").Preload("Campus").Preload("CampusClass").Preload("Stream").Preload("StudyPeriod").Preload("AcademicYear").First(&admission, "number = ?", admissionNumber).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch admission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admission})
}
