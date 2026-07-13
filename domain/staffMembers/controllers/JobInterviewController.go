package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateJobInterview(ctx *gin.Context) {
	var request requests.JobInterviewRequest
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

	var existingInterview models.JobInterview
	if err := config.DB.Where("job_vacancy_number = ? AND LOWER(interview_title) = LOWER(?)", request.JobVacancyNumber, request.InterviewTitle).First(&existingInterview).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "A job interview with this title already exists for this vacancy!"})
		return
	}

	status := request.Status
	if status == "" {
		status = "Scheduled"
	}

	jobInterview := models.JobInterview{
		Number:           uuid.NewV4().String(),
		JobVacancyNumber: request.JobVacancyNumber,
		InterviewTitle:   request.InterviewTitle,
		InterviewDate:    request.InterviewDate,
		InterviewTime:    request.InterviewTime,
		LocationOrLink:   request.LocationOrLink,
		Status:           status,
		UserNumber:       ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Omit("JobVacancy").Create(&jobInterview).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job interview not saved: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job interview created successfully!", "data": jobInterview})
}

func FetchJobInterviews(ctx *gin.Context) {
	var jobInterviews []models.JobInterview

	query := config.DB.Model(&models.JobInterview{}).Preload("JobVacancy").Preload("JobVacancy.JobPosition")

	jobVacancyNumber := ctx.Query("job_vacancy_number")
	if jobVacancyNumber != "" {
		query = query.Where("job_vacancy_number = ?", jobVacancyNumber)
	}

	schoolNumber := ctx.Query("school_number")
	if schoolNumber != "" {
		query = query.Joins("JOIN job_vacancies ON job_vacancies.number = job_interviews.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&jobInterviews).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch job interviews"})
		return
	}

	var responseData []map[string]interface{}
	for _, ji := range jobInterviews {
		responseData = append(responseData, map[string]interface{}{
			"number":             ji.Number,
			"job_vacancy_number": ji.JobVacancyNumber,
			"vacancy_title":      ji.JobVacancy.VacancyTitle,
			"job_position_name":  ji.JobVacancy.JobPosition.PostName,
			"interview_title":    ji.InterviewTitle,
			"interview_date":     ji.InterviewDate,
			"interview_time":     ji.InterviewTime,
			"location_or_link":   ji.LocationOrLink,
			"status":             ji.Status,
			"created_at":         ji.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func EditJobInterview(ctx *gin.Context) {
	var request requests.EditJobInterviewRequest
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

	var jobInterview models.JobInterview
	if err := config.DB.Where("number = ?", request.Number).First(&jobInterview).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job interview not found"})
		return
	}

	var existingInterview models.JobInterview
	if err := config.DB.Where("job_vacancy_number = ? AND LOWER(interview_title) = LOWER(?) AND number != ?", request.JobVacancyNumber, request.InterviewTitle, request.Number).First(&existingInterview).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "A job interview with this title already exists for this vacancy!"})
		return
	}

	jobInterview.JobVacancyNumber = request.JobVacancyNumber
	jobInterview.InterviewTitle = request.InterviewTitle
	jobInterview.InterviewDate = request.InterviewDate
	jobInterview.InterviewTime = request.InterviewTime
	jobInterview.LocationOrLink = request.LocationOrLink
	jobInterview.Status = request.Status

	if err := config.DB.Save(&jobInterview).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job interview not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job interview updated successfully!"})
}

func DeleteJobInterview(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job interview number is required"})
		return
	}

	var jobInterview models.JobInterview
	if err := config.DB.Unscoped().Where("number = ?", id).First(&jobInterview).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job interview not found"})
		return
	}

	var relatedSchedulesCount int64
	config.DB.Model(&models.JobInterviewSchedule{}).Where("job_interview_number = ?", id).Count(&relatedSchedulesCount)
	if relatedSchedulesCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job interview. It has scheduled candidates."})
		return
	}

	var relatedOffersCount int64
	config.DB.Model(&models.JobOffer{}).Where("job_interview_number = ?", id).Count(&relatedOffersCount)
	if relatedOffersCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job interview. It has associated job offers."})
		return
	}

	if err := config.DB.Unscoped().Delete(&jobInterview).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Job interview not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job interview deleted successfully!"})
}

func ScheduleInterview(ctx *gin.Context) {
	var request requests.ScheduleInterviewRequest
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

	// 1. Create JobInterviewSchedule
	schedule := models.JobInterviewSchedule{
		Number:               uuid.NewV4().String(),
		JobApplicationNumber: request.ApplicationNumber,
		JobInterviewNumber:   request.InterviewNumber,
		Remarks:              request.Remarks,
		UserNumber:           ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Omit("JobApplication", "JobInterview").Create(&schedule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Could not schedule the interview: " + err.Error()})
		return
	}

	// 2. Update JobApplication Status to "INTERVIEW_SCHEDULED"
	config.DB.Model(&models.JobApplication{}).Where("number = ?", request.ApplicationNumber).Update("status", "INTERVIEW_SCHEDULED")

	ctx.JSON(http.StatusOK, gin.H{"message": "Candidate successfully scheduled for interview!"})
}

func FetchScheduledInterviews(ctx *gin.Context) {
	var schedules []models.JobInterviewSchedule

	query := config.DB.Model(&models.JobInterviewSchedule{}).
		Preload("JobApplication").
		Preload("JobInterview")

	interviewNumber := ctx.Query("interview_number")
	if interviewNumber != "" {
		query = query.Where("job_interview_number = ?", interviewNumber)
	}
	
	applicationNumber := ctx.Query("application_number")
	if applicationNumber != "" {
		query = query.Where("job_application_number = ?", applicationNumber)
	}

	schoolNumber := ctx.Query("school_number")
	if schoolNumber != "" {
		query = query.Joins("JOIN job_interviews ON job_interviews.number = job_interview_schedules.job_interview_number").
			Joins("JOIN job_vacancies ON job_vacancies.number = job_interviews.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&schedules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch scheduled interviews"})
		return
	}

	var responseData []map[string]interface{}
	for _, sch := range schedules {
		responseData = append(responseData, map[string]interface{}{
			"number":                 sch.Number,
			"job_application_number": sch.JobApplicationNumber,
			"job_interview_number":   sch.JobInterviewNumber,
			"remarks":                sch.Remarks,
			"application_no":         sch.JobApplication.ApplicationNumber,
			"applicant_name":         sch.JobApplication.FirstName + " " + sch.JobApplication.LastName,
			"applicant_email":        sch.JobApplication.Email,
			"applicant_phone":        sch.JobApplication.PhoneNumber,
			"interview_title":        sch.JobInterview.InterviewTitle,
			"interview_date":         sch.JobInterview.InterviewDate,
			"interview_time":         sch.JobInterview.InterviewTime,
			"location_or_link":       sch.JobInterview.LocationOrLink,
			"comment":                sch.Comment,
			"status":                 sch.Status,
			"created_at":             sch.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func ScoreInterview(ctx *gin.Context) {
	var request requests.ScoreInterviewRequest
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

	var schedule models.JobInterviewSchedule
	if err := config.DB.Where("number = ?", request.ScheduledInterviewNumber).First(&schedule).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Scheduled interview not found"})
		return
	}

	schedule.Comment = request.Comment
	schedule.Status = request.Status

	if err := config.DB.Save(&schedule).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Score not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interview scored successfully!"})
}

func FetchCompletedInterviews(ctx *gin.Context) {
	var schedules []models.JobInterviewSchedule

	query := config.DB.Model(&models.JobInterviewSchedule{}).
		Preload("JobApplication").
		Preload("JobInterview").
		Where("job_interview_schedules.status = ?", "completed")

	schoolNumber := ctx.Query("school_number")
	if schoolNumber != "" {
		query = query.Joins("JOIN job_interviews ON job_interviews.number = job_interview_schedules.job_interview_number").
			Joins("JOIN job_vacancies ON job_vacancies.number = job_interviews.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&schedules).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch completed interviews"})
		return
	}

	var responseData []map[string]interface{}
	for _, sch := range schedules {
		responseData = append(responseData, map[string]interface{}{
			"number":                 sch.Number,
			"job_application_number": sch.JobApplicationNumber,
			"job_interview_number":   sch.JobInterviewNumber,
			"remarks":                sch.Remarks,
			"comment":                sch.Comment,
			"status":                 sch.Status,
			"application_no":         sch.JobApplication.ApplicationNumber,
			"applicant_name":         sch.JobApplication.FirstName + " " + sch.JobApplication.LastName,
			"applicant_email":        sch.JobApplication.Email,
			"applicant_phone":        sch.JobApplication.PhoneNumber,
			"interview_title":        sch.JobInterview.InterviewTitle,
			"interview_date":         sch.JobInterview.InterviewDate,
			"interview_time":         sch.JobInterview.InterviewTime,
			"location_or_link":       sch.JobInterview.LocationOrLink,
			"created_at":             sch.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}
