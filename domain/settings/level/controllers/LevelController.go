package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/responses"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewLevel(context *gin.Context) {

	type LevelInfo struct {
		Code             string `json:"code"`
		Name             string `json:"name"`
		CurriculumNumber string `json:"curriculum_number"`
	}

	var levelInfo LevelInfo
	err := context.BindJSON(&levelInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := levelInfo.Code
	name := levelInfo.Name
	curriculumNumber := levelInfo.CurriculumNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

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

	var duplicateCode levelModels.LearningLevel
	if duplicate := config.DB.
		Where("code = ?", code).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateLevel levelModels.LearningLevel
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateLevel).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()
	number := uuid.NewV4().String()
	///save user data
	saveLevel := levelModels.LearningLevel{
		Number:           number,
		Name:             name,
		Code:             codeNo,
		CurriculumNumber: curriculumNumber,
		UserNumber:       userNumber,
	}

	if err := tx.Create(&saveLevel).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create level"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add level"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Level Added Successfully!", "data": saveLevel})
}

func FetchLearningLevel(context *gin.Context) {

	curriculumNumber := context.Query("curriculum_number")
	var levels []levelModels.LearningLevel
	var levels_ []responses2.LearningLevelList
	if err := config.DB.
		Where("curriculum_number = ?", curriculumNumber).
		Order("code asc").
		Find(&levels).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	for _, list := range levels {
		levels_ = append(levels_, responses2.LearningLevelList{
			Number:     list.Number,
			Name:       list.Name,
			Code:       list.Code,
			UserNumber: list.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": levels_})
}
