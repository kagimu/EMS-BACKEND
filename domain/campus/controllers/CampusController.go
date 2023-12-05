package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewCampus(context *gin.Context) {

	type CampusInfo struct {
		Name                   string   `json:"name"`
		Phone                  string   `json:"phone"`
		Email                  string   `json:"email"`
		Location               string   `json:"location"`
		SchoolNumber           string   `json:"school_number"`
		SchoolCurriculumNumber string   `json:"school_curriculum_number"`
		CurriculumLevelNumber  []string `json:"curriculum_level_number"`
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
	curriculumLevelNumber := campusInfo.CurriculumLevelNumber
	schoolNumber := campusInfo.SchoolNumber
	schoolCurriculumNumber := campusInfo.SchoolCurriculumNumber

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

	if len(curriculumLevelNumber) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School number is required!"})
		return
	}

	if schoolCurriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Curriculum is required!"})
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

	var duplicateCampus campusModel.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	number := uuid.NewV4().String()
	saveCampus := campusModel.Campus{
		Number:                 number,
		UserNumber:             userNumber,
		Name:                   name,
		Phone:                  phone,
		Email:                  email,
		Location:               location,
		SchoolNumber:           schoolNumber,
		SchoolCurriculumNumber: schoolCurriculumNumber,
	}

	if err := tx.Create(&saveCampus).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create campus"})
		return
	}

	//save Level
	for _, item := range curriculumLevelNumber {

		var curriculumLevel schoolLevelModels.CurriculumLevel
		if err := tx.
			Where("number = ?", item).
			First(&curriculumLevel).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found!"})
			return
		}

		saveLevel := schoolLevelModels.SchoolLevel{
			Number:                uuid.NewV4().String(),
			Code:                  curriculumLevel.Code,
			CurriculumLevelNumber: item,
			CampusNumber:          saveCampus.Number,
		}

		if err := tx.Create(&saveLevel).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			//context.JSON(http.StatusBadRequest, gin.H{"message": err})
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create level"})
			return
		}

	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add campus!!"})
		return
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

	var duplicateCampus campusModel.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		Where("number != ?", campusNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}

	var campusDataInfo campusModel.Campus
	updateLevel := config.DB.Model(&campusDataInfo).
		Where("number = ?", campusNumber).
		Updates(campusModel.Campus{Name: name, Phone: phone, UserNumber: userNumber, Email: email, Location: location})
	if updateLevel.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	//save Level
	/*	for _, item := range level {
		levelNumber := uuid.NewV4().String()
		saveLevel := campusModel.Level{
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

	var campuses []campusModel.Campus
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

	var campuses campusModel.Campus
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

	var campus campusModel.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Preload("Level.LearningLevel").
		Preload("Level").
		Order("created_at desc").
		First(&campus).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Campus found"})
		return
	}

	/*	var level_ []responses.CampusLevelList
		for _, lvl := range campus.Level {
			level_ = append(level_, responses.CampusLevelList{
				Number:              lvl.Number,
				LearningLevelNumber: lvl.LearningLevelNumber,
				LevelName:           lvl.LearningLevel.Name,
				CampusNumber:        lvl.CampusNumber,
			})
		}*/

	campus_ := responses.CampusDataList{
		Number:       campus.Number,
		Name:         campus.Name,
		Phone:        campus.Phone,
		Email:        campus.Email,
		Location:     campus.Location,
		UserNumber:   campus.UserNumber,
		SchoolNumber: campus.SchoolNumber,
		//Level:        level_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campus_})

}

func FetchCampusLevel(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus Number is required!"})
		return
	}

	var levelInfo []schoolLevelModels.Level
	var levelInfo_ []schoolLevelModels.LevelDataInfo
	if err := config.DB.
		//Select("name", "number").
		Where("campus_number = ?", campusNumber).
		Preload("LearningLevel").
		Order("created_at desc").
		Find(&levelInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level found"})
		return
	}

	// Populate the selected columns from the fetched data
	/*for _, level := range levelInfo {
		levelInfo_ = append(levelInfo_, campusModel.LevelDataInfo{
			Name:   level.LearningLevel.Name,
			Code:   level.LearningLevel.Code,
			Number: level.Number,
		})
	}*/

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levelInfo_})

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

	var campusData campusModel.Campus

	campusDataInfo := config.DB.
		Where("number = ?", campusNumber).
		First(&campusData)
	if campusDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "School not found!"})
		return
	}

	//delete Campus
	var levelData schoolLevelModels.Level
	db.Where("campus_number = ?", campusNumber).Delete(&levelData)

	//delete Campus
	db.Where("number = ?", campusNumber).Delete(&campusData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus deleted successfully!"})

}
