package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/responses"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func NewClass(context *gin.Context) {

	var classInfo requests.ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := classInfo.Code
	name := classInfo.Name
	abbreviation := strings.ToUpper(classInfo.Abbreviation)
	level := classInfo.Level

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
		return
	}

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

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
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

	var duplicateCode models.Class
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists!"})
		return
	}

	var duplicateClass models.Class
	if duplicate := config.DB.
		Where("name = ?", name).
		First(&duplicateClass).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class already exists!"})
		return
	}

	var duplicateAbbreviation models.Class
	if duplicate := config.DB.
		Where("abbreviation = ?", abbreviation).
		First(&duplicateAbbreviation).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Abbreviation already exists!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save user data
	saveClass := models.Class{
		Number:       uuid.NewV4().String(),
		Name:         name,
		Abbreviation: abbreviation,
		UserNumber:   userNumber,
		LevelNumber:  level,
		Code:         codeNo,
	}

	if err := tx.Create(&saveClass).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create class"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add class"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func AddClassToLevel(context *gin.Context) {

	type ClassInfo struct {
		ClassNumber  string `json:"class_number"`
		CampusNumber string `json:"campus_number"`
		LevelNumber  string `json:"level_number"`
		Code         int64  `json:"class_code"`
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
	code := classInfo.Code

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
		return
	}

	number := uuid.NewV4().String()
	saveClass := models.ClassLevel{
		Number:       number,
		ClassNumber:  classNumber,
		CampusNumber: campusNumber,
		LevelNumber:  levelNumber,
		UserNumber:   userNumber,
		Code:         code,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func AddCurriculumClass(context *gin.Context) {

	type CurriculumClassInfo struct {
		LevelNumber string `json:"level_number"`
		ClassNumber string `json:"class_number"`
		Code        int64  `json:"code"`
	}

	var curriculumClassInfo CurriculumClassInfo
	err := context.BindJSON(&curriculumClassInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classNumber := curriculumClassInfo.ClassNumber
	levelNumber := curriculumClassInfo.LevelNumber
	code := curriculumClassInfo.Code

	userNumber := context.Request.Header["User-Number"][0]

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
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

	var duplicateClass models.CurriculumClass
	db := config.DB
	if duplicate := config.DB.
		Where("class_number = ?", classNumber).
		Where("curriculum_level_number = ?", levelNumber).
		First(&duplicateClass).Error; duplicate == nil {

		//check if this class has children such that it is not deleted

		/*var classChildrenChecker models.ClassLevelSubject
		if err := db.
			Where("class_level_number = ?", duplicateClass.Number).
			First(&classChildrenChecker).Error; err == nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Class, it has Children"})
			return
		}*/

		//------------------------------------------------------
		//delete class if it is removed from the level
		db.Unscoped().Where("class_number = ?", classNumber).
			Where("curriculum_level_number = ?", levelNumber).
			Delete(&duplicateClass)
		context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class removed successfully!"})
		return
	}

	saveClass := models.CurriculumClass{
		Number:                uuid.NewV4().String(),
		ClassNumber:           classNumber,
		CurriculumLevelNumber: levelNumber,
		UserNumber:            userNumber,
		Code:                  code,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveClass})
}

func FetchClasses(context *gin.Context) {

	level := context.Query("level")

	var levels []levelModels.Level
	var levels_ []responses2.ClassLevelList
	if err := config.DB.
		Scopes(levelDataFilter(level)).
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {

		var classes []models.Class
		var classes_ []responses.ClassDataList
		if err := config.DB.
			Where("level_number = ?", list.Number).
			Order("code asc").
			Find(&classes).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
			return
		}

		if len(classes) > 0 {
			for _, lst := range classes {
				classes_ = append(classes_, responses.ClassDataList{
					Number:       lst.Number,
					Name:         lst.Name,
					Code:         lst.Code,
					Abbreviation: lst.Abbreviation,
					UserNumber:   lst.UserNumber,
					LevelNumber:  lst.LevelNumber,
				})
			}

			levels_ = append(levels_, responses2.ClassLevelList{
				Number:        list.Number,
				Name:          list.Name,
				Code:          list.Code,
				UserNumber:    list.UserNumber,
				ClassDataList: classes_,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levels_})
}

func FetchCurriculumClasses(context *gin.Context) {

	level := context.Query("level")
	curriculumNumber := context.Query("curriculum_number")

	var levels []levelModels.CurriculumLevel
	var levels_ []responses2.ClassLevelList
	if err := config.DB.
		Scopes(levelInfoFilter(level)).
		Where("curriculum_Number = ?", curriculumNumber).
		Preload("Level").
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {

		var classes []models.CurriculumClass
		var classes_ []responses.ClassDataList
		if err := config.DB.
			Where("curriculum_level_number = ?", list.Number).
			Preload("Class").
			Order("code asc").
			Find(&classes).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
			return
		}

		if len(classes) > 0 {
			for _, lst := range classes {
				classes_ = append(classes_, responses.ClassDataList{
					Number:       lst.Number,
					Name:         lst.Class.Name,
					ClassNumber:  lst.Class.Number,
					Code:         lst.Code,
					Abbreviation: lst.Class.Abbreviation,
					UserNumber:   lst.UserNumber,
				})
			}

			levels_ = append(levels_, responses2.ClassLevelList{
				Number:        list.Number,
				Name:          list.Level.Name,
				Code:          list.Code,
				UserNumber:    list.UserNumber,
				ClassDataList: classes_,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levels_})
}

func FilterCurriculumClasses(context *gin.Context) {

	level := context.Query("level")

	var classes []models.Class
	var classes_ []responses.ClassDataList
	if err := config.DB.
		Where("level_number = ?", level).
		Order("code asc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range classes {
		classes_ = append(classes_, responses.ClassDataList{
			Number:       lst.Number,
			Name:         lst.Name,
			Code:         lst.Code,
			Abbreviation: lst.Abbreviation,
			UserNumber:   lst.UserNumber,
			LevelNumber:  lst.LevelNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes_})
}

func FetchSchemeOfWorkClasses(context *gin.Context) {

	level := context.Query("level")
	curriculumNumber := context.Query("curriculum_number")

	var classes []models.CurriculumClass
	var classes_ []responses.ClassDataList
	if err := config.DB.
		Scopes(levelDataFilter(level)).
		Where("curriculum_levels.curriculum_number = ?", curriculumNumber).
		Joins("join curriculum_levels on curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Joins("join classes on curriculum_classes.class_number = classes.number").
		Preload("CurriculumLevel").
		Preload("Class").
		Order("code asc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, lst := range classes {
		classes_ = append(classes_, responses.ClassDataList{
			Number:       lst.Number,
			Name:         lst.Class.Name,
			Code:         lst.Code,
			Abbreviation: lst.Class.Abbreviation,
			UserNumber:   lst.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes_})
}

func FetchClassesCurriculum(context *gin.Context) {
	curriculumNumber := context.Query("curriculum_number")
	var classes []models.CurriculumClass
	var classes_ []responses.ClassDataList
	if err := config.DB.
		Where("curriculum_levels.curriculum_number = ?", curriculumNumber).
		Joins("join curriculum_levels on curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Joins("join classes on curriculum_classes.class_number = classes.number").
		Preload("CurriculumLevel").
		Preload("Class").
		Order("code asc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range classes {
		classes_ = append(classes_, responses.ClassDataList{
			Number:       lst.Number,
			Name:         lst.Class.Name,
			ClassNumber:  lst.Class.Number,
			Code:         lst.Code,
			Abbreviation: lst.Class.Abbreviation,
			UserNumber:   lst.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes_})
}

func FetchClassesToAdd(context *gin.Context) {

	level := context.Query("level")

	var campusLevelData levelModels.CampusLevel
	if err := config.DB.
		Where("number = ?", level).
		First(&campusLevelData).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campusLevelData})

	/*var levelData levelModels.Level
	if err := config.DB.
		Where("number = ?", level).
		First(&levelData).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	var classes []models.Class
	if err := config.DB.
		//Scopes(levelFilter(levelData.LearningLevelNumber)).
		Order("code asc").
		Find(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes})*/
}

func levelDataFilter(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("curriculum_levels.number = ?", level)
		}
		return db
	}
}

func levelInfoFilter(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("level_number = ?", level)
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
		Order("code asc").
		Find(&classLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range classLevel {
		var stream_ []responses.StreamInfo

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

	if streamRequest.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if streamRequest.ClassLevelNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Class Level Number is required!"})
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
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add stream !"})
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Stream Added Successfully!"})

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
