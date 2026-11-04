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

	if err := config.DB.Omit("JobApplication", "JobInterview").Create(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job offer not saved: " + err.Error()})
		return
	}

	// 3. Update JobApplication status to OFFERED (or request.Status)
	config.DB.Model(&models.JobApplication{}).Where("number = ?", application.Number).Update("status", "OFFERED")

	// 4. Update JobInterviewSchedule status to OFFERED
	config.DB.Model(&models.JobInterviewSchedule{}).Where("job_application_number = ? AND job_interview_number = ?", application.Number, request.InterviewNumber).Update("status", "OFFERED")

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

	// ── PRE-FLIGHT CHECKS ──

	// Block if a deployment still exists for this offer
	var deploymentCount int64
	config.DB.Model(&models.StaffDeployment{}).Where("job_offer_number = ?", id).Count(&deploymentCount)
	if deploymentCount > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Cannot delete job offer. A staff deployment exists for this offer.",
			"hint":    "Delete the staff deployment first.",
		})
		return
	}


	// ── ALL CHECKS PASSED — delete in a transaction ──
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction: " + tx.Error.Error()})
		return
	}

	rollbackWithError := func(msg string, err error) {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg + ": " + err.Error()})
	}

	// Delete BioData leaf children first, then parents
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_family_children
		WHERE family_number IN (
			SELECT number FROM staff_bio_data_families WHERE job_offer_number = ?
		)`, id).Error; err != nil {
		rollbackWithError("Failed to delete bio data family children", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_families WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data family", err)
		return
	}
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_education_items
		WHERE education_number IN (
			SELECT number FROM staff_bio_data_educations WHERE job_offer_number = ?
		)`, id).Error; err != nil {
		rollbackWithError("Failed to delete bio data education items", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_educations WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data education", err)
		return
	}
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_language_items
		WHERE language_record_number IN (
			SELECT number FROM staff_bio_data_languages WHERE job_offer_number = ?
		)`, id).Error; err != nil {
		rollbackWithError("Failed to delete bio data language items", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_languages WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data language", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_contacts WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data contact", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_medicals WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data medical", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_references WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data references", err)
		return
	}
	if err := tx.Exec("DELETE FROM staff_bio_data_personals WHERE job_offer_number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete bio data personal", err)
		return
	}

	// Reset application status back to 'Interviewed' (undoing the OFFERED status)
	if jobOffer.JobApplicationNumber != "" {
		if err := tx.Model(&models.JobApplication{}).Where("number = ?", jobOffer.JobApplicationNumber).Update("status", "Interviewed").Error; err != nil {
			rollbackWithError("Failed to reset application status", err)
			return
		}
	}

	// Delete the job offer itself
	if err := tx.Unscoped().Where("number = ?", id).Delete(&models.JobOffer{}).Error; err != nil {
		rollbackWithError("Failed to delete job offer", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job offer and all staff bio-data deleted successfully."})
}
