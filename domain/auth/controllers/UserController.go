package controllers

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/requests"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	uuid "github.com/satori/go.uuid"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func FetchUsers(c *gin.Context) {

	entityNumber, _ := commons.GetUserEntityNumber(c)

	var paginatedResponse = pagination.PaginatedResponse{}
	var users []models.User
	if err := config.DB.Scopes(
		filters.ByStatus(c.Query("status")),
		pagination.Paginate(c, &paginatedResponse),
		filters.ByAnyField("entity_users.entity_number", entityNumber),
	).Joins("JOIN entity_users ON entity_users.user_number = users.number").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	type roleObject struct {
		RoleName   string `json:"role_name"`
		RoleNumber string `json:"role_number"`
		ExpiryDate string `json:"expiry_date"`
		Status     string `json:"status"`
	}

	type CLeanUser struct {
		models.User
		Roles []roleObject `json:"roles"`
	}

	var cleanUsers = make([]CLeanUser, 0)

	for _, user := range users {
		var userRoles []permissionModels.UserRole

		var roles []roleObject
		if err := config.DB.Preload("Role").Where("user_number = ?", user.Number).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user roles"})
			return
		}

		for _, userRole := range userRoles {

			var status string
			if userRole.ExpiryDate.Before(time.Now()) {
				status = "Expired"
			} else {
				status = "Active"
			}

			roles = append(roles, roleObject{RoleName: userRole.Role.Name,
				RoleNumber: userRole.Role.Number, ExpiryDate: userRole.ExpiryDate.Format("2006-01-02"), Status: status})
		}

		cleanUsers = append(cleanUsers, CLeanUser{User: user, Roles: roles})

	}

	paginatedResponse.Content = cleanUsers

	c.JSON(http.StatusOK, paginatedResponse)
}
func FetchAllUsers(c *gin.Context) {

	//var paginatedResponse = pagination.PaginatedResponse{}
	var users []models.User
	if err := config.DB.Scopes(
		filters.ByStatus(c.Query("status")),
	).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}

	type CLeanUser struct {
		models.User
		Roles []string `json:"roles"`
	}

	var cleanUsers = make([]CLeanUser, 0)

	for _, user := range users {
		var userRoles []permissionModels.UserRole
		var roles []string
		if err := config.DB.Preload("Role").Where("user_number = ?", user.Number).Find(&userRoles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user roles"})
			return
		}

		for _, userRole := range userRoles {
			roles = append(roles, userRole.Role.Name)
		}

		cleanUsers = append(cleanUsers, CLeanUser{User: user, Roles: roles})

	}

	c.JSON(http.StatusOK, cleanUsers)
}

func FetchUser(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	entityNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || entityNumber == "" {
		c.JSON(http.StatusConflict, gin.H{"message": "No active organisation is selected"})
		return
	}

	var user models.User
	if err := config.DB.Where("number = ?", c.Param("user_number")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// A user outside the actor's active organisation is reported as not found
	// so its existence is not disclosed.
	var membership eosModels.EntityUser
	if err := config.DB.Where("user_number = ? AND entity_number = ?", user.Number, entityNumber).
		First(&membership).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch user"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func actorCanAssignRole(actor *models.User, organisation, roleNumber string) (bool, error) {
	var rows []permissionModels.RolePermission
	if err := config.DB.Where("role_number = ?", roleNumber).Find(&rows).Error; err != nil {
		return false, err
	}

	rolePermissions := make([]string, 0, len(rows))
	for _, row := range rows {
		rolePermissions = append(rolePermissions, row.Permission)
	}

	if commons.ContainsControlPlanePermission(rolePermissions) && !commons.IsSoftwareOwner(actor.Number) {
		return false, nil
	}
	if actor.Type == "Admin" {
		return true, nil
	}

	granted, err := repository.FetchPermissionsByUserAndEntity(actor.Number, organisation)
	if err != nil {
		return false, err
	}
	allowed := make(map[string]struct{}, len(granted))
	for _, permission := range granted {
		allowed[permission] = struct{}{}
	}
	for _, permission := range rolePermissions {
		if _, ok := allowed[permission]; !ok {
			return false, nil
		}
	}
	return true, nil
}

func CreateUser(c *gin.Context) {

	var user requests.StaffAccountRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if err := user.ValidateStaffAccountRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	actor, err := commons.GetUserFromHeader(c)
	if err != nil || actor == nil || actor.Number == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	activeEntity, err := commons.GetUserEntityNumber(c)
	if err != nil || activeEntity == "" {
		c.JSON(http.StatusConflict, gin.H{"message": "No active organisation is selected"})
		return
	}

	if user.EntityNumber == "" {
		user.EntityNumber = activeEntity
	}
	if user.EntityNumber != activeEntity {
		c.JSON(http.StatusForbidden, gin.H{"message": "You can only add users to your current organisation"})
		return
	}

	for _, role := range user.Roles {
		var roleModel permissionModels.Role
		if err := config.DB.Where("number = ?", role.RoleNumber).First(&roleModel).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Role not found"})
			return
		}
		if roleModel.EntityNumber != activeEntity {
			c.JSON(http.StatusForbidden, gin.H{"message": "You can only assign roles from your current organisation"})
			return
		}

		allowed, err := actorCanAssignRole(actor, activeEntity, role.RoleNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify role authority"})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"message": "You cannot assign a role with permissions you do not hold"})
			return
		}
	}

	tx := config.DB.Begin()
	var userObject = models.User{
		Number:     uuid.NewV4().String(),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Gender:     user.Gender,
		Type:       constants.TypeStaff,
		Status:     "Active",
	}

	if err := tx.Create(&userObject).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// check if email or phone already exists

	if err := tx.Where("contact = ?", user.Email).First(&models.Contact{}).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "This Email was already used to create an account"})
		return
	}

	if err := tx.Where("contact = ?", user.Phone).First(&models.Contact{}).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "This Phone was already used to create an account"})
		return
	}

	//after you need to save user email
	saveEmailContact := models.Contact{
		Contact:     user.Email,
		ContactType: "Email",
		UserNumber:  userObject.Number,
	}

	if err := tx.Create(&saveEmailContact).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	savePhoneContact := models.Contact{
		Contact:     user.Phone,
		ContactType: "Phone",
		UserNumber:  userObject.Number,
	}
	if err := tx.Create(&savePhoneContact).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// append roles

	for _, role := range user.Roles {
		var userRole permissionModels.UserRole
		var roleModel permissionModels.Role
		if err := tx.Where("number = ?", role.RoleNumber).First(&roleModel).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch role"})
			return
		}

		userRole.UserNumber = userObject.Number
		userRole.RoleNumber = role.RoleNumber
		if role.ExpiryDate == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Expiry date is required"})
			return

		}
		date, err := commons.ConvertStringToDate(role.ExpiryDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to append role to user" + err.Error()})
			return
		}
		userRole.ExpiryDate = date
		userRole.AddedBy = actor.Number

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create role to user"})
			return

		}
	}

	// set password
	//hashPassword, hashErr := hassPassTokenGen.Hash(user.Password)
	//if hashErr != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to encrypt password"})
	//	return
	//}
	//
	//var userPassword = models.UserPassword{
	//	UserNumber:   userObject.Number,
	//	UserPassword: hashPassword,
	//}
	//
	//if err := tx.Create(&userPassword).Error; err != nil {
	//	tx.Rollback()
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user" + err.Error()})
	//	return
	//}

	// attach user to this entity
	if user.EntityNumber != "" {
		if err := licensingServices.EnsureMembershipSeat(tx, actor.Number, user.EntityNumber, userObject.Number); err != nil {
			tx.Rollback()
			if errors.Is(err, licensingServices.ErrSeatCapacityReached) || errors.Is(err, licensingServices.ErrSeatScopeInactive) {
				c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reserve a user seat"})
			}
			return
		}
		var entityUser = eosModels.EntityUser{
			EntityNumber:  user.EntityNumber,
			UserNumber:    userObject.Number,
			AddedByNumber: actor.Number,
		}

		if err := tx.Create(&entityUser).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user" + err.Error()})
			return
		}

	}

	//create / generate a token for the user to reset their password
	//userToken := tokenGenerator()
	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userObject.Number,
		Token:      userToken,
		Status:     "Inactive",
	}

	if err := tx.Create(&saveUserPasswordToken).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to create token"})
		return
	}
	//email sender
	Token := userToken
	d := emailSender.Info{Email: user.Email, Token: Token, MailType: "WholesomeEd: User Account Confirmation", UserName: userObject.SecondName, CurrentYear: time.Now().Year(), Type: "registration"}
	err = d.SendEmail()
	if err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//fmt.Println(fmt.Printf("Sending mail error %s\n", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to send email\n" + err.Error()})
		return
		//fmt.Printf("Sending mail error %s\n", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Failed :"+err.Error()))
	}

	c.JSON(http.StatusCreated, commons.Feedback("User created successfully"))
}

func CreateUserAcc(c *gin.Context) {
	userRepo := repository.NewSqlUserRepository(config.DB)
	fmt.Println(userRepo)
	//userService := services.NewUserService(userRepo)
}

// FetchPendingUsers returns a paginated list of Staff users who have registered
// but have not yet activated their account (status = "Inactive").
// Results are scoped to the calling user's entity.
// Supports ?page, ?per_page, and ?search query parameters.
func FetchPendingUsers(c *gin.Context) {

	schoolNumber := c.Query("school_number")
	search := c.Query("search")

	var paginatedResponse = pagination.PaginatedResponse{}
	var users []models.User

	query := config.DB.Model(&models.User{}).Scopes(
		filters.ByStatus("Inactive"),
		filters.ByType("JobApplicant"),
		filters.BySearchingAnyField([]string{"users.first_name", "users.second_name", "users.last_name"}, search),
	)

	if schoolNumber != "" {
		query = query.Where(
			"users.number IN (SELECT user_number FROM entity_users WHERE entity_number = ?)", schoolNumber,
		)
	}

	// Apply pagination AFTER all filters so the count is accurate
	if err := query.Scopes(pagination.Paginate(c, &paginatedResponse)).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch pending users"})
		return
	}

	type PendingUser struct {
		models.User
		Email     string `json:"email"`
	}

	var result = make([]PendingUser, 0)

	for _, user := range users {
		var emailContact models.Contact
		email := ""
		if err := config.DB.
			Where("user_number = ? AND contact_type = ?", user.Number, "Email").
			First(&emailContact).Error; err == nil {
			email = emailContact.Contact
		}
		result = append(result, PendingUser{User: user, Email: email})
	}

	paginatedResponse.Content = result
	c.JSON(http.StatusOK, paginatedResponse)
}

// ResendActivationEmail resends the account activation email to a pending user.
// It reuses the existing Inactive token from password_resets, or creates a new one
// if none exists. The user must be in "Inactive" status.
func ResendActivationEmail(c *gin.Context) {

	type ResendRequest struct {
		UserNumber string `json:"user_number" binding:"required"`
	}

	var req ResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_number is required"})
		return
	}

	// Fetch user and verify they are pending
	var user models.User
	if err := config.DB.
		Where("number = ?", req.UserNumber).
		Preload("Contacts").
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if user.Status != "Inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User account is already active"})
		return
	}

	// Resolve the user's email from contacts
	userEmail := ""
	for _, contact := range user.Contacts {
		if contact.ContactType == "Email" {
			userEmail = contact.Contact
			break
		}
	}
	if userEmail == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No email address found for this user"})
		return
	}

	// Fetch or create activation token
	var passwordReset models.PasswordReset
	tokenErr := config.DB.
		Where("user_number = ? AND status = ?", req.UserNumber, "Inactive").
		First(&passwordReset).Error

	if tokenErr != nil {
		// No existing token — create a fresh one inside a transaction
		newToken := uuid.NewV4().String()
		passwordReset = models.PasswordReset{
			Number:     uuid.NewV4().String(),
			UserNumber: req.UserNumber,
			Token:      newToken,
			Status:     "Inactive",
		}
		if err := config.DB.Create(&passwordReset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate activation token"})
			return
		}
	}

	// Send activation email
	d := emailSender.Info{
		Email:       userEmail,
		Token:       passwordReset.Token,
		MailType:    "WholesomeEd: User Account Confirmation",
		UserName:    user.SecondName,
		CurrentYear: time.Now().Year(),
		Type:        "registration",
	}
	if err := d.SendEmail(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send activation email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activation email sent successfully"})
}

// DeletePendingUser permanently removes a pending (Inactive) user and all their
// associated data (password_resets, contacts) from the database.
// Only users with status "Inactive" can be deleted through this endpoint.
func DeletePendingUser(c *gin.Context) {

	userNumber := c.Param("user_number")
	if userNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_number is required"})
		return
	}

	// Fetch and validate the user
	var user models.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if user.Status != "Inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only pending (Inactive) users can be deleted through this endpoint"})
		return
	}

	tx := config.DB.Begin()

	// Delete password reset tokens
	if err := tx.Where("user_number = ?", userNumber).Delete(&models.PasswordReset{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove user token data"})
		return
	}

	// Delete contacts (email, phone)
	if err := tx.Where("user_number = ?", userNumber).Delete(&models.Contact{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove user contact data"})
		return
	}

	// Delete the user record (hard delete — bypass GORM soft-delete)
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pending user deleted successfully"})
}
