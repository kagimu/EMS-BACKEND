package controllers

import (
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateJobVacancy(ctx *gin.Context) {
	var request requests.JobVacancyRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	jobVacancy := models.JobVacancy{
		Number:            uuid.NewV4().String(),
		JobPositionNumber: request.JobPositionNumber,
		VacancyTitle:      request.VacancyTitle,
		ReferenceNumber:   request.ReferenceNumber,
		DepartmentNumber:  request.DepartmentNumber,
		SchoolNumber:      request.SchoolNumber,
		BranchLocation:    request.BranchLocation,
		NumberOfPositions: request.NumberOfPositions,
		AppointmentTerms:  request.AppointmentTerms,
		Status:            request.Status,
		UserNumber:        ctx.Request.Header.Get("User-Number"),
	}

	if request.OpeningDate != "" {
		parsedDate, _ := time.Parse("2006-01-02", request.OpeningDate)
		jobVacancy.OpeningDate = parsedDate
	}
	
	if request.ClosingDate != "" {
		parsedDate, _ := time.Parse("2006-01-02", request.ClosingDate)
		jobVacancy.ClosingDate = parsedDate
	}

	if err := config.DB.Create(&jobVacancy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job vacancy not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job vacancy created successfully!", "data": jobVacancy})
}

func FetchJobVacancies(ctx *gin.Context) {
	var jobVacancies []models.JobVacancy
	
	schoolNumber := ctx.Query("school_number")
	departmentNumber := ctx.Query("department_number")
	
	query := config.DB.Model(&models.JobVacancy{}).Preload("School").Preload("Department").Preload("JobPosition").Preload("Campus")
	
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	
	if departmentNumber != "" {
		query = query.Where("department_number = ?", departmentNumber)
	}

	if err := query.Find(&jobVacancies).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch job vacancies"})
		return
	}

	var responseData []map[string]interface{}
	for _, jv := range jobVacancies {
		openingDate := ""
		if !jv.OpeningDate.IsZero() {
			openingDate = jv.OpeningDate.Format("2006-01-02")
		}
		
		closingDate := ""
		if !jv.ClosingDate.IsZero() {
			closingDate = jv.ClosingDate.Format("2006-01-02")
		}

		responseData = append(responseData, map[string]interface{}{
			"number":              jv.Number,
			"job_position_number": jv.JobPositionNumber,
			"job_position_name":   jv.JobPosition.PostName,
			"job_code":            jv.JobPosition.JobCode,
			"vacancy_title":       jv.VacancyTitle,
			"reference_number":    jv.ReferenceNumber,
			"school_number":       jv.SchoolNumber,
			"school_name":         jv.School.Name,
			"department_number":   jv.DepartmentNumber,
			"department_name":     jv.Department.Name,
			"branch_location":     jv.BranchLocation,
			"campus_name":         jv.Campus.Name,
			"number_of_positions": jv.NumberOfPositions,
			"appointment_terms":   jv.AppointmentTerms,
			"opening_date":        openingDate,
			"closing_date":        closingDate,
			"status":              jv.Status,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func EditJobVacancy(ctx *gin.Context) {
	var request requests.EditJobVacancyRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var jobVacancy models.JobVacancy
	if err := config.DB.Where("number = ?", request.Number).First(&jobVacancy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job vacancy not found"})
		return
	}

	jobVacancy.JobPositionNumber = request.JobPositionNumber
	jobVacancy.VacancyTitle = request.VacancyTitle
	jobVacancy.ReferenceNumber = request.ReferenceNumber
	jobVacancy.DepartmentNumber = request.DepartmentNumber
	jobVacancy.SchoolNumber = request.SchoolNumber
	jobVacancy.BranchLocation = request.BranchLocation
	jobVacancy.NumberOfPositions = request.NumberOfPositions
	jobVacancy.AppointmentTerms = request.AppointmentTerms
	jobVacancy.Status = request.Status

	if request.OpeningDate != "" {
		parsedDate, _ := time.Parse("2006-01-02", request.OpeningDate)
		jobVacancy.OpeningDate = parsedDate
	}
	
	if request.ClosingDate != "" {
		parsedDate, _ := time.Parse("2006-01-02", request.ClosingDate)
		jobVacancy.ClosingDate = parsedDate
	}
	
	if err := config.DB.Save(&jobVacancy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job vacancy not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job vacancy updated successfully!"})
}

func DeleteJobVacancy(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job vacancy number is required"})
		return
	}

	var jobVacancy models.JobVacancy
	if err := config.DB.Unscoped().Where("number = ?", id).First(&jobVacancy).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job vacancy not found"})
		return
	}

	var relatedAppsCount int64
	config.DB.Model(&models.JobApplication{}).Where("job_vacancy_number = ?", id).Count(&relatedAppsCount)
	if relatedAppsCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job vacancy. It has associated applications."})
		return
	}

	var relatedInterviewsCount int64
	config.DB.Model(&models.JobInterview{}).Where("job_vacancy_number = ?", id).Count(&relatedInterviewsCount)
	if relatedInterviewsCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job vacancy. It has associated interviews."})
		return
	}

	if err := config.DB.Unscoped().Delete(&jobVacancy).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Job vacancy not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job vacancy deleted successfully!"})
}
