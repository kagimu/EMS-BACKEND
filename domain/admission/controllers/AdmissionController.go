package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	var student studentModels.Student
	if err := config.DB.First(&student, "number = ?", payload.StudentNumber).Error; err != nil {
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

	generateADN, err := studentControllers.CreateADN(student.APNStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to generate ADN"})
		return
	}

	updateStudent := studentModels.Student{
		APNStatus: "Archive",
		ADN:       generateADN,
	}

	if err := config.DB.Model(&student).Where("number= ?", payload.StudentNumber).Updates(updateStudent).Error; err != nil {
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

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplicationStatus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
		return
	}

	for _, student := range payload.StudentNumber {
		saveAdmission := models.Admission{
			Number:            uuid.NewV4().String(),
			StudyPeriodNumber: payload.StudyPeriodNumber,
			InterviewNumber:   payload.InterviewNumber,
			StudentNumber:     student,
			CampusClassNumber: payload.CampusClassNumber,
			CampusNumber:      payload.CampusNumber,
			StreamNumber:      payload.StreamNumber,
			UserNumber:        user.Number,
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

	paginationParams := helpers.ParsePaginationParams(c)

	var admissions []models.Admission
	query := config.DB.Model(&models.Admission{}).Preload("Student").Preload("CampusClass").Preload("Campus").Preload("Stream")

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
		})
		return
	}

	if err := query.Find(&admissions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Admission fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    admissions,
		"message": "Admissions fetched successfully",
	})
}

func FetchAdmittedStudentsByClass(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	classNumber := c.Query("class_number")
	var class campusClassModels.CampusClass

	if classNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Class number is missing"})
		return
	}

	if err := config.DB.First(&class, "number = ?", classNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Class not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	var admissions []models.Admission

	if err := config.DB.Where("class_number = ?", classNumber).Preload("Campus").Preload("Stream").Preload("Student").Find(&admissions).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to fetch admission list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admissions})
}
