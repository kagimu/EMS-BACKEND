package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAdmissionOffer(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AdmissionOfferRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	studentNumber := payload.StudentNumber
	var student studentModel.Student
	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	}

	classNumber := payload.ClassNumber
	var class classModel.Class
	if err := config.DB.First(&class, "number = ?", classNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Class not found"})
		return
	}

	campusNumber := payload.CampusNumber
	var campus campusModel.Campus
	if err := config.DB.First(&campus, "number = ?", campusNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campus not found"})
		return
	}

	schoolNumber := payload.SchoolNumber
	var school schoolModel.School
	if err := config.DB.First(&school, "number = ?", schoolNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "School not found"})
		return
	}

	saveAdmissionOffer := models.AdmissionOffer{
		Number:        uuid.NewV1().String(),
		SchoolNumber:  payload.SchoolNumber,
		ClassNumber:   payload.ClassNumber,
		CampusNumber:  payload.CampusNumber,
		StudentNumber: payload.StudentNumber,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveAdmissionOffer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAdmissionOffer})

}

func FetchAdmissionOffers(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	paginationParams := helpers.ParsePaginationParams(c)

	var admissionOffers []models.AdmissionOffer
	query := config.DB.Model(&models.AdmissionOffer{}).Preload("Student").Preload("Class").Preload("Campus").Preload("School")

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
		})
		return
	}

	if err := query.Find(&admissionOffers).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Admission offer fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, admissionOffers)
}
