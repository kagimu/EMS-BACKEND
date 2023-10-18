package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/responses"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewSchemeOfWork(context *gin.Context) {

	type SubtopicListInfo struct {
		Number        string `json:"number"`
		Name          string `json:"name"`
		Code          int64  `json:"code"`
		SortCode      int64  `json:"sort_code"`
		UserNumber    string `json:"user_number"`
		EstimatedTime int64  `json:"estimated_time"`
		TopicNumber   string `json:"topic_number"`
	}

	type SchemeOfWorkInfo struct {
		Name         string             `json:"name"`
		ClassNumber  string             `json:"class_number"`
		SubtopicList []SubtopicListInfo `json:"subtopicList"`
	}

	var schemeOfWorkInfo SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	name := schemeOfWorkInfo.Name
	classNumber := schemeOfWorkInfo.ClassNumber
	subtopicList := schemeOfWorkInfo.SubtopicList
	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if len(subtopicList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic list is required!"})
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

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("class_number = ?", classNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	number := uuid.NewV4().String()
	///save sheme of work data
	saveSchemeOfWork := schemeOfWorkModel.SchemeOfWork{
		Number:      number,
		Name:        name,
		Status:      "Pending",
		ClassNumber: classNumber,
		UserNumber:  userNumber,
	}

	if err := tx.Create(&saveSchemeOfWork).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	for _, schemeList := range subtopicList {
		schemeNo := uuid.NewV4().String()
		saveSchemeOfWorkSubtopic := schemeOfWorkModel.SchemeOfWorkSubTopic{
			Number:             schemeNo,
			Code:               schemeList.SortCode,
			EstimatedTime:      schemeList.EstimatedTime,
			SchemeOfWorkNumber: saveSchemeOfWork.Number,
			SubtopicNumber:     schemeList.Number,
			UserNumber:         userNumber,
		}
		if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
			return
		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme Added Successfully!", "data": saveSchemeOfWork})
}

func FetchSchemeOfWork(context *gin.Context) {

	classNumber := context.Query("class_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	var schemeOfWrkInfo schemeOfWorkModel.SchemeOfWork
	if err := config.DB.
		Where("class_number = ?", classNumber).
		Order("created_at desc").
		First(&schemeOfWrkInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	var schemeOfWrkSubTopicInfo []schemeOfWorkModel.SchemeOfWorkSubTopic
	var schemeOfWrkSubTopicInfo_ []responses2.SchemeOfWorkSubtopicInfo
	if err := config.DB.
		Where("scheme_of_work_number = ?", schemeOfWrkInfo.Number).
		Preload("SubTopic").
		Order("code asc").
		Find(&schemeOfWrkSubTopicInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	for _, subTopicList := range schemeOfWrkSubTopicInfo {
		schemeOfWrkSubTopicInfo_ = append(schemeOfWrkSubTopicInfo_, responses2.SchemeOfWorkSubtopicInfo{
			Number:             subTopicList.Number,
			Code:               subTopicList.Code,
			Name:               subTopicList.SubTopic.Name,
			EstimatedTime:      subTopicList.EstimatedTime,
			SchemeOfWorkNumber: subTopicList.SchemeOfWorkNumber,
			SubtopicNumber:     subTopicList.SubtopicNumber,
			UserNumber:         subTopicList.UserNumber,
		})
	}

	schemeOfWrk_ := responses2.SchemeOfWorkResponseInfo{
		Number:                   schemeOfWrkInfo.Number,
		Name:                     schemeOfWrkInfo.Name,
		Status:                   schemeOfWrkInfo.Status,
		ClassNumber:              schemeOfWrkInfo.ClassNumber,
		UserNumber:               schemeOfWrkInfo.UserNumber,
		SchemeOfWorkSubtopicInfo: schemeOfWrkSubTopicInfo_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schemeOfWrk_})

}
