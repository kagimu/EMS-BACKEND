package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/curriculum/concept/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewConcept(context *gin.Context) {

	type ConceptInfo struct {
		Name           string `json:"name"`
		SubTopicNumber string `json:"sub_topic_number"`
	}

	var conceptInfo ConceptInfo
	err := context.BindJSON(&conceptInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := conceptInfo.Name
	subtopicNumber := conceptInfo.SubTopicNumber

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

	if subtopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
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
	}
	config.DB.Save(&saveConcept)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Concept Added Successfully!", "data": saveConcept})
}

func FetchConcept(context *gin.Context) {

	subtopicNumber := context.Query("sub_topic_number")

	if subtopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
		return
	}

	var concept []models3.Concept
	if err := config.DB.
		Where("sub_topic_number = ?", subtopicNumber).
		Order("created_at desc").
		Find(&concept).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Concept found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": concept})
}
