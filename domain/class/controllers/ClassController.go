package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func NewClass(context *gin.Context) {

	var classInfo requests.ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := classInfo.Name
	abbreviation := classInfo.Abbreviation
	level := classInfo.Level

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if level == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	if abbreviation == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Abbreviation is required!"})
		return
	}

	var userData authModels.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateClass models.Class
	if duplicate := config.DB.
		Where("name = ? OR abbreviation = ?", name, abbreviation).
		First(&duplicateClass).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class cant be added exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveClass := models.Class{
		Number:       number,
		Name:         name,
		Abbreviation: abbreviation,
		UserNumber:   userNumber,
		Level:        level,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func AddClassToLevel(context *gin.Context) {

	type ClassInfo struct {
		ClassNumber  string `json:"class_number"`
		CampusNumber string `json:"campus_number"`
		LevelNumber  string `json:"level_number"`
	}

	var classInfo ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classNumber := classInfo.ClassNumber
	campusNumber := classInfo.CampusNumber
	levelNumber := classInfo.LevelNumber

	userNumber := context.Request.Header["User-Number"][0]

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus is required!"})
		return
	}

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var userData authModels.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateClass models.ClassLevel
	db := config.DB
	if duplicate := config.DB.
		Where("class_number = ?", classNumber).
		Where("campus_number = ?", campusNumber).
		Where("level_number = ?", levelNumber).
		First(&duplicateClass).Error; duplicate == nil {

		//check if this class has children such that it is not deleted

		var classChildrenChecker models.ClassLevelSubject
		if err := db.
			Where("class_level_number = ?", duplicateClass.Number).
			First(&classChildrenChecker).Error; err == nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Class, it has Children"})
			return
		}

		//------------------------------------------------------
		//delete class if it is removed from the level
		db.Where("class_number = ?", classNumber).
			Where("campus_number = ?", campusNumber).
			Where("level_number = ?", levelNumber).
			Delete(&duplicateClass)
		context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class removed successfully!"})

		context.JSON(http.StatusBadRequest, gin.H{"message": "Class already exists"})
		return
	}

	number := uuid.NewV4().String()
	saveClass := models.ClassLevel{
		Number:       number,
		ClassNumber:  classNumber,
		CampusNumber: campusNumber,
		LevelNumber:  levelNumber,
		UserNumber:   userNumber,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func FetchClasses(context *gin.Context) {

	level := context.Query("level")

	var classes []models.Class
	if err := config.DB.
		Scopes(levelFilter(level)).
		Order("created_at desc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes})
}

func FetchClassesToAdd(context *gin.Context) {

	level := context.Query("level")

	var levelsInfo campusModels.Level
	if err := config.DB.
		Where("number = ?", level).
		First(&levelsInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	var classes []models.Class
	if err := config.DB.
		Scopes(levelFilter(levelsInfo.Name)).
		Order("created_at desc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes})
}

func levelFilter(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("level like ?", "%"+level+"%")
		}
		return db
	}
}

func FetchClassLevel(context *gin.Context) {

	levelNumber := context.Query("level_number")
	campusNumber := context.Query("campus_number")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level Number is required!"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "campus Number is required!"})
		return
	}

	var classLevel []models.ClassLevel
	var classLevels []responses.ClassLevelInfo

	if err := config.DB.
		Where("level_number = ?", levelNumber).
		Where("campus_number = ?", campusNumber).
		Preload("Streams").
		Preload("Class").
		Order("created_at desc").
		Find(&classLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range classLevel {
		var stream_ []responses.StreamInfo
		/*var streams []models.Stream
		config.DB.Where("class_level_number = ?", lst.Number).Find(&streams)*/

		for _, list := range lst.Streams {
			stream_ = append(stream_, responses.StreamInfo{
				StreamNumber: list.Number,
				Name:         list.Name,
			})
		}

		classLevels = append(classLevels, responses.ClassLevelInfo{
			Number:       lst.Number,
			Name:         lst.Class.Name,
			ClassNumber:  lst.Class.Number,
			Abbreviation: lst.Class.Abbreviation,
			Streams:      stream_,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classLevels})
}

func FetchClassSubjects(context *gin.Context) {

	levelNumber := context.Query("level_number")
	classNumber := context.Query("class_number")

	/*	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}*/

	var classSubjects []models.ClassLevelSubject
	var subjects_ []models.SubjectInfo
	if err := config.DB.
		Scopes(classFilter(classNumber)).
		Where("class_levels.level_number = ?", levelNumber).
		Joins("join class_levels on class_levels.number = class_level_subjects.class_level_number").
		Preload("ClassLevel.Class").
		Preload("ClassLevel").
		Preload("Subject").
		Order("created_at desc").
		Find(&classSubjects).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject!"})
		return
	}

	for _, item := range classSubjects {
		subjects_ = append(subjects_, models.SubjectInfo{
			SubjectNumber: item.Subject.Number,
			Name:          item.Subject.Name,
			Class:         item.ClassLevel.Class.Name,
			Abbreviation:  item.ClassLevel.Class.Abbreviation,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjects_})
}

func classFilter(classNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if classNumber != "" {
			return db.Where("class_level_number = ?", classNumber)
		}
		return db
	}
}

func NewClassSubjects(context *gin.Context) {

	type CLassSubjectInfo struct {
		ClassNumber    string `json:"class_number"`
		SubjectNumbers string `json:"subject_number"`
	}

	var subjectClassInfo CLassSubjectInfo
	err := context.BindJSON(&subjectClassInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classNumber := subjectClassInfo.ClassNumber
	subjectNumbers := subjectClassInfo.SubjectNumbers

	userNumber := context.Request.Header["User-Number"][0]

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumbers == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var userData authModels.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateSubjectClass models.ClassLevelSubject
	if duplicate := config.DB.
		Where("class_level_number = ?", classNumber).
		Where("subject_number = ?", subjectNumbers).
		First(&duplicateSubjectClass).Error; duplicate != nil {

		//		add a subject to a learning area once no duplicates
		number := uuid.NewV4().String()
		saveLevel := models.ClassLevelSubject{
			Number:           number,
			ClassLevelNumber: classNumber,
			SubjectNumber:    subjectNumbers,
			UserNumber:       userNumber,
		}
		config.DB.Save(&saveLevel)

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": "Subjects added"})

}

func AddClassStream(ctx *gin.Context) {

	var streamRequest requests.StreamRequest

	err := ctx.BindJSON(&streamRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to serialize  request error!"})
		return
	}

	streams := strings.Split(streamRequest.Name, ",")

	//add streams to the class
	if streams != nil {
		for _, stream := range streams {
			streamObj := models.Stream{
				Number:           uuid.NewV4().String(),
				Name:             stream,
				ClassLevelNumber: streamRequest.ClassLevelNumber,
				UserNumber:       ctx.GetHeader("User-Number"),
			}

			err = config.DB.Save(&streamObj).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusOK, "message": "Failed to add stream !"})
				return
			}

		}
	}
	/*} else {
		streamObj := models.Stream{
			Number:      uuid.NewV4().String(),
			Name:        saveClass.Name,
			ClassNumber: saveClass.Number,
			UserNumber:  userNumber,
		}

		tx.Save(&streamObj)
	}*/

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Stream Added Successfully!"})

}

func FetchStreamsByClassLevelNumber(ctx *gin.Context) {
	var streams []models.Stream
	classLevelNumber := ctx.Query("class_level_number")

	if classLevelNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Class Level Number is required!"})
		return
	}

	if err := config.DB.
		Where("class_level_number = ?", classLevelNumber).
		Find(&streams).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Streams!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": streams})

}
