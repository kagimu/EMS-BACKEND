package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	studentControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	studentNumber := c.Query("student_number")
	var student studentModels.Student
	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	}

	classNumber := c.Query("class_number")
	var class classModels.Class
	if err := config.DB.First(&class, "number = ?", classNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Class not found"})
		return
	}

	campusNumber := c.Query("campus_number")
	var campus campusModels.Campus
	if err := config.DB.First(&campus, "number = ?", campusNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Campus not found"})
		return
	}

	streamNumber := c.Query("stream_number")
	var stream classModels.Stream
	if err := config.DB.First(&stream, "number = ?", streamNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Stream not found"})
		return
	}

	saveAdmission := models.Admission{
		Number:           uuid.NewV1().String(),
		Status:           payload.Status,
		AdmissionFeePaid: payload.AdmissionFeePaid,
		StudentNumber:    payload.StudentNumber,
		ClassNumber:      payload.CampusNumber,
		CampusNumber:     payload.CampusNumber,
		StreamNumber:     payload.StreamNumber,
		UserNumber:       user.Number,
	}

	if err := config.DB.Create(&saveAdmission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission created successfully", "data": saveAdmission})

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
	query := config.DB.Model(&models.Admission{}).Preload("Student").Preload("Class").Preload("Campus").Preload("Stream")

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

func AdmitStudent(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Param("number")
	var student studentModels.Student
	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Student not found"})
		return
	}

	// Checks if the invoice is paid
	var invoice invoiceModel.Invoice
	if err := config.DB.First(&invoice, "student_number = ?", student.Number).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to fetch invoice for the student"})
		return
	}

	if invoice.Status != "PAID" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student invoice is not paid", "debug": invoice})
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

	if err := config.DB.Model(&student).Where("number= ?", studentNumber).Updates(updateStudent).Error; err != nil {
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
		Status: "INTERVIEW_SCHEDULED",
	}

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplicationStatus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func BulkAdmitStudents(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	type Input struct {
		StudentNumbers []string `json:"student_numbers"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	for _, studentNumber := range payload.StudentNumbers {
		var student studentModels.Student
		if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Student not found"})
			return
		}

		// Checks if the invoice is paid
		var invoice invoiceModel.Invoice
		if err := config.DB.First(&invoice, "student_number = ?", student.Number).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Failed to fetch invoice for the student"})
			return
		}

		if invoice.Status != "PAID" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Student invoice is not paid", "debug": invoice})
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

		if err := config.DB.Model(&student).Where("number= ?", studentNumber).Updates(updateStudent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update student", "err": err.Error()})
			return
		}

		// Update the application status for every student
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Students admitted successfully"})
}
