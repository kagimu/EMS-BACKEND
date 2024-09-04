package controllers

import (
	"fmt"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	campusClassModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"net/http"

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
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update student"})
			return
		}

		if err := tx.Save(&application).Error; err != nil {
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

		if err := tx.Create(&saveAdmission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission"})
			tx.Rollback()
			return
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
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	paginatedResponse := pagination.PaginatedResponse{}
	paginate := pagination.Paginate(c, &paginatedResponse)
	//adms := make([]models.Admission, 0)

	campusNumber := c.Query("campus_number")
	studyPeriodNumber := c.Query("study_period_number")
	streamNumber := c.Query("stream_number")
	interviewNumber := c.Query("interview_number")
	schoolNumber := c.Query("school_number")
	campusClassNumber := c.Query("campus_class_number")

	/*var admissions []models.Admission

	if err := config.DB.Scopes(
		filters.ByStudyPeriodNumber(studyPeriodNumber),
		filters.ByCampusNumber(campusNumber),
		filters.ByStreamNumber(streamNumber),
		filters.ByInterviewNumber(interviewNumber),
		filters.BySchoolNumber(schoolNumber),
		filters.ByCampusClassNumber(campusClassNumber),
		pagination.Paginate(c, &paginatedResponse),
	).Preload("Student").Preload("Campus").Preload("CampusClass").Preload("Stream").Preload("StudyPeriod").Preload("AcademicYear").Find(&admissions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch admissions"})
		return
	}

	if admissions != nil {
		adms = admissions
	}

	paginatedResponse.Content = adms*/

	admissions, err := admissionService.AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber, paginate)
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
