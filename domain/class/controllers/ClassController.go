package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/class/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewClass(context *gin.Context) {

	type ClassInfo struct {
		Name        string `json:"name"`
		LevelNumber string `json:"level_number"`
	}

	var classInfo ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := classInfo.Name
	levelNumber := classInfo.LevelNumber

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

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
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

	var duplicateClass models2.Class
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("level_number = ?", levelNumber).
		First(&duplicateClass).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveClass := models2.Class{
		Number:      number,
		Name:        name,
		UserNumber:  userNumber,
		LevelNumber: levelNumber,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func FetchClasses(context *gin.Context) {

	levelNumber := context.Query("level_number")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var classes []models2.Class
	if err := config.DB.
		Where("level_number = ?", levelNumber).
		Order("created_at desc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Campus!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes})
}

func FetchClassSubjects(context *gin.Context) {

	levelNumber := context.Query("level_number")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var classes []models2.Class
	var classes_ []models2.ClassInfo
	if err := config.DB.
		Where("level_number = ?", levelNumber).
		Order("created_at desc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Campus!"})
		return
	}

	for _, item := range classes {

		//fetch subjects under a learning area

		var subjClass []models2.SubjectCLass
		var subjects_ []models2.SubjectsInfo

		if err := config.DB.
			Where("class_number = ?", item.Number).
			Where("level_number = ?", levelNumber).
			Preload("Subject").
			Find(&subjClass).Error; err != nil {
			//context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Subj found"})
			//return
		}

		for _, lst := range subjClass {
			subjects_ = append(subjects_, models2.SubjectsInfo{
				Number: lst.Subject.Number,
				Name:   lst.Subject.Name,
			})
		}

		classes_ = append(classes_, models2.ClassInfo{
			Number:   item.Number,
			Name:     item.Name,
			Subjects: subjects_,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes_})
}

func NewClassSubjects(context *gin.Context) {

	type CLassSubjectInfo struct {
		ClassNumber    string   `json:"class_number"`
		CampusNumber   string   `json:"campus_number"`
		LevelNumber    string   `json:"level_number"`
		SubjectNumbers []string `json:"subject_number"`
	}

	var subjectClassInfo CLassSubjectInfo
	err := context.BindJSON(&subjectClassInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classNumber := subjectClassInfo.ClassNumber
	campusNumber := subjectClassInfo.CampusNumber
	levelNumber := subjectClassInfo.LevelNumber
	subjectNumbers := subjectClassInfo.SubjectNumbers

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning Area is required!"})
		return
	}

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
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

		var duplicateSubjectClass models2.SubjectCLass
		if duplicate := config.DB.
			Where("campus_number = ?", campusNumber).
			Where("level_number = ?", levelNumber).
			Where("subject_number = ?", item).
			First(&duplicateSubjectClass).Error; duplicate != nil {

			//		add a subject to a learning area once no duplicates
			number := uuid.NewV4().String()
			saveLevel := models2.SubjectCLass{
				Number:        number,
				ClassNumber:   classNumber,
				CampusNumber:  campusNumber,
				LevelNumber:   levelNumber,
				SubjectNumber: item,
				UserNumber:    userNumber,
			}
			config.DB.Save(&saveLevel)

		}

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Subjects added"})

}
