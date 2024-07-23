package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"net/http"
	"strconv"
	"strings"

	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/responses"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func NewClass(context *gin.Context) {

	var classInfo requests.ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classInfo.Prepare()

	if err := classInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := classInfo.Code
	name := classInfo.Name
	abbreviation := strings.ToUpper(classInfo.Abbreviation)
	level := classInfo.Level

	userNumber := context.Request.Header["User-Number"][0]

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

	var duplicateCode classModel.Class
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists!"})
		return
	}

	var duplicateClass classModel.Class
	if duplicate := config.DB.
		Where("name = ?", name).
		First(&duplicateClass).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class already exists!"})
		return
	}

	var duplicateAbbreviation classModel.Class
	if duplicate := config.DB.
		Where("abbreviation = ?", abbreviation).
		First(&duplicateAbbreviation).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Abbreviation already exists!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save user data
	saveClass := classModel.Class{
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

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class Added Successfully!", "data": saveClass})
}

func EditClass(context *gin.Context) {

	var classInfo requests.ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classInfo.Prepare()

	if err := classInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := classInfo.Code
	name := classInfo.Name
	abbreviation := strings.ToUpper(classInfo.Abbreviation)
	level := classInfo.Level
	classNumber := classInfo.ClassNumber

	userNumber := context.Request.Header["User-Number"][0]

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
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

	var duplicateCode classModel.Class
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("number != ?", classNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists!"})
		return
	}

	var duplicateClass classModel.Class
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("number != ?", classNumber).
		First(&duplicateClass).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class already exists!"})
		return
	}

	var duplicateAbbreviation classModel.Class
	if duplicate := config.DB.
		Where("abbreviation = ?", abbreviation).
		Where("number != ?", classNumber).
		First(&duplicateAbbreviation).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Abbreviation already exists!"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	var classDataInfo classModel.Class
	updateLevel := tx.Model(&classDataInfo).
		Where("number = ?", classNumber).
		Updates(classModel.Class{Name: name, UserNumber: userNumber, Abbreviation: abbreviation, Code: codeNo, LevelNumber: level})
	if updateLevel.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add class"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class Edited Successfully!", "data": classDataInfo})
}

func AddCampusClass(context *gin.Context) {

	type ClassInfo struct {
		Code              int64  `json:"class_code"`
		ClassNumber       string `json:"class_number"`
		CampusNumber      string `json:"campus_number"`
		CampusLevelNumber string `json:"campus_level_number"`
	}

	var classInfo ClassInfo
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := classInfo.Code
	classNumber := classInfo.ClassNumber
	campusNumber := classInfo.CampusNumber
	campusLevelNumber := classInfo.CampusLevelNumber

	userNumber := context.Request.Header["User-Number"][0]

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus is required!"})
		return
	}

	if campusLevelNumber == "" {
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

	var duplicateClass classModel.CampusClass
	db := config.DB
	if duplicate := config.DB.
		Where("curriculum_class_number = ?", classNumber).
		Where("campus_number = ?", campusNumber).
		Where("campus_level_number = ?", campusLevelNumber).
		First(&duplicateClass).Error; duplicate == nil {

		//check if this class has children such that it is not deleted

		/*var classChildrenChecker classModel.ClassLevelSubject
		if err := db.
			Where("class_level_number = ?", duplicateClass.Number).
			First(&classChildrenChecker).Error; err == nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Class, it has Children"})
			return
		}*/

		//------------------------------------------------------
		//delete class if it is removed from the level
		db.Unscoped().Where("curriculum_class_number = ?", classNumber).
			Where("campus_number = ?", campusNumber).
			Where("campus_level_number = ?", campusLevelNumber).
			Delete(&duplicateClass)
		context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class removed successfully!"})
		return
	}

	number := uuid.NewV4().String()
	saveClass := classModel.CampusClass{
		Number:                number,
		CurriculumClassNumber: classNumber,
		CampusNumber:          campusNumber,
		CampusLevelNumber:     campusLevelNumber,
		UserNumber:            userNumber,
		Code:                  code,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class Added Successfully!", "data": saveClass})
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

	var duplicateClass classModel.CurriculumClass
	db := config.DB
	if duplicate := config.DB.
		Where("class_number = ?", classNumber).
		Where("curriculum_level_number = ?", levelNumber).
		First(&duplicateClass).Error; duplicate == nil {

		//check if this class has children such that it is not deleted

		/*var classChildrenChecker classModel.ClassLevelSubject
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

	saveClass := classModel.CurriculumClass{
		Number:                uuid.NewV4().String(),
		ClassNumber:           classNumber,
		CurriculumLevelNumber: levelNumber,
		UserNumber:            userNumber,
		Code:                  code,
	}
	config.DB.Save(&saveClass)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Class Added Successfully!", "data": saveClass})
}

func FetchClasses(context *gin.Context) {

	level := context.Query("level")

	var levels []levelModels.Level
	var levels_ []responses2.ClassLevelList
	if err := config.DB.
		Scopes(filters.ByPrimaryKey(level)).
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {

		var classes []classModel.Class
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
		Scopes(filters.ByLevelNumber(level)).
		Where("curriculum_Number = ?", curriculumNumber).
		Preload("Level").
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {

		var classes []classModel.CurriculumClass
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

	var classes []classModel.Class
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
	campusNumber := context.Query("campus_number")

	//fetch curriculum under campus
	var campus campusModel.Campus
	if err := config.DB.
		Where("number = ?", campusNumber).
		Preload("SchoolCurriculum.Curriculum").
		Preload("SchoolCurriculum").
		Order("created_at desc").
		First(&campus).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No campus found"})
		return
	}

	levelNo := ""

	if level != "" {
		//fetch levels under campus
		var campusLevelInfo levelModels.CampusLevel
		if err := config.DB.
			Where("number = ?", level).
			First(&campusLevelInfo).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level found"})
			return
		}
		levelNo = campusLevelInfo.CurriculumLevelNumber
	}

	//fetch curriculum classes
	var classes []classModel.CurriculumClass
	var classes_ []responses.ClassDataList
	if err := config.DB.
		Scopes(filters.ByCurriculumLevelNumber(levelNo)).
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
	var classes []classModel.CurriculumClass
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

func FetchCampusClassesByLevel(context *gin.Context) {
	campusLevelNumber := context.Query("campus_level_number")

	var classLevel []classModel.CampusClass
	var classLevel_ []responses.CampusClassInfo
	if err := config.DB.
		Where("campus_level_number = ?", campusLevelNumber).
		Preload("CurriculumClass.Class").
		Preload("CurriculumClass").
		Preload("Streams").
		Order("code asc").
		Find(&classLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range classLevel {

		classLevel_ = append(classLevel_, responses.CampusClassInfo{
			Number:                lst.Number,
			Name:                  lst.CurriculumClass.Class.Name,
			ClassNumber:           lst.CurriculumClass.Class.Number,
			Abbreviation:          lst.CurriculumClass.Class.Abbreviation,
			CurriculumClassNumber: lst.CurriculumClassNumber,
			Code:                  lst.Code,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classLevel_})
}
func FetchClassesToAdd(context *gin.Context) {

	levelNumber := context.Query("level_number")

	var campusLevelData levelModels.CampusLevel
	if err := config.DB.
		Where("number = ?", levelNumber).
		Preload("CurriculumLevel").
		First(&campusLevelData).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	var curriculumClasses []classModel.CurriculumClass
	var curriculumClasses_ []responses.CampusClassInfo
	if err := config.DB.
		Where("curriculum_level_number = ?", campusLevelData.CurriculumLevelNumber).
		Preload("Class").
		Order("code asc").
		Find(&curriculumClasses).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	for _, lst := range curriculumClasses {
		curriculumClasses_ = append(curriculumClasses_, responses.CampusClassInfo{
			Number:       lst.Number,
			Name:         lst.Class.Name,
			ClassNumber:  lst.Class.Number,
			Code:         lst.Code,
			Abbreviation: lst.Class.Abbreviation,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculumClasses_})
}

func FetchClassLevel(context *gin.Context) {

	campusCurriculumNumber := context.Query("campus_curriculum_number")
	campusLevelNumber := context.Query("campus_level_number")

	if campusCurriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	var campusLevel []levelModels.CampusLevel
	var campusLevel_ []responses3.CampusLevelDataList
	if err := config.DB.
		Scopes(levelClassFilter(campusLevelNumber)).
		Where("campus_curriculum_number = ?", campusCurriculumNumber).
		Preload("CurriculumLevel.Level").
		Preload("CurriculumLevel").
		Find(&campusLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	for _, items := range campusLevel {

		var classLevel []classModel.CampusClass
		var classLevel_ []responses.CampusClassInfo
		if err := config.DB.
			Where("campus_level_number = ?", items.Number).
			Preload("CurriculumClass.Class").
			Preload("CurriculumClass").
			Preload("Streams").
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

			classLevel_ = append(classLevel_, responses.CampusClassInfo{
				Number:                lst.Number,
				Name:                  lst.CurriculumClass.Class.Name,
				ClassNumber:           lst.CurriculumClass.Class.Number,
				Abbreviation:          lst.CurriculumClass.Class.Abbreviation,
				CurriculumClassNumber: lst.CurriculumClassNumber,
				Code:                  lst.Code,
				Streams:               stream_,
			})
		}

		campusLevel_ = append(campusLevel_, responses3.CampusLevelDataList{
			Number:                items.Number,
			CurriculumLevelNumber: items.CurriculumLevelNumber,
			Code:                  items.Code,
			Name:                  items.CurriculumLevel.Level.Name,
			Classes:               classLevel_,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campusLevel_})
}

func levelClassFilter(levelNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if levelNumber != "" {
			return db.Where("number = ?", levelNumber)
		}
		return db
	}
}

func FetchClassSubjects(context *gin.Context) {

	//levelNumber := context.Query("level_number")
	classNumber := context.Query("class_number")

	/*	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}*/

	var classSubjects []classModel.ClassLevelSubject
	var subjects_ []classModel.SubjectInfo
	if err := config.DB.
		Scopes(classFilter(classNumber)).
		//Where("campus_classes.level_number = ?", levelNumber).
		Joins("join campus_classes on campus_classes.number = class_level_subjects.class_level_number").
		Preload("CampusClass.CurriculumClass.Class").
		Preload("CampusClass.CurriculumClass").
		Preload("CampusClass").
		Preload("Subject").
		Order("created_at desc").
		Find(&classSubjects).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject!"})
		return
	}

	for _, item := range classSubjects {
		subjects_ = append(subjects_, classModel.SubjectInfo{
			SubjectNumber: item.Subject.Number,
			Name:          item.Subject.Name,
			Class:         item.ClassLevel.CurriculumClass.Class.Name,
			Abbreviation:  item.ClassLevel.CurriculumClass.Class.Abbreviation,
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

	var duplicateSubjectClass classModel.ClassLevelSubject
	if duplicate := config.DB.
		Where("class_level_number = ?", classNumber).
		Where("subject_number = ?", subjectNumbers).
		First(&duplicateSubjectClass).Error; duplicate != nil {

		//		add a subject to a learning area once no duplicates
		number := uuid.NewV4().String()
		saveLevel := classModel.ClassLevelSubject{
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
			streamObj := classModel.Stream{
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
	var streams []classModel.Stream
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

func FetchClassesByCampusNumber(ctx *gin.Context) {

	campusNumber := ctx.Param("campus_number")

	var campusClasses []classModel.CampusClass

	err := config.DB.
		Where("campus_number = ?", campusNumber).
		Find(&campusClasses).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Classes!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": campusClasses})

}

func DeleteClass(context *gin.Context) {

	var classInfo requests.ClassDelRequest
	err := context.BindJSON(&classInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	classNumber := classInfo.Number

	var classChecker classModel.CurriculumClass
	if err := config.DB.
		Where("class_number = ?", classNumber).
		First(&classChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this class, it has classes under it on curriculum."})
		return
	}

	var classData classModel.Class
	academicYearDataInfo := config.DB.
		Where("number = ?", classNumber).
		First(&classData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Class not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", classNumber).
		Delete(&classData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete level"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete level"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level deleted successfully!"})
}

func ClassInfo(context *gin.Context) {
	classNumber := context.Query("class_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	var class classModel.Class
	if err := config.DB.
		Where("number = ?", classNumber).
		First(&class).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class found"})
		return
	}

	class_ := responses.ClassDataList{
		Number:       class.Number,
		Name:         class.Name,
		Code:         class.Code,
		Abbreviation: class.Abbreviation,
		UserNumber:   class.UserNumber,
		LevelNumber:  class.LevelNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": class_})
}

func FetchClassesBySchoolNumber(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	campusLevelNumber := context.Query("campus_level_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Choose school"})
		return
	}

	classes, err := services.FindSchoolClasses(schoolNumber, campusLevelNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": classes})
}
