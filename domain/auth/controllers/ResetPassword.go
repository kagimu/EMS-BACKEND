package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/commons/emailSender"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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
	status := "Inactive"
	var userData models.User

	//then we proceed to update the user

	//var userData models.User
	updateStatus := config.DB.Model(&userData).
		Where("number = ?", userNumber).
		Update("status", status)
	if updateStatus.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	userToken := uuid.NewV4().String()
	saveUserPasswordToken := models.PasswordReset{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
		Token:      userToken,
		Status:     status,
	}

	config.DB.Save(&saveUserPasswordToken)

	//email sender

	d := emailSender.Info{Email: Email, Token: userToken, MailType: "TutorView: User Account Confirmation", UserName: contact.User.SecondName, CurrentYear: time.Now().Year(), Type: "reset-password"}
	d.SendEmail()

	context.JSON(http.StatusOK, gin.H{"message": "Please check your email to continue resetting your password"})

}
