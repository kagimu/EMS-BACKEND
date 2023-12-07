package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/responses"
	schoolCurriculumModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func NewSchool(context *gin.Context) {

	type SchoolInfo struct {
		Name       string   `json:"name"`
		Curriculum []string `json:"curriculum"`
	}

	var schoolInfo SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := schoolInfo.Name
	curriculum := schoolInfo.Curriculum

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if len(curriculum) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
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

	var duplicateSchool schoolCurriculumModels.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	tx := config.DB.Begin()
	///save user data
	saveSchool := schoolCurriculumModels.School{
		Number:     uuid.NewV4().String(),
		Name:       name,
		UserNumber: userNumber,
	}

	if err := tx.Create(&saveSchool).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school"})
		return
	}

	for _, item := range curriculum {
		saveSchoolCurriculum := schoolCurriculumModels.SchoolCurriculum{
			Number:           uuid.NewV4().String(),
			CurriculumNumber: item,
			SchoolNumber:     saveSchool.Number,
		}

		if err := tx.Create(&saveSchoolCurriculum).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			//context.JSON(http.StatusBadRequest, gin.H{"message": err})
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save curriculum"})
			return
		}

	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school!!"})
		return
	}

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

	userNumber := context.Request.Header["User-Number"][0]

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

	var duplicateSchool schoolCurriculumModels.School
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		Where("number != ?", schoolNumber).
		First(&duplicateSchool).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School already exists"})
		return
	}

	var schoolData schoolCurriculumModels.School
	updateSchoolName := config.DB.Model(&schoolData).
		Where("number = ?", schoolNumber).
		Update("name", name)
	if updateSchoolName.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School Edited Successfully!", "data": schoolData})
}

func DeleteSchool(context *gin.Context) {
	type SchoolInfo struct {
		SchoolNumber string `json:"school_number"`
	}

	var schoolInfo SchoolInfo
	err := context.BindJSON(&schoolInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	schoolNumber := schoolInfo.SchoolNumber

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

	var schoolChecker models2.Campus
	if err := db.
		Where("school_number = ?", schoolNumber).
		First(&schoolChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this School, it has Children"})
		return
	}

	var schoolData schoolCurriculumModels.School

	classLevelInfo := config.DB.
		Where("number = ?", schoolNumber).
		First(&schoolData)
	if classLevelInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "School not found!"})
		return
	}

	db.Where("number = ?", schoolNumber).Delete(&schoolData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "School deleted successfully!"})

}

func FetchSchoolInfo(context *gin.Context) {

	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "School Number is required!"})
		return
	}

	var schools schoolCurriculumModels.School
	var curriculum_ []responses.SchoolCurriculumResponseInfo
	if err := config.DB.
		Where("number = ?", schoolNumber).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		First(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	for _, schoolCurr := range schools.SchoolCurriculum {
		curriculum_ = append(curriculum_, responses.SchoolCurriculumResponseInfo{
			Name:   schoolCurr.Curriculum.Name,
			Number: schoolCurr.Number,
		})
	}

	schools_ := responses.SchoolResponseInfo{
		Number:           schools.Number,
		Name:             schools.Name,
		SchoolCurriculum: curriculum_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools_})
}

func FetchSchools(context *gin.Context) {

	var schools []schoolCurriculumModels.School
	if err := config.DB.
		Order("created_at desc").
		Find(&schools).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Schools!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schools})
}
