package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	rolesResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/responses"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/hassPassTokenGen"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func Login(context *gin.Context) {

	type LoginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var login LoginUser

	err := context.BindJSON(&login)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	Email := login.Email
	Password := login.Password

	if Email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
		return
	}

	if Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Password is required!"})
		return
	}

	Email = strings.TrimSpace(strings.ToLower(Email))

	//email validation
	var rxEmail = regexp.MustCompile(".+@.+\\..+")
	match := rxEmail.Match([]byte(Email))
	if match == false {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email!"})
		return
	}

	//	check if user exists
	var contact models.Contact

	userContact := config.DB.
		Where("contact = ?", Email).
		First(&contact)
	if userContact.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Email not found!"})
		return
	}

	userNumber := contact.UserNumber

	//get user password from the password table
	var userPassword models.UserPassword
	getUserPassword := config.DB.
		Where("user_number = ?", userNumber).
		First(&userPassword)
	if getUserPassword.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Password!"})
		return
	}

	//check if user is active / not blocked
	var user models.User
	config.DB.Where("number = ?", userNumber).First(&user)
	if user.Status != "Active" {
		context.JSON(http.StatusNotFound, gin.H{"message": "Account not Active, please reset password to activate your account!"})
		return
	}

	userHashPassword := userPassword.UserPassword

	//	we need to validate hashed password
	if !hassPassTokenGen.CheckHash(Password, userHashPassword) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Password!"})
		return
	}

	//return User profile
	var userProfile models.User
	if err := config.DB.
		Where("number = ?", userNumber).
		Preload("Contacts").
		First(&userProfile).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	var contactListInfo []models.ContactList

	for _, item := range userProfile.Contacts {
		contactListInfo = append(contactListInfo, models.ContactList{
			ID:          item.ID,
			Contact:     item.Contact,
			ContactType: item.ContactType,
		})
	}

	//prepare roles and permissions
	//get user roles
	var userRoles []permissionModels.UserRole
	var userRolesList []rolesResponses.RoleView
	var permissions []string

	config.DB.Where("user_number = ? AND DATE(expiry_date) >= ?", userNumber, time.Now()).Find(&userRoles)

	for _, item := range userRoles {

		userRolesList = append(userRolesList, rolesResponses.RoleView{
			Number:      item.Role.Number,
			Name:        item.Role.Name,
			Description: item.Role.Description,
		})

		//get role permissions
		var rolePermissions []permissionModels.RolePermission
		config.DB.Where("role_number = ?", item.RoleNumber).Find(&rolePermissions)

		for _, permission := range rolePermissions {
			permissions = append(permissions, permission.Permission)
		}

	}

	scopes := filters.GetUserScopedItems(context, "stream", userNumber)

	scopedCampuses := filters.GetUserScopedItems(context, "campus", userNumber)
	scopedSchools := filters.GetUserScopedItems(context, "school", userNumber)
	campuses := make([]responses.AuthCampusResponse, 0)
	schools := make([]responses.ModelResponse, 0)

	for _, school := range scopedSchools {
		var schoolResponse responses.ModelResponse
		school, err := repositories2.GetSchoolByNumber(school)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": "School not found"})
			return
		}
		schoolResponse.Name = school.Name
		schoolResponse.Number = school.Number
		schools = append(schools, schoolResponse)
	}

	for _, campus := range scopedCampuses {

		var campusResponse responses.AuthCampusResponse
		campus, err := repositories.GetCampusByNumber(campus)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": "Campus not found"})
			return
		}
		campusResponse.CampusName = campus.Name
		classes, err := services2.FetchCampusClassesByCampusNumber(campus.Number)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": "Classes not found"})
			return
		}
		for _, class := range classes {
			campusResponse.Classes = append(campusResponse.Classes, responses2.ClassResponse{
				Name:   class.ClassName,
				Number: class.Number,
			})
		}

		campuses = append(campuses, campusResponse)

	}

	userView := responses.UserView{
		Number:      userProfile.Number,
		FirstName:   userProfile.FirstName,
		SecondName:  userProfile.SecondName,
		LastName:    userProfile.LastName,
		Gender:      userProfile.Gender,
		Type:        userProfile.Type,
		Roles:       userRolesList,
		Permissions: permissions,
		Scopes:      scopes,
		Campuses:    campuses,
		Schools:     schools,
	}

	token, err := services.JWTAuthService().GenerateToken(user.Number)

	context.JSON(http.StatusOK, gin.H{"message": "Login Successful", "token": token, "data": userView})

}

func Profile(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var userProfile models.User

	if err := config.DB.
		Where("number = ?", user.Number).
		Preload("Contacts").
		First(&userProfile).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	var contactListInfo []models.ContactList

	for _, item := range userProfile.Contacts {
		contactListInfo = append(contactListInfo, models.ContactList{
			ID:          item.ID,
			Contact:     item.Contact,
			ContactType: item.ContactType,
			//UserNumber:  item.UserNumber,
		})
	}

	_userProfile := responses.GetUserList{
		Number:     userProfile.Number,
		FirstName:  userProfile.FirstName,
		SecondName: userProfile.SecondName,
		LastName:   userProfile.LastName,
		Gender:     userProfile.Gender,
		Contacts:   contactListInfo,
		Type:       userProfile.Type,
		Status:     userProfile.Status,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _userProfile})
}

func Register(context *gin.Context) {

	//	we need to first initialize a struct that caters for posted data
	//then after this data is separated

	var userPostData requests.UserDataRequest
	err := context.BindJSON(&userPostData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	userPostData.Prepare()

	if err := userPostData.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	firstName := userPostData.FirstName
	secondName := userPostData.SecondName
	lastName := userPostData.LastName
	userGender := userPostData.Gender
	userEmail := userPostData.Email
	userPhone := userPostData.Phone

	//remove white space from email and converts to lower case
	userEmail = strings.TrimSpace(strings.ToLower(userEmail))

	//email validation
	var rxEmail = regexp.MustCompile(".+@.+\\..+")
	match := rxEmail.Match([]byte(userEmail))
	if match == false {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email!"})
		return
	}

	//	duplicate Email validation
	var contact models.Contact
	if duplicateEmail := config.DB.
		Where("contact = ?", userEmail).
		Where("contact_type = ?", "email").
		First(&contact).Error; duplicateEmail == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
		return
	}

	userPhone = commons.CheckValidPhone(strings.TrimSpace(userPhone))
	if userPhone == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Phone:" + userPhone})
		return
	}

	//	duplicate Phone validation
	if duplicatePhone := config.DB.
		Where("contact = ?", userPhone).
		Where("contact_type = ?", "phone").
		First(&contact).Error; duplicatePhone == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Phone already exists"})
		return
	}

	///save user data

	//first time registration user type must be admin
	var user models.User
	userType := "User"
	userStatus := "Inactive"

	if userTypeChecker := config.DB.
		First(&user).Error; userTypeChecker != nil {
		userType = constants.TypeAdmin
	}

	//creates a secure uuid user number

	tx := config.DB.Begin()
	userNumber := uuid.NewV4().String()
	///save user data
	saveUser := models.User{
		Number:     userNumber,
		FirstName:  firstName,
		SecondName: secondName,
		LastName:   lastName,
		Gender:     userGender,
		Type:       userType,
		Status:     userStatus,
	}

	if err := tx.Create(&saveUser).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create user"})
		return
	}

	//after you need to save user email
	saveEmailContact := models.Contact{
		Contact:     userEmail,
		ContactType: "Email",
		UserNumber:  saveUser.Number,
	}

	if err := tx.Create(&saveEmailContact).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create user email " + err.Error()})
		return
	}

	//after you need to save user phone
	savePhoneContact := models.Contact{
		Contact:     userPhone,
		ContactType: "Phone",
		UserNumber:  saveUser.Number,
	}

	if err := tx.Create(&savePhoneContact).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create user phone " + err.Error()})
		return
	}

	//create / generate a token for the user to reset their password
	//userToken := tokenGenerator()
	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
		Token:      userToken,
		Status:     userStatus,
	}

	if err := tx.Create(&saveUserPasswordToken).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create token"})
		return
	}

	//email sender
	Token := userToken
	d := emailSender.Info{Email: userEmail, Token: Token, MailType: "WholesomeEd: User Account Confirmation", UserName: secondName, CurrentYear: time.Now().Year(), Type: "registration"}
	err = d.SendEmail()
	if err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//fmt.Println(fmt.Printf("Sending mail error %s\n", err.Error()))
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to send email\n" + err.Error()})
		return
		//fmt.Printf("Sending mail error %s\n", err.Error())
	}

	var collection []interface{}
	collection = append(collection, saveEmailContact)
	collection = append(collection, savePhoneContact)

	type UserInfo struct {
		models.User
		Contacts []interface{}
	}

	UserDataInfo := UserInfo{
		saveUser,
		collection,
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Registered Successfully!",
		"data":      UserDataInfo,
		"tokenInfo": saveUserPasswordToken.Token,
	})

}

func ResetPassword(context *gin.Context) {

	var userMail requests.UserEmail

	err := context.BindJSON(&userMail)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	userMail.Prepare()

	if err := userMail.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	Email := userMail.Email

	Email = strings.TrimSpace(strings.ToLower(Email))

	//email validation
	var rxEmail = regexp.MustCompile(".+@.+\\..+")
	match := rxEmail.Match([]byte(Email))
	if match == false {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Email!"})
		return
	}

	//	check if user exists
	var contact models.Contact

	userContact := config.DB.
		Where("contact = ?", Email).
		Joins("join users on users.number = contacts.user_number").
		Preload("User").
		First(&contact)
	if userContact.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Email not found!"})
		return
	}

	userNumber := contact.UserNumber
	userStatus := contact.User.Status
	userToken := ""

	//passwordStatus := contact.User.Status

	status := "Inactive"
	statusActive := "Active"
	var userData models.User

	//this is done only when user already sent an activation link
	tx := config.DB.Begin()
	if userStatus == statusActive {
		//then we proceed to update the user
		//var userData models.User
		updateStatus := tx.Model(&userData).
			Where("number = ?", userNumber).
			Update("status", status)
		if updateStatus.Error != nil {
			tx.Rollback()
			context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
			return
		}

		userToken = uuid.NewV4().String()
		saveUserPasswordToken := models.PasswordReset{
			Number:     uuid.NewV4().String(),
			UserNumber: userNumber,
			Token:      userToken,
			Status:     status,
		}

		if err := tx.Create(&saveUserPasswordToken).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create Learning area"})
			return
		}

	} else {

		var passwordResetInfo models.PasswordReset
		if err := tx.
			Where("user_number = ?", userNumber).
			Where("status = ?", status).
			First(&passwordResetInfo).Error; err != nil {
			tx.Rollback()
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Token not found"})
			return
		}

		userToken = passwordResetInfo.Token

	}

	//email sender
	d := emailSender.Info{Email: Email, Token: userToken, MailType: "Wholesome: User Account Confirmation", UserName: contact.User.SecondName, CurrentYear: time.Now().Year(), Type: "reset-password"}
	err = d.SendEmail()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Error " + err.Error()})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create reset password"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Please check your email to continue resetting your password", "token": userToken})
}

func UserPassword(context *gin.Context) {

	var userPass requests.UserPassword
	err := context.BindJSON(&userPass)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := userPass.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token := userPass.Token
	userPassword := userPass.Password
	//confirmPassword := userPass.ConfirmPassword

	//	validate if token is provided
	/*	if token == "" {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Token is required!"})
			return
		}

		if userPassword == "" {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Password is required!"})
			return
		}

		if userPassword != confirmPassword {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Password Mismatch!"})
			return
		}*/

	//	check if token is valid for the user to set the password
	//this we need to fetch a token which is active
	var passwordReset models.PasswordReset
	checkUserToken := config.DB.
		Where("token = ?", token).
		First(&passwordReset)
	if checkUserToken.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token!"})
		return
	}

	if passwordReset.Status == "Active" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Token has expired! Reset Password, to get new password!"})
		return
	}

	//this gets user number to validate
	userNumber := passwordReset.UserNumber

	//before saving password we need to first hash the password
	hashPassword, hashErr := hassPassTokenGen.Hash(userPassword)
	if hashErr != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid password"})
		return
	}

	tx := config.DB.Begin()
	var userData models.User

	//check if user is a valid user
	validateUser := tx.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	//if user already exist the user password should be updated
	var userPasswordInfo models.UserPassword

	getUserPassword := tx.
		Where("user_number = ?", userNumber).
		First(&userPasswordInfo)
	if getUserPassword.Error == nil {
		//this will just update user password data
		updateStatus := tx.Model(&userPasswordInfo).
			Where("user_number = ?", userNumber).
			Update("user_password", hashPassword)
		if updateStatus.Error != nil {
			tx.Rollback()
			context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
			return
		}
	} else {
		userPasswordData := models.UserPassword{
			UserPassword: hashPassword,
			UserNumber:   userNumber,
		}

		if err := tx.Create(&userPasswordData).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create password"})
			return
		}

	}

	updateUserStatus := tx.Model(&userData).
		Where("number = ?", userNumber).
		Update("status", "Active")
	if updateUserStatus.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	updateTokenStatus := tx.Model(&passwordReset).
		Where("number = ?", passwordReset.Number).
		Update("status", "Active")
	if updateTokenStatus.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create password"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Your account is now active, Please Login"})
}
