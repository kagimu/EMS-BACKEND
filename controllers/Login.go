package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/config"
	"github.com/hamsoftug/wholesomeEd/controllers/profile"
	"github.com/hamsoftug/wholesomeEd/hassPassTokenGen"
	"github.com/hamsoftug/wholesomeEd/models"
	"net/http"
	"regexp"
	"strings"
	"time"
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
		context.JSON(http.StatusNotFound, gin.H{"message": "Account not Active, please reset password to activate your account"})
		return
	}

	userHashPassword := userPassword.UserPassword

	//	we need to validate hashed password
	if !hassPassTokenGen.CheckHash(Password, userHashPassword) {
		context.JSON(http.StatusNotFound, gin.H{"message": "Invalid Password!"})
		return
	}

	//update last seen for user
	updateLastSeen := config.DB.Model(&user).
		Where("number = ?", userNumber).
		Update("seen", time.Now())
	if updateLastSeen.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
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
			//UserNumber:  item.UserNumber,
		})
	}

	_userProfile := models.GetUserList{
		ID:          userProfile.ID,
		Number:      userProfile.Number,
		Name:        userProfile.Name,
		DateOfBirth: userProfile.DateOfBirth,
		Gender:      userProfile.Gender,
		Type:        userProfile.Type,
		//CoverPic:    profile.CoverProfilePic(context, "cover_photo", userNumber, "defaultCoverPic.png"),
		ProfilePic:  profile.CoverProfilePic(context, "profile", userNumber, "default.png"),
		ContactList: contactListInfo,
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
