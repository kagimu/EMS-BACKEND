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

	config.DB.Where("staff_id_number LIKE ?", prefix+"%").Order("staff_id_number desc").First(&lastDeployment)

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

func ensureStaffUserAndSendWelcomeEmail(applicant models.JobApplication) {
	applicantEmail := strings.TrimSpace(applicant.Email)
	applicantPhone := strings.TrimSpace(applicant.PhoneNumber)
	if applicantEmail == "" {
		return
	}

	var userNumber string
	var contact authModels.Contact
	if applicant.UserNumber != "" {
		userNumber = applicant.UserNumber
	} else if err := config.DB.Where("contact = ?", applicantEmail).First(&contact).Error; err == nil {
		userNumber = contact.UserNumber
	} else {
		userNumber = uuid.NewV4().String()
		userObj := authModels.User{
			Number:     userNumber,
			FirstName:  applicant.FirstName,
			SecondName: applicant.FirstName,
			LastName:   applicant.LastName,
			Gender:     "Other",
			Type:       constants.TypeStaff,
			Status:     "Active",
		}
		if err := config.DB.Omit("Contacts").Create(&userObj).Error; err != nil {
			return
		}
		emailContact := authModels.Contact{
			Contact:     applicantEmail,
			ContactType: "Email",
			UserNumber:  userNumber,
		}
		config.DB.Create(&emailContact)

		if applicantPhone != "" {
			phoneContact := authModels.Contact{
				Contact:     applicantPhone,
				ContactType: "Phone",
				UserNumber:  userNumber,
			}
			config.DB.Create(&phoneContact)
		}
	}

	config.DB.Model(&authModels.User{}).Where("number = ?", userNumber).Update("type", constants.TypeStaff)
	config.DB.Model(&models.JobApplication{}).Where("number = ?", applicant.Number).Update("user_number", userNumber)

	mailInfo := emailSender.Info{
		Email:       applicantEmail,
		MailType:    "WholesomeEd: Congratulations on your Job Deployment!",
		UserName:    applicant.FirstName,
		CurrentYear: time.Now().Year(),
		Type:        "staff-deployed",
	}
	mailInfo.SendEmail()
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
	if err := config.DB.Preload("JobApplication").Where("number = ?", request.JobOfferNumber).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	// 2. Generate new STN ID
	staffIdNumber := generateStaffIdNumber()

	userNumber := ctx.Request.Header.Get("User-Number")
	if userNumber == "" {
		userNumber = jobOffer.UserNumber
	}

	// 3. Create Deployment
	deployment := models.StaffDeployment{
		Number:         uuid.NewV4().String(),
		JobOfferNumber: jobOffer.Number,
		StaffIdNumber:  staffIdNumber,
		Status:         "Deployed",
		UserNumber:     userNumber,
	}

	if err := config.DB.Omit("JobOffer").Create(&deployment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Staff deployment not saved: " + err.Error()})
		return
	}

	// 4. Update JobOffer Status
	config.DB.Model(&models.JobOffer{}).Where("number = ?", jobOffer.Number).Update("status", "DEPLOYED")

	// 5. Ensure user account exists and send confirmation/welcome email
	ensureStaffUserAndSendWelcomeEmail(jobOffer.JobApplication)

	deployment.JobOffer = jobOffer

	ctx.JSON(http.StatusOK, gin.H{"message": "Staff deployed successfully!", "data": deployment})
}

func FetchStaffDeployments(ctx *gin.Context) {
	var deployments []models.StaffDeployment

	schoolNumber := ctx.Query("school_number")

	query := config.DB.Model(&models.StaffDeployment{}).
		Preload("JobOffer").
		Preload("JobOffer.JobApplication").
		Preload("JobOffer.JobInterview").
		Preload("JobOffer.JobInterview.JobVacancy").
		Preload("JobOffer.JobInterview.JobVacancy.JobPosition")

	if schoolNumber != "" {
		query = query.Joins("JOIN job_offers ON job_offers.number = staff_deployments.job_offer_number").
			Joins("JOIN job_applications ON job_applications.number = job_offers.job_application_number").
			Joins("JOIN job_vacancies ON job_vacancies.number = job_applications.job_vacancy_number").
			Where("job_vacancies.school_number = ?", schoolNumber)
	}

	if err := query.Find(&deployments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch staff deployments"})
		return
	}

	var responseData []map[string]interface{}
	for _, dep := range deployments {
		var userStatus string = "No Account"
		var assignedRoles []map[string]interface{}
		var userNumber string
		var contact authModels.Contact
		if err := config.DB.Preload("User").Where("contact = ?", dep.JobOffer.JobApplication.Email).First(&contact).Error; err == nil {
			userNumber = contact.UserNumber
			if contact.User.Status != "" {
				userStatus = contact.User.Status
			} else {
				userStatus = "Active"
			}
		} else if dep.JobOffer.JobApplication.UserNumber != "" {
			var u authModels.User
			if err := config.DB.Where("number = ?", dep.JobOffer.JobApplication.UserNumber).First(&u).Error; err == nil {
				userNumber = u.Number
				userStatus = u.Status
			}
		}

		if userNumber != "" {
			var userRoles []permissionModels.UserRole
			if err := config.DB.Preload("Role").Where("user_number = ?", userNumber).Find(&userRoles).Error; err == nil {
				for _, ur := range userRoles {
					assignedRoles = append(assignedRoles, map[string]interface{}{
						"role_number": ur.RoleNumber,
						"role_name":   ur.Role.Name,
						"expiry_date": ur.ExpiryDate.Format("2006-01-02"),
					})
				}
			}
		}
		if assignedRoles == nil {
			assignedRoles = []map[string]interface{}{}
		}

		responseData = append(responseData, map[string]interface{}{
			"number":           dep.Number,
			"job_offer_number": dep.JobOfferNumber,
			"staff_id_number":  dep.StaffIdNumber,
			"status":           dep.Status,
			"user_status":      userStatus,
			"assigned_roles":   assignedRoles,
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

	if err := config.DB.Model(&deployment).Update("status", request.Status).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Deployment not updated: " + err.Error()})
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

	// Load the deployment with all relations needed to resolve IDs
	var deployment models.StaffDeployment
	if err := config.DB.
		Preload("JobOffer").
		Preload("JobOffer.JobApplication").
		Preload("JobOffer.JobInterview").
		Where("number = ?", id).
		First(&deployment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Deployment not found"})
		return
	}

	jobOfferNumber := deployment.JobOfferNumber
	jobApplicationNumber := deployment.JobOffer.JobApplicationNumber
	jobInterviewNumber := deployment.JobOffer.JobInterviewNumber

	// Resolve the user number: prefer from deployment, fallback to job application
	userNumber := deployment.UserNumber
	if userNumber == "" {
		userNumber = deployment.JobOffer.JobApplication.UserNumber
	}
	// If still empty, try to find via contact email
	if userNumber == "" && deployment.JobOffer.JobApplication.Email != "" {
		var contact authModels.Contact
		if err := config.DB.Where("contact = ?", strings.TrimSpace(deployment.JobOffer.JobApplication.Email)).First(&contact).Error; err == nil {
			userNumber = contact.UserNumber
		}
	}

	// --- Begin single atomic transaction ---
	tx := config.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction: " + tx.Error.Error()})
		return
	}

	// Disable FK checks for this session so delete order never causes constraint violations.
	// Safe here because we remove ALL related rows — no orphans remain after commit.
	if err := tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to disable FK checks: " + err.Error()})
		return
	}

	// Helper to rollback (also re-enables FK checks) and respond on error
	rollbackWithError := func(msg string, err error) {
		tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg + ": " + err.Error()})
	}

	// ── STEP 1: BioData children (leaf tables) ──

	// 1a. StaffBioDataFamilyChild — keyed via family_number
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_family_children
		WHERE family_number IN (
			SELECT number FROM staff_bio_data_families WHERE job_offer_number = ?
		)`, jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data family children", err)
		return
	}

	// 1b. StaffBioDataFamily
	if err := tx.Exec("DELETE FROM staff_bio_data_families WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data family", err)
		return
	}

	// 1c. StaffBioDataEducationItem — keyed via education_number
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_education_items
		WHERE education_number IN (
			SELECT number FROM staff_bio_data_educations WHERE job_offer_number = ?
		)`, jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data education items", err)
		return
	}

	// 1d. StaffBioDataEducation
	if err := tx.Exec("DELETE FROM staff_bio_data_educations WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data education", err)
		return
	}

	// 1e. StaffBioDataLanguageItem — keyed via language_record_number
	if err := tx.Exec(`
		DELETE FROM staff_bio_data_language_items
		WHERE language_record_number IN (
			SELECT number FROM staff_bio_data_languages WHERE job_offer_number = ?
		)`, jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data language items", err)
		return
	}

	// 1f. StaffBioDataLanguage
	if err := tx.Exec("DELETE FROM staff_bio_data_languages WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data language", err)
		return
	}

	// 1g. StaffBioDataContact
	if err := tx.Exec("DELETE FROM staff_bio_data_contacts WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data contact", err)
		return
	}

	// 1h. StaffBioDataMedical
	if err := tx.Exec("DELETE FROM staff_bio_data_medicals WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data medical", err)
		return
	}

	// 1i. StaffBioDataReferences
	if err := tx.Exec("DELETE FROM staff_bio_data_references WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data references", err)
		return
	}

	// 1j. StaffBioDataPersonal
	if err := tx.Exec("DELETE FROM staff_bio_data_personals WHERE job_offer_number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete bio data personal", err)
		return
	}

	// ── STEP 2: JobInterviewSchedule (all schedules for this application) ──
	if err := tx.Exec("DELETE FROM job_interview_schedules WHERE job_application_number = ?", jobApplicationNumber).Error; err != nil {
		rollbackWithError("Failed to delete job interview schedules", err)
		return
	}

	// ── STEP 3: StaffDeployment itself ──
	if err := tx.Exec("DELETE FROM staff_deployments WHERE number = ?", id).Error; err != nil {
		rollbackWithError("Failed to delete staff deployment", err)
		return
	}

	// ── STEP 4: JobOffer ──
	if err := tx.Exec("DELETE FROM job_offers WHERE number = ?", jobOfferNumber).Error; err != nil {
		rollbackWithError("Failed to delete job offer", err)
		return
	}

	// ── STEP 5: JobApplication ──
	if err := tx.Exec("DELETE FROM job_applications WHERE number = ?", jobApplicationNumber).Error; err != nil {
		rollbackWithError("Failed to delete job application", err)
		return
	}

	// ── STEP 6: JobInterview ──
	if jobInterviewNumber != "" {
		if err := tx.Exec("DELETE FROM job_interviews WHERE number = ?", jobInterviewNumber).Error; err != nil {
			rollbackWithError("Failed to delete job interview", err)
			return
		}
	}

	// ── STEP 7: User-related records (only if a user was created for this staff) ──
	if userNumber != "" {
		// 7a. UserRoles (permissions)
		if err := tx.Exec("DELETE FROM user_roles WHERE user_number = ?", userNumber).Error; err != nil {
			rollbackWithError("Failed to delete user roles", err)
			return
		}

		// 7b. EntityUsers (EOS)
		if err := tx.Exec("DELETE FROM entity_users WHERE user_number = ?", userNumber).Error; err != nil {
			rollbackWithError("Failed to delete entity users", err)
			return
		}

		// 7c. PasswordResets
		if err := tx.Exec("DELETE FROM password_resets WHERE user_number = ?", userNumber).Error; err != nil {
			rollbackWithError("Failed to delete password resets", err)
			return
		}

		// 7d. UserPasswords
		if err := tx.Exec("DELETE FROM user_passwords WHERE user_number = ?", userNumber).Error; err != nil {
			rollbackWithError("Failed to delete user passwords", err)
			return
		}

		// 7e. Contacts
		if err := tx.Exec("DELETE FROM contacts WHERE user_number = ?", userNumber).Error; err != nil {
			rollbackWithError("Failed to delete user contacts", err)
			return
		}

		// 7f. User row — only delete if user type is Staff to protect shared accounts
		var user authModels.User
		if err := tx.Where("number = ?", userNumber).First(&user).Error; err == nil {
			if user.Type == "Staff" {
				if err := tx.Exec("DELETE FROM users WHERE number = ?", userNumber).Error; err != nil {
					rollbackWithError("Failed to delete user", err)
					return
				}
			}
		}
	}

	// Re-enable FK checks before committing
	if err := tx.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to re-enable FK checks: " + err.Error()})
		return
	}

	// ── COMMIT ──
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit deletion transaction: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Staff deployment and all related records permanently deleted."})
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

	// Check JobApplication.UserNumber or look up user by email in contacts table first to get their true user number
	var contact authModels.Contact
	if deployment.JobOffer.JobApplication.UserNumber != "" {
		userNumber = deployment.JobOffer.JobApplication.UserNumber
	} else if err := tx.Where("contact = ?", applicantEmail).First(&contact).Error; err == nil {
		userNumber = contact.UserNumber
	} else {
		// Create new user profile
		userNumber = uuid.NewV4().String()
		userObj := authModels.User{
			Number:     userNumber,
			FirstName:  applicantFirstName,
			SecondName: applicantFirstName,
			LastName:   applicantLastName,
			Gender:     "Other",
			Type:       constants.TypeStaff,
			Status:     "Active",
		}
		if err := tx.Omit("Contacts").Create(&userObj).Error; err != nil {
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

	// Check if user has a password set already. If not, generate password reset token and send activation/welcome email
	var up authModels.UserPassword
	if err := tx.Where("user_number = ?", userNumber).First(&up).Error; err != nil {
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
	}

	var assignedRoleNumbers []string
	for _, roleItem := range request.Roles {
		date, err := commons.ConvertStringToDate(roleItem.ExpiryDate)
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format for role expiry date"})
			return
		}
		assignedRoleNumbers = append(assignedRoleNumbers, roleItem.RoleNumber)
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

	if len(assignedRoleNumbers) > 0 {
		tx.Where("user_number = ? AND role_number NOT IN (?)", userNumber, assignedRoleNumbers).Delete(&permissionModels.UserRole{})
	} else {
		tx.Where("user_number = ?", userNumber).Delete(&permissionModels.UserRole{})
	}

	if request.UserStatus != "" {
		tx.Model(&authModels.User{}).Where("number = ?", userNumber).Update("status", request.UserStatus)
	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Staff role and access details updated successfully!"})
}

