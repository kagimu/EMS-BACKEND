package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	models2 "github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/commons"
	"github.com/hamsoftug/wholesomeEd/domain/commons/emailSender"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Register(context *gin.Context) {

	//	we need to first initialize a struct that caters for posted data
	//then after this data is separated

	var userPostData models2.UserDataRequest
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
	var contact models2.Contact
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
	var user models2.User
	userType := "User"
	userStatus := "Inactive"

	if userTypeChecker := config.DB.
		First(&user).Error; userTypeChecker != nil {
		userType = "Admin"
	}

	//creates a secure uuid user number
	userNumber := uuid.NewV4().String()

	///save user data
	saveUser := models2.User{
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
	saveEmailContact := models2.Contact{
		Contact:     userEmail,
		ContactType: "Email",
		UserNumber:  saveUser.Number,
	}
	config.DB.Save(&saveEmailContact)

	//after you need to save user phone
	savePhoneContact := models2.Contact{
		Contact:     userPhone,
		ContactType: "Phone",
		UserNumber:  saveUser.Number,
	}
	config.DB.Save(&savePhoneContact)

	//create / generate a token for the user to reset their password
	//userToken := tokenGenerator()
	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models2.PasswordReset{
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
		models2.User
		Contacts []interface{}
	}

	UserDataInfo := UserInfo{
		saveUser,
		collection,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Registered Successfully!", "data": UserDataInfo})

}
