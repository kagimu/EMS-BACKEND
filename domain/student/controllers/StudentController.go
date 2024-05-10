package controllers

import (
	"errors"
	"fmt"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	assessmentMetricModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	enrollmentResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateStudent(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var payload request.StudentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("failed to validate input %s:", err.Error()),
		})
		return
	}

	applicationNumber := payload.ApplicationNumber
	var application applicationModel.Application
	if err := config.DB.Where("number = ?", applicationNumber).First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Application not found"})
		return
	}

	var campusClass classModels.CampusClass
	if err := config.DB.Where("number = ?", payload.CampusClassNumber).First(&campusClass).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Campus class not found"})
		return
	}

	tx := config.DB.Begin()
	studentAPN, err := CreateAPN(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create APN"})
		return
	}

	studentNumber := uuid.NewV4().String()
	userNumber := user.Number

	// find Student by application number

	var student models.Student
	tx.Where("application_number = ?", applicationNumber).First(&student)

	// NOTE: To avoid duplicate records, check if the student number is empty else, update the existing student
	if student.Number != "" {
		// update student details
		updateStudent := models.Student{
			Surname:                   payload.Surname,
			OtherName:                 payload.OtherName,
			Gender:                    payload.Gender,
			DateOfBirth:               payload.DateOfBirth,
			ADNStatus:                 payload.ADNStatus,
			ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
			PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
			SecondName:                payload.SecondName,
			PreferedCampusNumber:      payload.PreferredCampusNumber,
			UserNumber:                userNumber,
			SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
			ApplicationNumber:         applicationNumber,
			CampusClassNumber:         payload.CampusClassNumber,
		}

		if err := tx.Where("application_number = ?", applicationNumber).Updates(updateStudent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "failed to update student",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "updated student details successfully"})
		return
	}

	var createdUser = authModel.User{
		Number:     uuid.NewV4().String(),
		FirstName:  payload.SecondName,
		LastName:   payload.Surname,
		SecondName: payload.OtherName,
	}

	if err := tx.Create(&createdUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	savStudent := models.Student{
		Number:                    studentNumber,
		APN:                       studentAPN,
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		Gender:                    payload.Gender,
		DateOfBirth:               payload.DateOfBirth,
		ADNStatus:                 payload.ADNStatus,
		ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		SecondName:                payload.SecondName,
		PreferedCampusNumber:      payload.PreferredCampusNumber,
		UserNumber:                createdUser.Number,
		SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
		ApplicationNumber:         applicationNumber,
		CampusClassNumber:         payload.CampusClassNumber,
	}

	if err := tx.Create(&savStudent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": " Failed to save student"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "student created successfully",
		"data":    savStudent,
	})
}

func EditStudent(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var payload request.StudentUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to process input",
		})
		return
	}

	payload.Prepare()

	var student models.Student
	studentNumber := c.Param("student_number")
	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No student found with the specified number"})
		return
	}

	savedStudent := models.Student{
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		SecondName:                payload.SecondName,
		APNStatus:                 payload.APNStatus,
		DateOfBirth:               payload.DateOfBirth,
		ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		CampusClassNumber:         payload.CampusClassNumber,
	}

	if err := config.DB.Model(&student).Where("number = ?", studentNumber).Updates(savedStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})
}

func DeleteStudent(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Param("student_number")
	student := models.Student{}

	err := config.DB.First(&student, studentNumber).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no student found"})
		return
	}

	if err := config.DB.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student deleted successfully"})
}

func FetchStudentByNumber(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var student models.Student
	studentNumber := c.Param("student_number")

	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student number is missing"})
		return
	}

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "no student found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}
		return
	}

	var detailedStudent = responses.DetailedStudentProfile{}

	detailedStudent.Surname = student.Surname
	detailedStudent.OtherName = student.OtherName
	detailedStudent.SecondName = student.SecondName
	detailedStudent.ADNStatus = student.ADNStatus
	detailedStudent.DateOfBirth = student.DateOfBirth
	detailedStudent.ReasonForChoosingSchool = student.ReasonForChoosingSchool

	// find the latest enrollment for the student
	var enrollment enrollmentModel.Enrollment
	if err := config.DB.Where("student_number = ?", studentNumber).Order("created_at desc").First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})

		}
		return
	}

	var enrollmentResponse = enrollmentResponse.Enrollment{}
	enrollmentResponse.StudentNumber = enrollment.StudentNumber
	enrollmentResponse.EnrollmentDate = enrollment.EnrollmentDate
	enrollmentResponse.Class = enrollment.CampusClassNumber
	enrollmentResponse.Status = enrollment.Status

	detailedStudent.CurrentEnrollment = enrollmentResponse
	//enrollmentResponse.Stream =

	c.JSON(http.StatusOK, detailedStudent)
}

func FetchAllStudents(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var students []models.Student

	if err := config.DB.Preload("CampusClass").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": students})
}

func FetchStudentsByAssessmentMetricNumber(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	assessmentMetricNumber := c.Param("assessment_metric_number")

	if assessmentMetricNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric number is required"})
		return
	}

	var assessmentMetric assessmentMetricModel.AssessmentMetric

	if err := config.DB.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric not found"})
		return
	}

	//fetch enrollment students from assessment metric campus class number

	var enrollments []enrollmentModel.Enrollment

	if err := config.DB.Where("campus_class_number = ? AND status = ? ", assessmentMetric.CampusClassNumber, constants.STATUS_ACTIVE).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch enrollments"})
		return

	}

	var students []enrollmentResponse.EnrolledStudent

	for _, enrollment := range enrollments {
		var student models.Student
		if err := config.DB.Where("number = ?", enrollment.StudentNumber).First(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch student"})
			return
		}

		// find assesment for student and assessment metric

		var assessment assessmentMetricModel.Assessment

		if err := config.DB.Where("student_number = ? AND assessment_metric_number = ?", student.Number, assessmentMetric.Number).First(&assessment).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessments " + err.Error()})
				return
			}
		}

		var campusClass classModels.CampusClass
		if err := config.DB.Where("number = ?", enrollment.CampusClassNumber).First(&campusClass).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Campus class not found"})
			return
		}

		students = append(students, enrollmentResponse.EnrolledStudent{

			Number:            student.Number,
			StudentSurname:    student.Surname,
			StudentOtherName:  student.OtherName,
			StudentSecondName: student.SecondName,
			Comment:           assessment.Comment,
			Score:             assessment.ActualScore,
			Class:             campusClass.ClassName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": students})
}
