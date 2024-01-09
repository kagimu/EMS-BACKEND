package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
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
