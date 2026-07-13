package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateJobOffer(ctx *gin.Context) {
	var request requests.JobOfferRequest
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

	// 1. Find the actual JobApplication UUID using the ApplicationNumber (e.g. JOBAPLN2026001)
	var application models.JobApplication
	if err := config.DB.Where("application_number = ?", request.ApplicationNumber).First(&application).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job application not found"})
		return
	}

	// 2. Create JobOffer
	jobOffer := models.JobOffer{
		Number:                uuid.NewV4().String(),
		JobApplicationNumber:  application.Number, // use internal UUID
		JobInterviewNumber:    request.InterviewNumber,
		AppointmentLetterHtml: request.AppointmentLetterHtml,
		Status:                request.Status,
		UserNumber:            ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job offer not saved."})
		return
	}

	// 3. Update JobApplication status to OFFERED (or request.Status)
	application.Status = "OFFERED"
	config.DB.Save(&application)

	// 4. Update JobInterviewSchedule status to OFFERED
	// Looking up by application UUID and interview pool UUID
	var schedule models.JobInterviewSchedule
	if err := config.DB.Where("job_application_number = ? AND job_interview_number = ?", application.Number, request.InterviewNumber).First(&schedule).Error; err == nil {
		schedule.Status = "OFFERED"
		config.DB.Save(&schedule)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job offer created successfully!", "data": jobOffer})
}

func FetchJobOffers(ctx *gin.Context) {
	var jobOffers []models.JobOffer

	schoolNumber := ctx.Query("school_number")

	query := config.DB.Model(&models.JobOffer{}).
		Preload("JobApplication").
		Preload("JobInterview").
		Preload("JobInterview.JobVacancy").
		Preload("JobInterview.JobVacancy.JobPosition")

	if schoolNumber != "" {
		query = query.Joins("JOIN job_applications ON job_applications.number = job_offers.job_application_number").
			Joins("JOIN job_vacancies ON job_vacancies.number = job_applications.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&jobOffers).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch job offers"})
		return
	}

	var responseData []map[string]interface{}
	for _, jo := range jobOffers {
		responseData = append(responseData, map[string]interface{}{
			"number":                  jo.Number,
			"job_application_number":  jo.JobApplicationNumber,
			"job_interview_number":    jo.JobInterviewNumber,
			"appointment_letter_html": jo.AppointmentLetterHtml,
			"status":                  jo.Status,
			"application_no":          jo.JobApplication.ApplicationNumber,
			"applicant_name":          jo.JobApplication.FirstName + " " + jo.JobApplication.LastName,
			"applicant_email":         jo.JobApplication.Email,
			"applicant_phone":         jo.JobApplication.PhoneNumber,
			"interview_title":         jo.JobInterview.InterviewTitle,
			"vacancy_title":           jo.JobInterview.JobVacancy.VacancyTitle,
			"job_position_name":       jo.JobInterview.JobVacancy.JobPosition.PostName,
			"created_at":              jo.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func DeleteJobOffer(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var jobOffer models.JobOffer
	if err := config.DB.Unscoped().Where("number = ?", id).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	var relatedDeploymentsCount int64
	config.DB.Model(&models.StaffDeployment{}).Where("job_offer_number = ?", id).Count(&relatedDeploymentsCount)
	if relatedDeploymentsCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job offer. It has associated staff deployments."})
		return
	}

	// Force delete associated biodata records since staff is not deployed
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataPersonal{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataContact{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataFamily{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataEducation{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataMedical{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataLanguage{})
	config.DB.Unscoped().Where("job_offer_number = ?", id).Delete(&models.StaffBioDataReferences{})

	if err := config.DB.Unscoped().Delete(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Job offer not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job offer deleted successfully!"})
}
