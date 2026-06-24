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

func generateStaffIdNumber() string {
	var lastDeployment models.StaffDeployment
	year := time.Now().Year()
	prefix := fmt.Sprintf("STN%d", year)

	config.DB.Where("staff_id_number LIKE ?", prefix+"%").Order("created_at desc").First(&lastDeployment)

	if lastDeployment.StaffIdNumber == "" {
		return fmt.Sprintf("%s001", prefix)
	}

	numberPart := lastDeployment.StaffIdNumber[len(prefix):]
	sequence, err := strconv.Atoi(numberPart)
	if err != nil {
		sequence = 0
	}
	sequence++

	return fmt.Sprintf("%s%03d", prefix, sequence)
}

func CreateStaffDeployment(ctx *gin.Context) {
	var request requests.StaffDeploymentRequest
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

	// 1. Verify Job Offer exists
	var jobOffer models.JobOffer
	if err := config.DB.Where("number = ?", request.JobOfferNumber).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	// 2. Generate new STN ID
	staffIdNumber := generateStaffIdNumber()

	// 3. Create Deployment
	deployment := models.StaffDeployment{
		Number:         uuid.NewV4().String(),
		JobOfferNumber: jobOffer.Number,
		StaffIdNumber:  staffIdNumber,
		Status:         "Deployed",
		UserNumber:     ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Staff deployment not saved."})
		return
	}

	// 4. Update JobOffer Status
	jobOffer.Status = "DEPLOYED"
	config.DB.Save(&jobOffer)

	ctx.JSON(http.StatusOK, gin.H{"message": "Staff deployed successfully!", "data": deployment})
}

func FetchStaffDeployments(ctx *gin.Context) {
	var deployments []models.StaffDeployment

	query := config.DB.Model(&models.StaffDeployment{}).
		Preload("JobOffer").
		Preload("JobOffer.JobApplication").
		Preload("JobOffer.JobInterview").
		Preload("JobOffer.JobInterview.JobVacancy").
		Preload("JobOffer.JobInterview.JobVacancy.JobPosition")

	if err := query.Find(&deployments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch staff deployments"})
		return
	}

	var responseData []map[string]interface{}
	for _, dep := range deployments {
		responseData = append(responseData, map[string]interface{}{
			"number":           dep.Number,
			"job_offer_number": dep.JobOfferNumber,
			"staff_id_number":  dep.StaffIdNumber,
			"status":           dep.Status,
			"applicant_name":   dep.JobOffer.JobApplication.FirstName + " " + dep.JobOffer.JobApplication.LastName,
			"applicant_email":  dep.JobOffer.JobApplication.Email,
			"applicant_phone":  dep.JobOffer.JobApplication.PhoneNumber,
			"job_position":     dep.JobOffer.JobInterview.JobVacancy.JobPosition.PostName,
			"created_at":       dep.CreatedAt,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func EditStaffDeployment(ctx *gin.Context) {
	var request requests.EditStaffDeploymentRequest
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

	var deployment models.StaffDeployment
	if err := config.DB.Where("number = ?", request.Number).First(&deployment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Staff deployment not found"})
		return
	}

	deployment.Status = request.Status

	if err := config.DB.Save(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Deployment not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deployment updated successfully!"})
}

func DeleteStaffDeployment(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deployment number is required"})
		return
	}

	var deployment models.StaffDeployment
	if err := config.DB.Where("number = ?", id).First(&deployment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Deployment not found"})
		return
	}

	if err := config.DB.Delete(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Deployment not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deployment deleted successfully!"})
}
