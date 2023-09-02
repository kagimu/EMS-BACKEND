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
		Name string `json:"name"`
	}

	var learningAreaInfo LearningAreaInfo
	err := context.BindJSON(&learningAreaInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := learningAreaInfo.Name

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
		First(&duplicateLearningArea).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveLearningArea := models2.LearningArea{
		Number:     number,
		Name:       name,
		UserNumber: userNumber,
	}
	config.DB.Save(&saveLearningArea)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning Area Added Successfully!", "data": saveLearningArea})
}

func FetchLearningArea(context *gin.Context) {

	var learningAreas []models2.LearningArea
	if err := config.DB.
		Order("created_at desc").
		Find(&learningAreas).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Learning Area!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas})
}
