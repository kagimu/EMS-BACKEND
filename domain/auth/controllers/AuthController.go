package controllers

import (
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
	var userRoles []models.UserRole
	var userRolesList []responses.RoleView
	var permissions []string

	config.DB.Where("user_number = ?", userNumber).Find(&userRoles)

	for _, item := range userRoles {

		userRolesList = append(userRolesList, responses.RoleView{
			Number:      item.Role.Number,
			Name:        item.Role.Name,
			Description: item.Role.Description,
		})

		//get role permissions
		var rolePermissions []models.RolePermission
		config.DB.Where("role_number = ?", item.RoleNumber).Find(&rolePermissions)

		for _, permission := range rolePermissions {
			permissions = append(permissions, permission.Permission)
		}

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
	}

	//we generate jwt token
	//serialized := user.Serialize()
	/*date := time.Now().Add(time.Hour * 24 * 7)
	token, err := hassPassTokenGen.GenerateToken(user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}*/

	token, err := services.JWTAuthService().GenerateToken(user.Number)

	//context.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
	//context.JSON(http.StatusOK, gin.H{"message": "Login Successful", "data": user, "expires_in": date, "token": token})
	context.JSON(http.StatusOK, gin.H{"message": "Login Successful", "token": token, "data": userView})

}

func Profile(context *gin.Context) {

	userId := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userId = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found!"})
		return
	}

	var userProfile models.User

	if err := config.DB.
		Where("number = ?", userId).
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

	firstName := userPostData.FirstName
	secondName := userPostData.SecondName
	lastName := userPostData.LastName
	userGender := userPostData.Gender
	userEmail := userPostData.Email
	userPhone := userPostData.Phone

	// form validation for required fields
	if firstName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "First Name is required!"})
		return
	}

	if secondName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "First Name is required!"})
		return
	}

	if userGender == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Gender is required!"})
		return
	}

	if userEmail == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
		return
	}

	if userPhone == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Phone is required!"})
		return
	}

	validMale := userGender != "Male"
	validFemale := userGender != "Female"

	if validMale && validFemale {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Gender!"})
		return
	}

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

	//first time registration user type must be admin
	var user models.User
	userType := "User"
	userStatus := "Inactive"

	if userTypeChecker := config.DB.
		First(&user).Error; userTypeChecker != nil {
		userType = "Admin"
	}

	//creates a secure uuid user number
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

	config.DB.Save(&saveUser)

	//after you need to save user email
	saveEmailContact := models.Contact{
		Contact:     userEmail,
		ContactType: "Email",
		UserNumber:  saveUser.Number,
	}
	config.DB.Save(&saveEmailContact)

	//after you need to save user phone
	savePhoneContact := models.Contact{
		Contact:     userPhone,
		ContactType: "Phone",
		UserNumber:  saveUser.Number,
	}
	config.DB.Save(&savePhoneContact)

	//create / generate a token for the user to reset their password
	//userToken := tokenGenerator()
	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
		Token:      userToken,
		Status:     userStatus,
	}

	config.DB.Save(&saveUserPasswordToken)

	//email sender
	Token := userToken
	d := emailSender.Info{Email: userEmail, Token: Token, MailType: "WholesomeEd: User Account Confirmation", UserName: secondName, CurrentYear: time.Now().Year(), Type: "registration"}
	d.SendEmail()

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

	context.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Registered Successfully!",
		"data":      UserDataInfo,
		"tokenInfo": saveUserPasswordToken.Token,
	})

}

func ResetPassword(context *gin.Context) {

	type UserEmail struct {
		Email string `json:"email"`
	}
	var userMail UserEmail

	err := context.BindJSON(&userMail)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	Email := userMail.Email
	if Email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
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
	if userStatus == statusActive {
		//then we proceed to update the user
		//var userData models.User
		updateStatus := config.DB.Model(&userData).
			Where("number = ?", userNumber).
			Update("status", status)
		if updateStatus.Error != nil {
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

		config.DB.Save(&saveUserPasswordToken)

	} else {

		var passwordResetInfo models.PasswordReset
		if err := config.DB.
			Where("user_number = ?", userNumber).
			Where("status = ?", status).
			First(&passwordResetInfo).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Token not found"})
			return
		}

		userToken = passwordResetInfo.Token

	}

	//email sender
	d := emailSender.Info{Email: Email, Token: userToken, MailType: "Wholesome: User Account Confirmation", UserName: contact.User.SecondName, CurrentYear: time.Now().Year(), Type: "reset-password"}
	d.SendEmail()

	context.JSON(http.StatusOK, gin.H{"message": "Please check your email to continue resetting your password", "token": userToken})

}

func UserPassword(context *gin.Context) {

	type UserPassword struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var userPass UserPassword
	err := context.BindJSON(&userPass)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token := userPass.Token
	userPassword := userPass.Password
	confirmPassword := userPass.ConfirmPassword

	//	validate if token is provided
	if token == "" {
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
	}

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

	var userData models.User

	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	//if user already exist the user password should be updated
	var userPasswordInfo models.UserPassword

	getUserPassword := config.DB.
		Where("user_number = ?", userNumber).
		First(&userPasswordInfo)
	if getUserPassword.Error == nil {
		//this will just update user password data
		updateStatus := config.DB.Model(&userPasswordInfo).
			Where("user_number = ?", userNumber).
			Update("user_password", hashPassword)
		if updateStatus.Error != nil {
			context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
			return
		}
	} else {
		//do save user passwords if its a new password
		userPasswordData := models.UserPassword{
			UserPassword: hashPassword,
			UserNumber:   userNumber,
		}
		config.DB.Save(&userPasswordData)

	}
	/////////////////////////////////////////////////////////////////

	//then we proceed to update the user status
	updateUserStatus := config.DB.Model(&userData).
		Where("number = ?", userNumber).
		Update("status", "Active")
	if updateUserStatus.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	//////////////////////////////////////////////////////////////

	//then we update status of the token as well
	updateTokenStatus := config.DB.Model(&passwordReset).
		Where("number = ?", passwordReset.Number).
		Update("status", "Active")
	if updateTokenStatus.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Your account is now active, Please Login"})
}
