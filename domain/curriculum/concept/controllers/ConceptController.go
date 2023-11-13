package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewConcept(context *gin.Context) {

	type ConceptInfo struct {
		Code           string `json:"code"`
		Name           string `json:"name"`
		SubTopicNumber string `json:"sub_topic_number"`
	}

	var conceptInfo ConceptInfo
	err := context.BindJSON(&conceptInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := conceptInfo.Code
	name := conceptInfo.Name
	subtopicNumber := conceptInfo.SubTopicNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if subtopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
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

	var duplicateConcept models3.Concept
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("sub_topic_number = ?", subtopicNumber).
		First(&duplicateConcept).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveConcept := models3.Concept{
		Number:         number,
		Name:           name,
		SubTopicNumber: subtopicNumber,
		UserNumber:     userNumber,
		Code:           codeNo,
	}
	config.DB.Save(&saveConcept)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Concept Added Successfully!", "data": saveConcept})
}

func FetchConcept(context *gin.Context) {

	subtopicNumber := context.Query("sub_topic_number")

	var concept []models3.Concept
	var concept_ []responses.ConceptList
	if err := config.DB.
		Scopes(subTopic(subtopicNumber)).
		Preload("SubTopic.Topic.Strand.Subject").
		Preload("SubTopic.Topic.Strand").
		Preload("SubTopic.Topic").
		Preload("SubTopic").
		Order("code asc").
		Find(&concept).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Concept found!"})
		return
	}

	for _, list := range concept {
		concept_ = append(concept_, responses.ConceptList{
			Number:         list.Number,
			Name:           list.Name,
			Code:           list.Code,
			UserNumber:     list.UserNumber,
			SubTopicNumber: list.SubTopicNumber,
			SubTopicName:   list.SubTopic.Name,
			TopicNumber:    list.SubTopic.Topic.Number,
			TopicName:      list.SubTopic.Topic.Name,
			StrandNumber:   list.SubTopic.Topic.Strand.Number,
			StrandName:     list.SubTopic.Topic.Strand.Name,
			SubjectNumber:  list.SubTopic.Topic.Strand.Subject.Number,
			SubjectName:    list.SubTopic.Topic.Strand.Subject.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": concept_})
}

func FetchConceptBySubTopic(context *gin.Context) {

	subtopicNumber := context.Query("sub_topic_number")

	/*	if subtopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
		return
	}*/

	var concept []models3.Concept
	var concept_ []responses.ConceptList
	if err := config.DB.
		Scopes(subTopic(subtopicNumber)).
		Order("code asc").
		Find(&concept).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Concept found!"})
		return
	}

	for _, list := range concept {
		concept_ = append(concept_, responses.ConceptList{
			Number:         list.Number,
			Name:           list.Name,
			Code:           list.Code,
			UserNumber:     list.UserNumber,
			SubTopicNumber: list.SubTopicNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": concept_})
}

func subTopic(subtopicNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subtopicNumber != "" {
			return db.Where("sub_topic_number = ?", subtopicNumber)
		}
		return db
	}
}
