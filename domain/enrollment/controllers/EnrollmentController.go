package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/services"
	"gorm.io/gorm"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/gin-gonic/gin"
)

var enrolmentService = services2.NewEnrolmentService()

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

	err = enrolmentService.NewEnrolment(payload, user.Number)

	// validate if student exists on the same campus and school and campus_class_number and study period and   is active
	//var existingEnrollment models.Enrollment

	/*if err := config.DB.Where("student_number = ? AND study_period_number = ? AND status = ?", payload.UserNumber, payload.StudyPeriodNumber, "active").First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student is already enrolled for this study period"})
		return
	}*/

	/*var admission admissionModel.Admission
	if err := config.DB.Where("student_number = ?", payload.UserNumber).First(&admission).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch the student admission details"})
			return
		}
	}

	student, err := repositories2.FindStudentByNumber(payload.UserNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch student details"})
		return
	}*/

	/*saveEnrollment := models.Enrollment{
		Number:            uuid.NewV4().String(),
		UserNumber:     payload.UserNumber,
		StreamNumber:      payload.StreamNumber,
		CampusNumber:      admission.CampusNumber,
		SchoolNumber:      admission.SchoolNumber,
		EnrollmentDate:    payload.EnrollmentDate,
		StudyPeriodNumber: admission.StudyPeriodNumber,
		CampusClassNumber: admission.CampusClassNumber,
		Status:            constants.StatusActive,
		UserNumber:        user.Number,
		ContinuingStudent: payload.ContinuingStudent,
	}*/

	//tx := config.DB.Begin()

	/*if err := tx.Create(&saveEnrollment).Error; err != nil {
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

	}*/

	//if len(items) > 0 {
	//validate if all mandatory items have been checked
	//}

	/*for _, item := range payload.Items {
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
	}*/

	// update admission status to enrolled where student number is equal to payload.UserNumber

	/*if err := tx.Where("student_number = ?", payload.UserNumber).First(&admission).Error; err != nil {
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
		return
	}*/

	/*if admission.Number != "" {
		admission.Status = constants.StatusEnrolled
		if err := tx.Save(&admission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update admission status"})
			return
		}
	}*/

	// find all other enrollments for the student and set their status to inactive
	/*var enrollments []models.Enrollment
	if err := tx.Where("student_number = ? AND number != ? AND status = ? ", payload.UserNumber, saveEnrollment.Number, constants.StatusActive).Find(&enrollments).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch other enrollments"})
			return
		}

	}*/

	/*for _, enrollment := range enrollments {
		enrollment.Status = constants.StatusInactive
		if err := tx.Save(&enrollment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update other enrollments status"})
			return
		}
	}*/

	/*if payload.FeedingGroupNumber != "" {
		fg, err := repositories.GetFeedingGroupByNumber(payload.FeedingGroupNumber)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding group"})
			return
		}

		if fg.Number != "" {
			err := services.CreateFeedingGroupUser(student.UserNumber, payload.CampusNumber, payload.FeedingGroupNumber, payload.Status, tx)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create feeding group user " + err.Error()})
				return
			}
		}

	}*/

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
}

func FetchEnrolledStudents(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	paginatedResponse := pagination.PaginatedResponse{}

	schoolNumber := c.Query("school_number")
	streamNumber := c.Query("stream_number")
	campusNumber := c.Query("campus_number")
	campusClassNumber := c.Query("campus_class_number")
	studyPeriodNumber := c.Query("study_period_number")
	status := c.Query("status")

	if user.Type != constants.TypeClient {
		if schoolNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view enrollments for"})
			return
		}
	}

	page := pagination.Paginate(c, &paginatedResponse)

	students_, err := enrolmentService.Enrolments(schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, status, user, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	paginatedResponse.Content = students_

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

func ReEnrollStudent(context *gin.Context) {
	studentNumber := context.Param("student_number")

	// get the enrollment number
	var enrollment models.Enrollment
	if err := config.DB.Where("student_number = ? AND status = ?", studentNumber, constants.StatusTerminated).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Student doesnt have a terminated enrollment"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch student enrollment"})
		return
	}

	// update the enrollment status to active
	enrollment.Status = constants.StatusActive
	if err := config.DB.Save(&enrollment).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to re-enroll student"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Student re-enrolled successfully"})

}
