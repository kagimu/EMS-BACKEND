package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	"gorm.io/gorm"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateEnrollment(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.EnrollmentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// validate if student exists on the same campus and school and campus_class_number and study period and   is active
	var existingEnrollment models.Enrollment

	if err := config.DB.Where("student_number = ? AND study_period_number = ? AND status = ?", payload.StudentNumber, payload.StudyPeriodNumber, "active").First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student is already enrolled for this study period"})
		return
	}

	var admission admissionModel.Admission
	if err := config.DB.Where("student_number = ?", payload.StudentNumber).First(&admission).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch the student admission details"})
			return
		}
	}

	saveEnrollment := models.Enrollment{
		Number:            uuid.NewV4().String(),
		StudentNumber:     payload.StudentNumber,
		StreamNumber:      payload.StreamNumber,
		CampusNumber:      admission.CampusNumber,
		SchoolNumber:      admission.SchoolNumber,
		EnrollmentDate:    payload.EnrollmentDate,
		StudyPeriodNumber: admission.StudyPeriodNumber,
		CampusClassNumber: admission.CampusClassNumber,
		Status:            constants.StatusActive,
		UserNumber:        user.Number,
		ContinuingStudent: payload.ContinuingStudent,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveEnrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save enrollment"})
		return
	}

	var items []admissionModel.AdmissionCondition

	if err := config.DB.Preload("Item").Joins("JOIN items on items.number = admission_conditions.item_number AND items.mandatory=?", true).Where("admission_number = ?", admission.Number).Find(&items).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

	}

	if len(items) > 0 {
		//validate if all mandatory items have been checked
	}

	for _, item := range payload.Items {
		var enrolmentItem models.EnrollmentItem
		enrolmentItem.ItemNumber = item.ItemNumber
		enrolmentItem.EnrollmentNumber = saveEnrollment.Number
		enrolmentItem.Quantity = item.Quantity
		enrolmentItem.Brought = item.Brought
		enrolmentItem.UserNumber = user.Number
		if err := tx.Create(&enrolmentItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save enrollment"})
			return
		}
	}

	// update admission status to enrolled where student number is equal to payload.StudentNumber

	if err := tx.Where("student_number = ?", payload.StudentNumber).First(&admission).Error; err != nil {
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
		return
	}

	if admission.Number != "" {
		admission.Status = constants.StatusEnrolled
		if err := tx.Save(&admission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update admission status"})
			return
		}
	}

	// find all other enrollments for the student and set their status to inactive
	var enrollments []models.Enrollment
	if err := tx.Where("student_number = ? AND number != ? AND status = ? ", payload.StudentNumber, saveEnrollment.Number, constants.StatusActive).Find(&enrollments).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch other enrollments"})
			return
		}

	}

	for _, enrollment := range enrollments {
		enrollment.Status = constants.StatusInactive
		if err := tx.Save(&enrollment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update other enrollments status"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
}

func FetchEnrolledStudents(c *gin.Context) {

	paginatedResponse := pagination.PaginatedResponse{}
	students := make([]responses.Enrollment, 0)

	schoolNumber := c.Query("school_number")
	streamNumber := c.Query("stream_number")
	campusNumber := c.Query("campus_number")
	campusClassNumber := c.Query("campus_class_number")
	studyPeriodNumber := c.Query("study_period_number")
	status := c.Query("status")

	var enrolledStudents []models.Enrollment
	if err := config.DB.
		Scopes(
			filters.BySchoolNumber(schoolNumber),
			filters.ByStreamNumber(streamNumber),
			filters.ByCampusNumber(campusNumber),
			filters.ByCampusClassNumber(campusClassNumber),
			filters.ByStudyPeriodNumber(studyPeriodNumber),
			filters.ByStatus(status),
			pagination.Paginate(c, &paginatedResponse),
		).
		Preload("School").Preload("CampusClass").Preload("Campus").Preload("Student").Preload("StudyPeriod").Find(&enrolledStudents).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch enrolled students"})
		return
	}

	for _, enrollment := range enrolledStudents {
		students = append(students, responses.Enrollment{
			StudentSurname:    enrollment.Student.Surname,
			StudentSecondName: enrollment.Student.SecondName,
			StudentOtherName:  enrollment.Student.OtherName,
			StudentNumber:     enrollment.Student.Number,
			Number:            enrollment.Number,
			School:            enrollment.School.Name,
			Campus:            enrollment.Campus.Name,
			StudyPeriod:       enrollment.StudyPeriod.StudyPeriodName,
			Status:            enrollment.Status,
			EnrollmentDate:    enrollment.EnrollmentDate,
			ContinuingStudent: enrollment.ContinuingStudent,
			Class:             enrollment.CampusClass.ClassName + "(" + enrollment.CampusClass.ClassAbbreviation + ")",
			Adn:               enrollment.Student.ADN,
		})
	}

	paginatedResponse.Content = students

	c.JSON(http.StatusOK, paginatedResponse)
}

func TerminateStudentEnrollment(context *gin.Context) {

	studentNumber := context.Param("student_number")

	// get the enrollment number
	var enrollment models.Enrollment
	if err := config.DB.Where("student_number = ? AND status = ?", studentNumber, constants.StatusActive).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Student is not enrolled"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch student enrollment"})
		return
	}

	// update the enrollment status to inactive
	enrollment.Status = constants.StatusTerminated
	if err := config.DB.Save(&enrollment).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to terminate student enrollment"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Student enrollment terminated successfully"})
}
