package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/campus/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewCampus(context *gin.Context) {

	type CampusInfo struct {
		Name         string   `json:"name"`
		Phone        string   `json:"phone"`
		Email        string   `json:"email"`
		Location     string   `json:"location"`
		Level        []string `json:"level"`
		SchoolNumber string   `json:"school_number"`
	}

	var campusInfo CampusInfo
	err := context.BindJSON(&campusInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := campusInfo.Name
	phone := campusInfo.Phone
	email := campusInfo.Email
	location := campusInfo.Location
	level := campusInfo.Level
	schoolNumber := campusInfo.SchoolNumber

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

	if phone == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Phone is required!"})
		return
	}

	if email == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required!"})
		return
	}

	if location == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Location is required!"})
		return
	}

	if len(level) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
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

	number := uuid.NewV4().String()
	///save user data
	saveCampus := models2.Campus{
		Number:     number,
		UserNumber: userNumber,
		Name:       name,
		Phone:      phone,
		Email:      email,
		Location:   location,
		//Level:        level,
		SchoolNumber: schoolNumber,
	}

	config.DB.Save(&saveCampus)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
}

func FetchCampuses(context *gin.Context) {

	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var campuses []models2.Campus
	if err := config.DB.
		Where("school_number = ?", schoolNumber).
		Order("created_at desc").
		Find(&campuses).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No class found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campuses})

}
