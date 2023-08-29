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

	number := uuid.NewV4().String()
	///save user data
	saveCampus := models3.School{
		Number:     number,
		Name:       name,
		UserNumber: userNumber,
	}

	config.DB.Save(&saveCampus)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
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
