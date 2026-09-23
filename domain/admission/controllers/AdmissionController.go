package controllers

import (
	"fmt"
	"net/http"

	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	applicationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/services"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	studentControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var admissionService = services.NewAdmissionService()

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

	//var interview interviewModel.Interview
	//if err := config.DB.First(&interview, "number = ?", payload.InterviewNumber).Error; err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Interview not found"})
	//	return
	//}

	var academicYear academicYearModel.AcademicYear
	if err := config.DB.First(&academicYear, "number = ?", payload.AcademicYearNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Academic year not found"})
		return
	}

	for _, studentPayload := range payload.StudentNumber {

		var student studentModels.Student
		if err := config.DB.First(&student, "number = ?", studentPayload).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
			return
		}

		//validate if student has been admitted before on the same academic year and study period
		var admission models.Admission
		if err := config.DB.Where("student_number = ? AND academic_year_number = ? AND study_period_number = ?", studentPayload, payload.AcademicYearNumber, payload.StudyPeriodNumber).First(&admission).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Student has already been admitted for the selected academic year and study period"})
			return
		}

		// validate if student has been admitted before on the same academic year and study period and campus class

		if err := config.DB.Where("student_number = ? AND academic_year_number = ? AND study_period_number = ? AND campus_class_number = ?", studentPayload, payload.AcademicYearNumber, payload.StudyPeriodNumber, payload.CampusClassNumber).First(&admission).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Student has already been admitted for the selected academic year, study period and campus class"})
			return
		}

		// Update the application status
		applicationNumber := student.ApplicationNumber
		application, err := applicationController.FetchApplicationByNumber(c, applicationNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred while fetching application"})
			return
		}
		application.Status = constants.StatusAdmitted

		tx := config.DB.Begin()
		generateADN, err := studentControllers.CreateADN(student.APNStatus, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ADN"})
			return
		}

		student.APNStatus = "Archive"
		student.ADN = generateADN
		student.Lin = uuid.NewV4().String()
		// Update the student status
		student.Status = constants.StatusAdmitted
		if err := tx.Save(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update student details"})
			return
		}

		if err := tx.Save(&application).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application"})
			return
		}

		saveAdmission := models.Admission{
			Number:             uuid.NewV4().String(),
			StudyPeriodNumber:  payload.StudyPeriodNumber,
			StudentNumber:      studentPayload,
			CampusClassNumber:  payload.CampusClassNumber,
			CampusNumber:       payload.CampusNumber,
			StreamNumber:       payload.StreamNumber,
			UserNumber:         user.Number,
			SchoolNumber:       payload.SchoolNumber,
			AcademicYearNumber: payload.AcademicYearNumber,
		}

		if err := tx.Create(&saveAdmission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission"})
			tx.Rollback()
			return
		}

		// Create AdmissionInterview record if interview number is provided
		if payload.InterviewNumber != "" {
			admissionInterview := models.AdmissionInterview{
				Number:          uuid.NewV4().String(),
				InterviewNumber: payload.InterviewNumber,
				AdmissionNumber: saveAdmission.Number,
			}

			if err := tx.Create(&admissionInterview).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission interview record"})
				return
			}
		}

		for _, admissionOfferCondition := range payload.Items {
			saveAdmissionOfferConditions := models.AdmissionCondition{
				AdmissionNumber: saveAdmission.Number,
				ItemNumber:      admissionOfferCondition.ItemNumber,
				Quantity:        admissionOfferCondition.Quantity,
				UserNumber:      user.Number,
			}

			if err := tx.Create(&saveAdmissionOfferConditions).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission offer conditions"})
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction"})
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission created successfully"})

}

func FetchAdmissions(c *gin.Context) {

	// Heal any admission records previously saved as ADMITTED back to PENDING so they display and can be enrolled
	config.DB.Model(&models.Admission{}).Where("status = ?", constants.StatusAdmitted).Update("status", constants.StatusPending)

	user, err := commons.GetUserFromHeader(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	paginatedResponse := pagination.PaginatedResponse{}
	paginate := pagination.Paginate(c, &paginatedResponse)

	campusNumber := c.Query("campus_number")
	studyPeriodNumber := c.Query("study_period_number")
	streamNumber := c.Query("stream_number")
	schoolNumber := c.Query("school_number")
	campusClassNumber := c.Query("campus_class_number")

	if user.Type != constants.TypeClient {
		if schoolNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view admissions for"})
			return
		}
	}

	admissions, err := admissionService.AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, schoolNumber, campusClassNumber, user, paginate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	paginatedResponse.Content = admissions

	c.JSON(http.StatusOK, paginatedResponse)
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

func FetchItemsAddedOnAdmission(c *gin.Context) {
	admissionNumber := c.Param("admission_number")
	var items []models.AdmissionCondition

	if err := config.DB.Preload("Item").Where("admission_number = ?", admissionNumber).Find(&items).Error; err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	type Item struct {
		ItemNumber string  `json:"item_number"`
		Cost       float64 `json:"cost"`
		ItemName   string  `json:"item_name"`
		Quantity   float64 `json:"quantity"`
	}
	filteredItems := make([]Item, 0)

	for _, item := range items {
		filteredItems = append(filteredItems, Item{
			ItemNumber: item.ItemNumber,
			ItemName:   item.Item.Name,
			Cost:       item.Item.Cost,
			Quantity:   item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": filteredItems})
}

func CreateExistingAdmission(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	var payload requests.ExistingAdmissionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
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

	//var interview interviewModel.Interview
	//if err := config.DB.First(&interview, "number = ?", payload.InterviewNumber).Error; err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Interview not found"})
	//	return
	//}

	var academicYear academicYearModel.AcademicYear
	if err := config.DB.First(&academicYear, "number = ?", payload.AcademicYearNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Academic year not found"})
		return
	}

	var student studentModels.Student
	if err := config.DB.First(&student, "number = ?", payload.StudentNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		return
	}

	//validate if student has been admitted before on the same academic year and study period
	var admission models.Admission
	if err := config.DB.Where("student_number = ? AND academic_year_number = ? AND study_period_number = ?", payload.StudentNumber, payload.AcademicYearNumber, payload.StudyPeriodNumber).First(&admission).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student has already been admitted for the selected academic year and study period"})
		return
	}

	// validate if student has been admitted before on the same academic year and study period and campus class

	if err := config.DB.Where("student_number = ? AND academic_year_number = ? AND study_period_number = ? AND campus_class_number = ?", payload.StudentNumber, payload.AcademicYearNumber, payload.StudyPeriodNumber, payload.CampusClassNumber).First(&admission).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student has already been admitted for the selected academic year, study period and campus class"})
		return
	}

	// Update the application status
	applicationNumber := student.ApplicationNumber
	application, err := applicationController.FetchApplicationByNumber(c, applicationNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred while fetching application"})
		return
	}
	application.Status = constants.StatusAdmitted

	tx := config.DB.Begin()
	//generateADN, err := studentControllers.CreateADN(student.APNStatus, tx)
	//if err != nil {
	//	tx.Rollback()
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ADN"})
	//	return
	//}

	student.APNStatus = "Archive"
	//student.ADN = payload.ADN
	student.Lin = uuid.NewV4().String()
	// Update the student status
	student.Status = constants.StatusAdmitted
	if err := tx.Save(&student).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update student " + err.Error()})
		return
	}

	if err := tx.Save(&application).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update the application " + err.Error()})
		return
	}

	saveAdmission := models.Admission{
		Number:             uuid.NewV4().String(),
		StudyPeriodNumber:  payload.StudyPeriodNumber,
		StudentNumber:      payload.StudentNumber,
		CampusClassNumber:  payload.CampusClassNumber,
		CampusNumber:       payload.CampusNumber,
		StreamNumber:       payload.StreamNumber,
		UserNumber:         user.Number,
		SchoolNumber:       payload.SchoolNumber,
		AcademicYearNumber: payload.AcademicYearNumber,
	}

	if err := tx.Create(&saveAdmission).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission"})

		return
	}

	// Create AdmissionInterview record if interview number is provided
	//if payload.InterviewNumber != "" {
	//	admissionInterview := models.AdmissionInterview{
	//		Number:          uuid.NewV4().String(),
	//		InterviewNumber: payload.InterviewNumber,
	//		AdmissionNumber: saveAdmission.Number,
	//	}
	//
	//	if err := tx.Create(&admissionInterview).Error; err != nil {
	//		tx.Rollback()
	//		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission interview record"})
	//		return
	//	}
	//}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission created successfully"})

}

func BulkDirectAdmit(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	var payload requests.BulkDirectAdmitRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload: " + err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var school schoolModels.School
	if err := config.DB.First(&school, "number = ?", payload.SchoolNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "School not found"})
		return
	}

	var campus campusModels.Campus
	if err := config.DB.First(&campus, "number = ?", payload.CampusNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campus not found"})
		return
	}

	var stream campusClassModels.Stream
	if err := config.DB.First(&stream, "number = ?", payload.StreamNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Stream not found"})
		return
	}

	campusClassNumber := stream.CampusClassNumber
	if campusClassNumber == "" {
		campusClassNumber = payload.CampusClassNumber
	}

	var campusClass campusClassModels.CampusClass
	if err := config.DB.First(&campusClass, "number = ?", campusClassNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Campus class not found"})
		return
	}

	var academicYear academicYearModel.AcademicYear
	if err := config.DB.First(&academicYear, "number = ?", payload.AcademicYearNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Academic year not found"})
		return
	}

	var studyPeriod studyPeriodModel.StudyPeriod
	if err := config.DB.First(&studyPeriod, "number = ?", payload.StudyPeriodNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Study period not found"})
		return
	}

	var modeOfStudy schoolModels.ModeOfStudy
	if err := config.DB.First(&modeOfStudy, "school_number = ?", payload.SchoolNumber).Error; err != nil {
		config.DB.First(&modeOfStudy)
	}

	tx := config.DB.Begin()

	admittedCount := 0

	for i, l := range payload.Learners {
		rowNum := i + 1

		linValue := l.Lin
		if linValue != "" {
			var existingStudent studentModels.Student
			if err := tx.Where("lin = ?", linValue).First(&existingStudent).Error; err == nil {
				var existingAdm models.Admission
				if err := tx.Where("student_number = ? AND academic_year_number = ? AND study_period_number = ?", existingStudent.Number, payload.AcademicYearNumber, payload.StudyPeriodNumber).First(&existingAdm).Error; err == nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Row %d: Learner %s %s with LIN %s is already admitted for this study period", rowNum, l.Surname, l.SecondName, linValue)})
					return
				}
			}
		} else {
			linValue = uuid.NewV4().String()
		}

		studentUserNumber := uuid.NewV4().String()
		studentUser := authModels.User{
			Number:     studentUserNumber,
			FirstName:  l.SecondName,
			SecondName: l.OtherName,
			LastName:   l.Surname,
			Gender:     l.Gender,
			Type:       constants.TypeStudent,
			Status:     constants.StatusInactive,
		}
		if err := tx.Create(&studentUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to create student user account: %s", rowNum, err.Error())})
			return
		}

		appNumber := uuid.NewV4().String()
		application := applicationModels.Application{
			Number:                appNumber,
			UserNumber:            user.Number,
			Type:                  "direct-application",
			Status:                constants.StatusAdmitted,
			ApplicationFeePayment: true,
			InterviewMode:         "OFFLINE",
		}
		if err := tx.Create(&application).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to create application: %s", rowNum, err.Error())})
			return
		}

		apn, err := studentServices.CreateAPN(tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to generate APN: %s", rowNum, err.Error())})
			return
		}

		adnValue := l.Adn
		if adnValue == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Row %d: ADN (Admission Number) is required and cannot be empty", rowNum)})
			return
		}

		var existingStudentWithAdn studentModels.Student
		if err := tx.Where("adn = ?", adnValue).First(&existingStudentWithAdn).Error; err == nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Row %d: Learner with ADN '%s' already exists in the system", rowNum, adnValue)})
			return
		}
		adnStatus := "Active"

		studentNumber := uuid.NewV4().String()
		student := studentModels.Student{
			Number:                    studentNumber,
			APN:                       apn,
			ADN:                       adnValue,
			ADNStatus:                 adnStatus,
			Lin:                       linValue,
			APNStatus:                 "Archive",
			Surname:                   l.Surname,
			SecondName:                l.SecondName,
			OtherName:                 l.OtherName,
			Gender:                    l.Gender,
			DateOfBirth:               l.DateOfBirth,
			PreferredTypeOfEnrollment: "Day",
			PreferredCampusNumber:     payload.CampusNumber,
			CampusClassNumber:         campusClassNumber,
			SourceOfSchoolInfo:        "Direct Admission",
			UserNumber:                studentUserNumber,
			Status:                    constants.StatusAdmitted,
			SchoolNumber:              payload.SchoolNumber,
			ModeOfStudyNumber:         modeOfStudy.Number,
			ApplicationNumber:         appNumber,
			AddedBy:                   user.Number,
		}
		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to create student record: %s", rowNum, err.Error())})
			return
		}

		if l.ParentSurname != "" || l.ParentSecondName != "" || l.ParentPhone != "" {
			parentSurname := l.ParentSurname
			if parentSurname == "" {
				parentSurname = l.Surname
			}
			parentSecondName := l.ParentSecondName
			if parentSecondName == "" {
				parentSecondName = l.SecondName
			}
			parentEmail := l.ParentEmail
			if parentEmail == "" {
				parentEmail = fmt.Sprintf("parent_%s@school.local", studentNumber[:8])
			}

			parent := parentModels.Parent{
				Number:              uuid.NewV4().String(),
				Surname:             parentSurname,
				SecondName:          parentSecondName,
				OtherName:           "",
				ApplicationNumber:   appNumber,
				MainTelephoneNumber: l.ParentPhone,
				Email:               parentEmail,
				UserNumber:          user.Number,
			}
			if err := tx.Create(&parent).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to save parent details: %s", rowNum, err.Error())})
				return
			}
		}

		admission := models.Admission{
			Number:             uuid.NewV4().String(),
			StudentNumber:      studentNumber,
			CampusNumber:       payload.CampusNumber,
			CampusClassNumber:  campusClassNumber,
			StudyPeriodNumber:  payload.StudyPeriodNumber,
			AcademicYearNumber: payload.AcademicYearNumber,
			StreamNumber:       payload.StreamNumber,
			UserNumber:         user.Number,
			SchoolNumber:       payload.SchoolNumber,
			Status:             constants.StatusPending,
		}
		if err := tx.Create(&admission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Row %d: Failed to create admission record: %s", rowNum, err.Error())})
			return
		}

		admittedCount++
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully admitted %d learners", admittedCount),
		"count":   admittedCount,
	})
}

