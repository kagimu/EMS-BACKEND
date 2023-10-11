package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewSubTopic(context *gin.Context) {

	type TopicInfo struct {
		Name          string `json:"name"`
		TopicNumber   string `json:"topic_number"`
		EstimatedTime string `json:"estimated_time"`
	}

	var topicInfo TopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := topicInfo.Name
	topicNumber := topicInfo.TopicNumber
	estimatedTime := topicInfo.EstimatedTime

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if topicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic is required!"})
		return
	}

	if estimatedTime == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Estimated time is required!"})
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

	estTime, err := strconv.ParseInt(estimatedTime, 10, 64)

	if err != nil {
		// Handle the error if the string cannot be converted to int64
		//return err
	}

	number := uuid.NewV4().String()
	///save user data
	saveSubtopic := models2.SubTopic{
		Number:        number,
		Name:          name,
		EstimatedTime: estTime,
		TopicNumber:   topicNumber,
		UserNumber:    userNumber,
	}
	config.DB.Save(&saveSubtopic)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subtopic Added Successfully!", "data": saveSubtopic})
}

func FetchSubTopic(context *gin.Context) {

	topicNumber := context.Query("topic_number")

	/*if topicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic is required!"})
		return
	}*/

	var subtopic []models2.SubTopic
	var subTopicInfo_ []responses.SubTopicList
	if err := config.DB.
		Scopes(topic(topicNumber)).
		Order("code asc").
		Find(&subtopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!"})
		return
	}

	for idx, subtopicList := range subtopic {
		subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
			Number:        subtopicList.Number,
			Code:          subtopicList.Code,
			SortCode:      strconv.Itoa(idx + 1),
			Name:          subtopicList.Name,
			UserNumber:    subtopicList.UserNumber,
			EstimatedTime: subtopicList.EstimatedTime,
			TopicNumber:   subtopicList.TopicNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subTopicInfo_})
}

func topic(topicNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if topicNumber != "" {
			return db.Where("topic_number = ?", topicNumber)
		}
		return db
	}
}
