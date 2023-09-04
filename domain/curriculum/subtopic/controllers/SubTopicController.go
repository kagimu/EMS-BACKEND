package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/curriculum/subtopic/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewSubTopic(context *gin.Context) {

	type TopicInfo struct {
		Name        string `json:"name"`
		TopicNumber string `json:"topic_number"`
	}

	var topicInfo TopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := topicInfo.Name
	topicNumber := topicInfo.TopicNumber

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

	if topicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic is required!"})
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

	var duplicateSubTopic models2.SubTopic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("topic_number = ?", topicNumber).
		First(&duplicateSubTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStrand := models2.SubTopic{
		Number:      number,
		Name:        name,
		TopicNumber: topicNumber,
		UserNumber:  userNumber,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subtopic Added Successfully!", "data": saveStrand})
}

func FetchSubTopic(context *gin.Context) {

	topicNumber := context.Query("topic_number")

	if topicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic is required!"})
		return
	}

	var subtopic []models2.SubTopic
	if err := config.DB.
		Where("topic_number = ?", topicNumber).
		Order("created_at desc").
		Find(&subtopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subtopic})
}
