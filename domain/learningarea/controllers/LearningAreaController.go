package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/learningarea/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewLearningArea(context *gin.Context) {

	type LearningAreaInfo struct {
		Name         string `json:"name"`
		CampusNumber string `json:"campus_number"`
	}

	var learningAreaInfoInfo LearningAreaInfo
	err := context.BindJSON(&learningAreaInfoInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := learningAreaInfoInfo.Name
	campusNumber := learningAreaInfoInfo.CampusNumber

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

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus is required!"})
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

	var duplicateLearningArea models2.LearningArea
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("campus_number = ?", campusNumber).
		First(&duplicateLearningArea).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveLearningArea := models2.LearningArea{
		Number:       number,
		Name:         name,
		UserNumber:   userNumber,
		CampusNumber: campusNumber,
	}
	config.DB.Save(&saveLearningArea)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning Area Added Successfully!", "data": saveLearningArea})
}

func FetchLearningArea(context *gin.Context) {

	campusNumber := context.Query("campus_number")

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus is required!"})
		return
	}

	var learningAreas []models2.LearningArea
	if err := config.DB.
		Where("campus_number = ?", campusNumber).
		Order("created_at desc").
		Find(&learningAreas).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Learning Area!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas})
}

func AddSubjectToLearningArea(context *gin.Context) {

	type SubjectLearningAreaInfo struct {
		LearningAreaNumber string   `json:"learning_area_number"`
		SubjectNumbers     []string `json:"subject_number"`
	}

	var subjectLearningAreaInfoInfo SubjectLearningAreaInfo
	err := context.BindJSON(&subjectLearningAreaInfoInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	learningAreaNumber := subjectLearningAreaInfoInfo.LearningAreaNumber
	subjectNumbers := subjectLearningAreaInfoInfo.SubjectNumbers

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
		return
	}

	if learningAreaNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning Area is required!"})
		return
	}

	if len(subjectNumbers) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject Area is required!"})
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

	for _, item := range subjectNumbers {

		var duplicateSubLearningArea models2.SubjectLearningArea
		if duplicate := config.DB.
			Where("learning_area_number = ?", learningAreaNumber).
			Where("subject_number = ?", item).
			First(&duplicateSubLearningArea).Error; duplicate != nil {

			//		add a subject to a learning area once no duplicates

			number := uuid.NewV4().String()
			saveLevel := models2.SubjectLearningArea{
				Number:             number,
				LearningAreaNumber: learningAreaNumber,
				SubjectNumber:      item,
				UserNumber:         userNumber,
			}
			config.DB.Save(&saveLevel)

		}

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Subjects added"})
}
