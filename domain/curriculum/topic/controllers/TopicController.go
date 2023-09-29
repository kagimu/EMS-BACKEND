package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func NewTopic(context *gin.Context) {

	type TopicInfo struct {
		Name         string `json:"name"`
		StrandNumber string `json:"strand_number"`
		ClassNumber  string `json:"Class_number"`
	}

	var topicInfo TopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := topicInfo.Name
	strandNumber := topicInfo.StrandNumber
	classNumber := topicInfo.ClassNumber

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
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
		Where("class_number = ?", classNumber).
		First(&duplicateTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveTopic := models3.Topic{
		Number:       number,
		Name:         name,
		StrandNumber: strandNumber,
		UserNumber:   userNumber,
		ClassNumber:  classNumber,
	}
	config.DB.Save(&saveTopic)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Topic Added Successfully!", "data": saveTopic})
}

func FetchTopic(context *gin.Context) {

	strandNumber := context.Query("strand_number")

	/*	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}*/

	var topic []models3.Topic
	if err := config.DB.
		Scopes(strand(strandNumber)).
		Order("created_at desc").
		Find(&topic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": topic})
}

func strand(strandNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if strandNumber != "" {
			return db.Where("strand_number = ?", strandNumber)
		}
		return db
	}
}
