package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/school/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewSchool(context *gin.Context) {

	type SchoolInfo struct {
		Name string `json:"name"`
	}

	var schoolInfo SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := schoolInfo.Name

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
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

	var duplicateSchool models3.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveSchool := models3.School{
		Number:     number,
		Name:       name,
		UserNumber: userNumber,
	}

	config.DB.Save(&saveSchool)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School Added Successfully!", "data": saveSchool})
}

func EditSchool(context *gin.Context) {

	type SchoolInfo struct {
		Name         string `json:"name"`
		SchoolNumber string `json:"school_number"`
	}

	var schoolInfo SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := schoolInfo.Name
	schoolNumber := schoolInfo.SchoolNumber

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School number is required!"})
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

	var duplicateSchool models3.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		Where("number != ?", schoolNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	var schoolData models3.School
	updateSchoolName := config.DB.Model(&schoolData).
		Where("number = ?", schoolNumber).
		Update("name", name)
	if updateSchoolName.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School Edited Successfully!", "data": schoolData})
}

func FetchSchoolInfo(context *gin.Context) {

	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var schools models3.School
	if err := config.DB.
		Where("number = ?", schoolNumber).
		Find(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools})
}

func FetchSchools(context *gin.Context) {

	var schools []models3.School
	if err := config.DB.
		Order("created_at desc").
		Find(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools})
}
