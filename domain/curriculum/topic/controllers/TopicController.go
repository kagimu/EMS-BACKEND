package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/curriculum/topic/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewTopic(context *gin.Context) {

	type TopicInfo struct {
		Name         string `json:"name"`
		StrandNumber string `json:"strand_number"`
	}

	var topicInfo TopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := topicInfo.Name
	strandNumber := topicInfo.StrandNumber

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

	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand is required!"})
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

	var duplicateTopic models3.Topic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("strand_number = ?", strandNumber).
		First(&duplicateTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStrand := models3.Topic{
		Number:       number,
		Name:         name,
		StrandNumber: strandNumber,
		UserNumber:   userNumber,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveStrand})
}

func FetchTopic(context *gin.Context) {

	strandNumber := context.Query("strand_number")

	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var topic []models3.Topic
	if err := config.DB.
		Where("strand_number = ?", strandNumber).
		Order("created_at desc").
		Find(&topic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": topic})
}
