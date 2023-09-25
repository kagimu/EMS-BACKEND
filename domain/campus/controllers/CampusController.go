package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/campus/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/class/models"
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

	userNumber := context.Request.Header["User-Number"][0]

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

	var duplicateCampus models2.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	number := uuid.NewV4().String()
	saveCampus := models2.Campus{
		Number:       number,
		UserNumber:   userNumber,
		Name:         name,
		Phone:        phone,
		Email:        email,
		Location:     location,
		SchoolNumber: schoolNumber,
	}

	config.DB.Save(&saveCampus)

	//save Level
	for _, item := range level {
		levelNumber := uuid.NewV4().String()
		saveLevel := models2.Level{
			Number:       levelNumber,
			Name:         item,
			CampusNumber: saveCampus.Number,
		}
		config.DB.Save(&saveLevel)
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
}

func EditCampus(context *gin.Context) {
	type CampusInfo struct {
		Name         string   `json:"name"`
		Phone        string   `json:"phone"`
		Email        string   `json:"email"`
		Location     string   `json:"location"`
		Level        []string `json:"level"`
		SchoolNumber string   `json:"school_number"`
		CampusNumber string   `json:"campus_number"`
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
	campusNumber := campusInfo.CampusNumber

	userNumber := context.Request.Header["User-Number"][0]

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

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required!"})
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

	var duplicateCampus models2.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		Where("number != ?", campusNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	var campusDataInfo models2.Campus
	updateLevel := config.DB.Model(&campusDataInfo).
		Where("number = ?", campusNumber).
		Updates(models2.Campus{Name: name, Phone: phone, UserNumber: userNumber, Email: email, Location: location})
	if updateLevel.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	//save Level
	/*	for _, item := range level {
		levelNumber := uuid.NewV4().String()
		saveLevel := models2.Level{
			Number:       levelNumber,
			Name:         item,
			CampusNumber: saveCampus.Number,
		}
		config.DB.Save(&saveLevel)
	}*/

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Edited Successfully!", "data": campusDataInfo})
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
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campuses})

}

func CampusInfo(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var campuses models2.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Preload("Level").
		Order("created_at desc").
		First(&campuses).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campuses})

}

func FetchLevel(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus Number is required!"})
		return
	}

	var campus models2.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Preload("Level").
		Order("created_at desc").
		First(&campus).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Campus found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campus})

}

func DeleteCampus(context *gin.Context) {

	type CampusInfo struct {
		CampusNumber string `json:"campus_number"`
	}

	var campusInfo CampusInfo
	err := context.BindJSON(&campusInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	campusNumber := campusInfo.CampusNumber

	userNumber := context.Request.Header["User-Number"][0]

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	db := config.DB

	var campusChecker models3.ClassLevel
	if err := db.
		Where("campus_number = ?", campusNumber).
		First(&campusChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Campus, it has Children"})
		return
	}

	var campusData models2.Campus

	campusDataInfo := config.DB.
		Where("number = ?", campusNumber).
		First(&campusData)
	if campusDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "School not found!"})
		return
	}

	//delete Campus
	var levelData models2.Level
	db.Where("campus_number = ?", campusNumber).Delete(&levelData)

	//delete Campus
	db.Where("number = ?", campusNumber).Delete(&campusData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus deleted successfully!"})

}
