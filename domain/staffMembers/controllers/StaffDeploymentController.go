package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"


	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
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

func AssignStaffRole(ctx *gin.Context) {
	var request requests.AssignStaffRoleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx := config.DB.Begin()

	var deployment models.StaffDeployment
	if err := tx.Preload("JobOffer.JobApplication").Where("number = ?", request.DeploymentNumber).First(&deployment).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Staff deployment not found"})
		return
	}

	var userNumber string
	var applicantEmail = strings.TrimSpace(deployment.JobOffer.JobApplication.Email)
	var applicantPhone = strings.TrimSpace(deployment.JobOffer.JobApplication.PhoneNumber)
	var applicantFirstName = deployment.JobOffer.JobApplication.FirstName
	var applicantLastName = deployment.JobOffer.JobApplication.LastName

	if applicantEmail == "" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Staff applicant email is empty, cannot send activation email"})
		return
	}

	// Look up user by main email in contacts table first to get their true user number (avoiding admin ID)
	var contact authModels.Contact
	if err := tx.Where("contact = ?", applicantEmail).First(&contact).Error; err == nil {
		userNumber = contact.UserNumber
	} else {
		// Create new user profile
		userNumber = uuid.NewV4().String()
		userObj := authModels.User{
			Number:     userNumber,
			FirstName:  applicantFirstName,
			SecondName: applicantFirstName,
			LastName:   applicantLastName,
			Type:       constants.TypeStaff,
			Status:     "Active",
		}
		if err := tx.Create(&userObj).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create staff user profile"})
			return
		}
		// Save Email contact
		emailContact := authModels.Contact{
			Contact:     applicantEmail,
			ContactType: "Email",
			UserNumber:  userNumber,
		}
		if err := tx.Create(&emailContact).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save email contact"})
			return
		}
		// Save Phone contact if present
		if applicantPhone != "" {
			phoneContact := authModels.Contact{
				Contact:     applicantPhone,
				ContactType: "Phone",
				UserNumber:  userNumber,
			}
			if err := tx.Create(&phoneContact).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save phone contact"})
				return
			}
		}
		// Attach to entity if available
		entityNumber, _ := commons.GetUserEntityNumber(ctx)
		if entityNumber != "" {
			entityUser := eosModels.EntityUser{
				EntityNumber:  entityNumber,
				UserNumber:    userNumber,
				AddedByNumber: ctx.GetHeader("User-Number"),
			}
			tx.Create(&entityUser)
		}
	}

	// Ensure user type is set to Staff and link to JobApplication
	tx.Model(&authModels.User{}).Where("number = ?", userNumber).Update("type", constants.TypeStaff)
	tx.Model(&deployment.JobOffer.JobApplication).Update("user_number", userNumber)

	// Always generate password reset token and send activation/welcome email to main email of deployed staff
	resetToken := uuid.NewV4().String()
	pwReset := authModels.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
		Token:      resetToken,
		Status:     "Inactive",
	}
	if err := tx.Create(&pwReset).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create password reset token"})
		return
	}

	// Send welcome/activation email to the main email of the deployed staff
	mailInfo := emailSender.Info{
		Email:       applicantEmail,
		Token:       resetToken,
		MailType:    "WholesomeEd: Staff Account Confirmation",
		UserName:    applicantFirstName,
		CurrentYear: time.Now().Year(),
		Type:        "registration",
	}
	if err := mailInfo.SendEmail(); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to send account activation email: " + err.Error()})
		return
	}


	for _, roleItem := range request.Roles {
		date, err := commons.ConvertStringToDate(roleItem.ExpiryDate)
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format for role expiry date"})
			return
		}
		var existingRole permissionModels.UserRole
		err = tx.Where("user_number = ? AND role_number = ?", userNumber, roleItem.RoleNumber).First(&existingRole).Error
		if err == nil {
			existingRole.ExpiryDate = date
			if err := tx.Save(&existingRole).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role assignment"})
				return
			}
		} else {
			newRole := permissionModels.UserRole{
				UserNumber: userNumber,
				RoleNumber: roleItem.RoleNumber,
				AddedBy:    ctx.Request.Header.Get("User-Number"),
				ExpiryDate: date,
			}
			if err := tx.Create(&newRole).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign role"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Staff role assigned successfully and account activated if new!"})
}

