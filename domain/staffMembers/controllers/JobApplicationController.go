package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func generateApplicationNumber() string {
	currentYear := time.Now().Format("2006")
	prefix := "JOBAPLN" + currentYear

	var lastApplication models.JobApplication
	config.DB.Where("application_number LIKE ?", prefix+"%").Order("created_at desc").First(&lastApplication)

	nextNumber := 1
	if lastApplication.ApplicationNumber != "" && len(lastApplication.ApplicationNumber) > len(prefix) {
		lastSuffix := lastApplication.ApplicationNumber[len(prefix):]
		if num, err := strconv.Atoi(lastSuffix); err == nil {
			nextNumber = num + 1
		}
	}

	return fmt.Sprintf("%s%03d", prefix, nextNumber)
}

func CreateJobApplication(ctx *gin.Context) {
	var request requests.JobApplicationRequest
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

	jobApplication := models.JobApplication{
		Number:            uuid.NewV4().String(),
		ApplicationNumber: generateApplicationNumber(),
		JobVacancyNumber:  request.JobVacancyNumber,
		FirstName:         request.FirstName,
		LastName:          request.LastName,
		Email:             request.Email,
		PhoneNumber:       request.PhoneNumber,
		Status:            "Pending",
		UserNumber:        ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&jobApplication).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job application not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job application submitted successfully!", "data": jobApplication})
}

func FetchJobApplications(ctx *gin.Context) {
	var jobApplications []models.JobApplication

	jobVacancyNumber := ctx.Query("job_vacancy_number")
	schoolNumber := ctx.Query("school_number")

	query := config.DB.Model(&models.JobApplication{}).Preload("JobVacancy").Preload("JobVacancy.JobPosition")

	if jobVacancyNumber != "" {
		query = query.Where("job_vacancy_number = ?", jobVacancyNumber)
	}

	if schoolNumber != "" {
		query = query.Joins("JOIN job_vacancies ON job_vacancies.number = job_applications.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&jobApplications).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch job applications"})
		return
	}

	var responseData []map[string]interface{}
	for _, ja := range jobApplications {
		responseData = append(responseData, map[string]interface{}{
			"number":             ja.Number,
			"application_number": ja.ApplicationNumber,
			"job_vacancy_number": ja.JobVacancyNumber,
			"vacancy_title":      ja.JobVacancy.VacancyTitle,
			"job_position_name":  ja.JobVacancy.JobPosition.PostName,
			"first_name":         ja.FirstName,
			"last_name":          ja.LastName,
			"email":              ja.Email,
			"phone_number":       ja.PhoneNumber,
			"status":             ja.Status,
			"created_at":         ja.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func EditJobApplication(ctx *gin.Context) {
	var request requests.EditJobApplicationRequest
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

	var jobApplication models.JobApplication
	if err := config.DB.Where("number = ?", request.Number).First(&jobApplication).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job application not found"})
		return
	}

	jobApplication.JobVacancyNumber = request.JobVacancyNumber
	jobApplication.FirstName = request.FirstName
	jobApplication.LastName = request.LastName
	jobApplication.Email = request.Email
	jobApplication.PhoneNumber = request.PhoneNumber
	jobApplication.Status = request.Status

	if err := config.DB.Save(&jobApplication).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job application not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job application details updated successfully!"})
}

func DeleteJobApplication(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job application number is required"})
		return
	}

	var jobApplication models.JobApplication
	if err := config.DB.Unscoped().Where("number = ?", id).First(&jobApplication).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job application not found"})
		return
	}

	var relatedSchedulesCount int64
	config.DB.Model(&models.JobInterviewSchedule{}).Where("job_application_number = ?", id).Count(&relatedSchedulesCount)
	if relatedSchedulesCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job application. It is scheduled for an interview."})
		return
	}

	var relatedOffersCount int64
	config.DB.Model(&models.JobOffer{}).Where("job_application_number = ?", id).Count(&relatedOffersCount)
	if relatedOffersCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job application. It has been offered a job."})
		return
	}

	if err := config.DB.Unscoped().Delete(&jobApplication).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Job application not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job application deleted successfully!"})
}
