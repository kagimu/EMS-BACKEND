package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"html"
	"net/http"
	"strconv"
)

func NewLearningArea(context *gin.Context) {

	var learningAreaInfo requests.LearningAreaInfo
	err := context.BindJSON(&learningAreaInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	learningAreaInfo.Prepare()

	if err := learningAreaInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	name := learningAreaInfo.Name
	code := learningAreaInfo.Code
	curriculumNumber := learningAreaInfo.CurriculumNumber
	userNumber := context.Request.Header["User-Number"][0]

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
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

	var duplicateCode models2.LearningArea
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateLearningArea models2.LearningArea
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateLearningArea).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists!!"})
		return
	}

	///save user data
	tx := config.DB.Begin()
	saveLearningArea := models2.LearningArea{
		Number:           uuid.NewV4().String(),
		Name:             name,
		Code:             codeNo,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}

	if err := tx.Create(&saveLearningArea).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create Learning area"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create learning area"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning Area Added Successfully!", "data": saveLearningArea})
}

func FetchLearningArea(context *gin.Context) {
	curriculumNumber := context.Query("curriculum_number")
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum("", curriculumNumber)
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses2.LearningAreaInfoData

	for _, lst := range learningAreaResp {
		learningAreas_ = append(learningAreas_, responses2.LearningAreaInfoData{
			Number:     lst.Number,
			Name:       html.UnescapeString(lst.Name),
			Code:       lst.Code,
			UserNumber: lst.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})
}

func DeleteLearningArea(context *gin.Context) {

	var learningAreaInfo requests.LearningAreaDelRequest
	err := context.BindJSON(&learningAreaInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	learningAreaNumber := learningAreaInfo.Number

	var learningAreaChecker models3.SubjectCluster
	if err := config.DB.
		Where("learning_area_number = ?", learningAreaNumber).
		First(&learningAreaChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this learning area, it has subject cluster."})
		return
	}

	var learningAreaData models2.LearningArea
	academicYearDataInfo := config.DB.
		Where("number = ?", learningAreaNumber).
		First(&learningAreaData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Learning Area not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", learningAreaNumber).
		Delete(&learningAreaData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete learning area"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete learning area"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning area deleted successfully!"})
}

func LearningAreaInfo(context *gin.Context) {
	learningAreaNumber := context.Query("learning_area_number")

	if learningAreaNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning Area is required!"})
		return
	}

	var learningArea models2.LearningArea
	if err := config.DB.
		Where("number = ?", learningAreaNumber).
		First(&learningArea).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No learning area found"})
		return
	}

	learningArea_ := responses2.LearningAreaInfoData{
		Number:     learningArea.Number,
		Name:       html.UnescapeString(learningArea.Name),
		Code:       learningArea.Code,
		UserNumber: learningArea.UserNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningArea_})
}

func LearningAreaEdit(context *gin.Context) {

	var learningAreaInfo requests.LearningAreaInfo
	err := context.BindJSON(&learningAreaInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	learningAreaInfo.Prepare()

	if err := learningAreaInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	name := learningAreaInfo.Name
	code := learningAreaInfo.Code
	curriculumNumber := learningAreaInfo.CurriculumNumber
	learningAreaNumber := learningAreaInfo.LearningAreaNumber
	userNumber := context.Request.Header["User-Number"][0]

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
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

	var duplicateCode models2.LearningArea
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", learningAreaNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateLearningArea models2.LearningArea
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", learningAreaNumber).
		First(&duplicateLearningArea).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists!!"})
		return
	}

	///save user data
	tx := config.DB.Begin()

	var learningAreaDataInfo models2.LearningArea
	updateLevel := tx.Model(&learningAreaDataInfo).
		Where("number = ?", learningAreaNumber).
		Updates(models2.LearningArea{Name: name, UserNumber: userNumber, Code: codeNo})
	if updateLevel.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create learning area"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning Area Edited Successfully!", "data": learningAreaDataInfo})
}
