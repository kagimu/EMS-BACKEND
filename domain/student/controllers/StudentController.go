package controllers

import (
	"errors"
	"fmt"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/services"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"net/http"

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

	campusClass, err := services.ResolveCampusClass(payload.ClassNumber, payload.CampusNumber, payload.CurriculumNumber, payload.CampusCurriculumNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error checking class and campus: " + err.Error()})
		return
	}

	if campusClass.Number == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Selected class is not active on the selected campus"})
		return
	}

	/*var campusClass classModels.CampusClass
	if err := config.DB.Where("number = ?", campusClass.Number).First(&campusClass).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Campus class not found"})
		return
	}*/

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
			PreferredCampusNumber:     payload.PreferredCampusNumber,
			SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
			ApplicationNumber:         applicationNumber,
			CampusClassNumber:         campusClass.Number,
			SchoolNumber:              payload.SchoolNumber,
			Lin:                       payload.Lin,
			ModeOfStudyNumber:         payload.ModeOfStudyNumber,
			AddedBy:                   userNumber,
		}

		if payload.ADN != "" {
			updateStudent.ADN = payload.ADN
		}

		if err := tx.Where("application_number = ?", applicationNumber).Updates(updateStudent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "failed to update student",
			})
			return
		}

		// update the user details
		updateUser := authModel.User{
			FirstName:  payload.SecondName,
			LastName:   payload.Surname,
			SecondName: payload.OtherName,
		}

		if err := tx.Where("number = ?", student.UserNumber).Updates(updateUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "failed to update user",
			})
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "failed to commit transaction",
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
		Gender:     payload.Gender,
		Type:       constants.TypeStudent,
		Status:     constants.StatusInactive,
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
		PreferredCampusNumber:     payload.PreferredCampusNumber,
		UserNumber:                createdUser.Number,
		SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
		ApplicationNumber:         applicationNumber,
		CampusClassNumber:         campusClass.Number,
		SchoolNumber:              payload.SchoolNumber,
		Lin:                       payload.Lin,
		ModeOfStudyNumber:         payload.ModeOfStudyNumber,
	}

	if payload.ADN != "" {
		savStudent.ADN = payload.ADN
	}

	if err := tx.Create(&savStudent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Saving student error: " + err.Error()})
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
		SchoolNumber:              payload.SchoolNumber,
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
	detailedStudent.UserNumber = student.UserNumber
	detailedStudent.SecondName = student.SecondName
	detailedStudent.ADNStatus = student.ADNStatus
	detailedStudent.DateOfBirth = student.DateOfBirth
	detailedStudent.ReasonForChoosingSchool = student.ReasonForChoosingSchool
	detailedStudent.Photo = commons.GenerateDownloadPath(student.Photo, true)

	// find the latest enrollment for the student
	var enrollment enrollmentModel.Enrollment
	if err := config.DB.Where("student_number = ?", studentNumber).
		Preload("CampusClass").Preload("Stream").Preload("StudyPeriod").
		Preload("StudyPeriod.AcademicYear").
		Preload("Campus").Preload("School").
		Order("created_at desc").First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})

		}
		return
	}

	studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(enrollment.StudyPeriod.StudyPeriodNameNumber)

	campusLevel, _ := repositories4.GetCampusLevelByNumber(enrollment.CampusClass.CampusLevelNumber)

	curriculumLevel, _ := repositories4.GetCurriculumLevelByNumber(campusLevel.CurriculumLevelNumber)

	var response = enrollmentResponse.Enrollment{}
	response.LevelNumber = curriculumLevel.LevelNumber
	response.CampusCurriculumNumber = campusLevel.CampusCurriculumNumber
	response.CurriculumNumber = curriculumLevel.CurriculumNumber
	response.StudentNumber = enrollment.StudentNumber
	response.EnrollmentDate = enrollment.EnrollmentDate
	response.Class = enrollment.CampusClass.ClassName
	response.ClassNumber = enrollment.CampusClass.ClassNumber
	response.Status = enrollment.Status
	response.Stream = enrollment.Stream.Name
	response.StreamNumber = enrollment.Stream.Number
	response.School = enrollment.School.Name
	response.SchoolNumber = enrollment.School.Number
	response.Campus = enrollment.Campus.Name
	response.CampusNumber = enrollment.Campus.Number
	response.StudyPeriod = studyPeriodName.Name
	response.StudyPeriodNumber = enrollment.StudyPeriod.Number
	response.CampusClassNumber = enrollment.CampusClass.Number
	response.ContinuingStudent = enrollment.ContinuingStudent
	response.AcademicYear = enrollment.StudyPeriod.AcademicYear.AcademicYearName
	detailedStudent.CurrentEnrollment = response
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

	if err := config.DB.Where("campus_class_number = ? AND status = ? AND study_period_number = ?", assessmentMetric.CampusClassNumber, constants.StatusActive, assessmentMetric.StudyPeriodNumber).Find(&enrollments).Error; err != nil {
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
			ActualScore:       assessment.ActualScore,
			Score:             assessment.Score,
			Class:             campusClass.ClassName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": students})
}

func CreateMedicalRecord(context *gin.Context) {
	var payload request.MedicalRecordRequest
	err := context.BindJSON(&payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}

	err = payload.Validate()

	err = services2.CreateMedicalRecord(payload)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Student created"})
}

func CreateFamilyStructure(context *gin.Context) {
	var payload request.FamilyStructureRequest
	err := context.BindJSON(&payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}

	err = payload.Validate()

	err = services2.CreateFamilyStructure(payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Student created"})
}

func CreateEmergencyInformation(context *gin.Context) {

	var payload request.EmergencyInformationRequest
	err := context.BindJSON(&payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}

	err = payload.Validate()

	err = services2.CreateEmergencyInformation(payload)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to process input"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Student created"})

}
