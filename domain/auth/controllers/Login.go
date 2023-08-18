package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	models2 "github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/commons/hassPassTokenGen"
	"net/http"
	"regexp"
	"strings"
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
	var contact models2.Contact

	userContact := config.DB.
		Where("contact = ?", Email).
		First(&contact)
	if userContact.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Email not found!"})
		return
	}

	userNumber := contact.UserNumber

	//get user password from the password table
	var userPassword models2.UserPassword
	getUserPassword := config.DB.
		Where("user_number = ?", userNumber).
		First(&userPassword)
	if getUserPassword.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Password!"})
		return
	}

	//check if user is active / not blocked
	var user models2.User
	config.DB.Where("number = ?", userNumber).First(&user)
	if user.Status != "Active" {
		context.JSON(http.StatusNotFound, gin.H{"message": "Account not Active, please reset password to activate your account"})
		return
	}

	userHashPassword := userPassword.UserPassword

	//	we need to validate hashed password
	if !hassPassTokenGen.CheckHash(Password, userHashPassword) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Password!"})
		return
	}

	//return User profile
	var userProfile models2.User
	if err := config.DB.
		Where("number = ?", userNumber).
		Preload("Contacts").
		First(&userProfile).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	var contactListInfo []models2.ContactList

	for _, item := range userProfile.Contacts {
		contactListInfo = append(contactListInfo, models2.ContactList{
			ID:          item.ID,
			Contact:     item.Contact,
			ContactType: item.ContactType,
		})
	}

	_userProfile := models2.GetUserList{
		ID:         userProfile.ID,
		Number:     userProfile.Number,
		FirstName:  userProfile.FirstName,
		SecondName: userProfile.SecondName,
		LastName:   userProfile.LastName,
		Gender:     userProfile.Gender,
		Type:       userProfile.Type,
	}

	//we generate jwt token
	//serialized := user.Serialize()
	/*date := time.Now().Add(time.Hour * 24 * 7)
	token, err := hassPassTokenGen.GenerateToken(user)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}*/

	//context.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
	//context.JSON(http.StatusOK, gin.H{"message": "Login Successful", "data": user, "expires_in": date, "token": token})
	context.JSON(http.StatusOK, gin.H{"message": "Login Successful", "data": _userProfile})

}
