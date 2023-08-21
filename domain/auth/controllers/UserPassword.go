package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/hamsoftug/wholesomeEd/domain/commons/hassPassTokenGen"
	"net/http"
)

func UserPassword(context *gin.Context) {

	type UserPassword struct {
		Token           string `json:"token"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var userPass UserPassword
	err := context.BindJSON(&userPass)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
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
	////////////////////////////////////////////////////////////////

	//then we proceed to update the user status
	updateUserStatus := config.DB.Model(&userData).
		Where("number = ?", userNumber).
		Update("status", "Active")
	if updateUserStatus.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	/////////////////////////////////////////////////////////////

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
