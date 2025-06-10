package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	CurriculumClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/responses"
	curriculumLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var levelService = services.NewLevelService()

func NewLevel(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var levelRequest requests2.LevelInfo
	if err := context.BindJSON(&levelRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	levelRequest.Prepare()

	if err := levelRequest.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveLevel, err := levelService.SaveLevel(user.Number, levelRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level Added Successfully!", "data": saveLevel})
}

func NewCurriculumLevel(context *gin.Context) {

	type LevelInfo struct {
		Code             int64  `json:"code"`
		LevelNumber      string `json:"level_number"`
		CurriculumNumber string `json:"curriculum_number"`
	}

	var levelInfo LevelInfo
	err := context.BindJSON(&levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := levelInfo.Code
	levelNumber := levelInfo.LevelNumber
	curriculumNumber := levelInfo.CurriculumNumber

	userNumber := context.Request.Header["User-Number"][0]

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	if curriculumNumber == "" {
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

	db := config.DB

	var duplicateLevel curriculumLevelModels.CurriculumLevel
	if duplicate := db.
		Where("level_number = ?", levelNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateLevel).Error; duplicate == nil {

		var levelChildrenChecker CurriculumClassModel.CurriculumClass
		if err := db.
			Where("curriculum_level_number = ?", duplicateLevel.Number).
			First(&levelChildrenChecker).Error; err == nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not remove this Level, it has Children"})
			return
		}

		//	delete level on curriculum
		db.Unscoped().Where("curriculum_number = ?", curriculumNumber).
			Where("level_number = ?", levelNumber).
			Delete(&duplicateLevel)

		context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level removed successfully!"})

		//context.JSON(http.StatusBadRequest, gin.H{"message": "Level already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	///save user data
	saveLevel := curriculumLevelModels.CurriculumLevel{
		Number:           uuid.NewV4().String(),
		Code:             code,
		LevelNumber:      levelNumber,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}

	if err := tx.Create(&saveLevel).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		//context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create level"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add level"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level Added Successfully!"})
}

func FetchLearningLevel(context *gin.Context) {

	var levels []curriculumLevelModels.Level
	var levels_ []responses2.LevelList
	if err := config.DB.
		Scopes(filters.ByEntityNumber(context.Query("entity_number"))).
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {
		levels_ = append(levels_, responses2.LevelList{
			Number:     list.Number,
			Name:       list.Name,
			Code:       list.Code,
			UserNumber: list.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levels_})
}

func FetchCurriculumLevel(context *gin.Context) {

	curriculumNumber := context.Query("curriculum_number")
	var curriculumLevels []curriculumLevelModels.CurriculumLevel
	var curriculumLevels_ []responses2.CurriculumLevelList
	if err := config.DB.
		Where("curriculum_number = ?", curriculumNumber).
		Preload("Level").
		Preload("Curriculum").
		Order("code asc").
		Find(&curriculumLevels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level found!"})
		return
	}

	for _, list := range curriculumLevels {
		curriculumLevels_ = append(curriculumLevels_, responses2.CurriculumLevelList{
			Number:           list.Number,
			Name:             list.Level.Name,
			LevelNumber:      list.Level.Number,
			Code:             list.Code,
			UserNumber:       list.UserNumber,
			CurriculumName:   list.Curriculum.Name,
			CurriculumNumber: list.Curriculum.Number,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": curriculumLevels_})
}

func DeleteLevel(context *gin.Context) {

	var levelInfo requests2.LevelDelRequest
	err := context.BindJSON(&levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	levelNumber := levelInfo.Number

	var classChecker CurriculumClassModel.Class
	if err := config.DB.
		Where("level_number = ?", levelNumber).
		First(&classChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Level, it has classes under it."})
		return
	}

	var levelData curriculumLevelModels.Level
	academicYearDataInfo := config.DB.
		Where("number = ?", levelNumber).
		First(&levelData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Level not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", levelNumber).
		Delete(&levelData).Error; err != nil {
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

func LevelInfo(context *gin.Context) {

	levelNumber := context.Query("level_number")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level Number is required!"})
		return
	}

	var level curriculumLevelModels.Level
	if err := config.DB.
		Where("number = ?", levelNumber).
		First(&level).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No level found"})
		return
	}

	level_ := responses2.LevelList{
		Number:     level.Number,
		Name:       level.Name,
		Code:       level.Code,
		UserNumber: level.UserNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": level_})
}

func EditLevel(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var levelInfo requests2.LevelInfo
	err = context.BindJSON(&levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	levelInfo.Prepare()

	if err := levelInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	levelDataInfo, err := levelService.EditLevel(user.Number, levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level updated Successfully!", "data": levelDataInfo})
}
