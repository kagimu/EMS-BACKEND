package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewTopic(context *gin.Context) {

	type TopicInfo struct {
		Code         string `json:"code"`
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

	code := topicInfo.Code
	name := topicInfo.Name
	strandNumber := topicInfo.StrandNumber
	classNumber := topicInfo.ClassNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
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

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
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

	var duplicateCode models3.Topic
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("strand_number = ?", strandNumber).
		Where("class_number = ?", classNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
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
		Code:         codeNo,
	}
	config.DB.Save(&saveTopic)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Topic Added Successfully!", "data": saveTopic})
}

func FetchTopic(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")
	strandNumber := context.Query("strand_number")
	classNumber := context.Query("class_number")
	curriculumNumber := context.Query("curriculum_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber, curriculumNumber)
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses4.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(lstLA.Number, subjectClusterNumber)
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses3.SubjectClusterInfo
		for _, lstSC := range subjectClusterResp {

			subjectResp, err := curriculum.FetchSubjectCurriculum(lstSC.Number, subjectNumber)
			if err != "" {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			}

			var subjects_ []responses5.SubjectInfoList
			for _, lstSB := range subjectResp {

				strandResp, err := curriculum.FetchStrandCurriculum(lstSB.Number, strandNumber)
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				var strand_ []responses2.StrandDataList
				for _, lstST := range strandResp {

					topicResp, err := curriculum.FetchTopicCurriculum(lstST.Number, "", classNumber)
					if err != "" {
						context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
						return
					}

					var topic_ []responses.TopicDataList
					for _, list := range topicResp {
						topic_ = append(topic_, responses.TopicDataList{
							Number:      list.Number,
							Name:        list.Name,
							Code:        list.Code,
							UserNumber:  list.UserNumber,
							ClassNumber: list.ClassNumber,
						})
					}

					strand_ = append(strand_, responses2.StrandDataList{
						Number:        lstST.Number,
						Name:          lstST.Name,
						Code:          lstST.Code,
						UserNumber:    lstST.UserNumber,
						SubjectNumber: lstST.SubjectNumber,
						Topic:         topic_,
					})
				}

				subjects_ = append(subjects_, responses5.SubjectInfoList{
					Number:               lstSB.Number,
					Name:                 lstSB.Name,
					Code:                 lstSB.Code,
					UserNumber:           lstSB.UserNumber,
					SubjectClusterNumber: lstSB.SubjectClusterNumber,
					Strand:               strand_,
				})

			}

			subjectCluster_ = append(subjectCluster_, responses3.SubjectClusterInfo{
				Number:             lstSC.Number,
				Name:               lstSC.Name,
				Code:               lstSC.Code,
				LearningAreaNumber: lstSC.LearningAreaNumber,
				Subject:            subjects_,
			})

		}

		learningAreas_ = append(learningAreas_, responses4.LearningAreaInfoData{
			Number:         lstLA.Number,
			Name:           lstLA.Name,
			Code:           lstLA.Code,
			UserNumber:     lstLA.UserNumber,
			SubjectCluster: subjectCluster_,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})
}

func FetchTopicList(context *gin.Context) {

	strandNumber := context.Query("strand_number")

	topicResp, err := curriculum.FetchTopicCurriculum(strandNumber, "", "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var topic_ []responses.TopicDataList
	for _, list := range topicResp {
		topic_ = append(topic_, responses.TopicDataList{
			Number:      list.Number,
			Name:        list.Name,
			Code:        list.Code,
			UserNumber:  list.UserNumber,
			ClassNumber: list.ClassNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": topic_})
}

func strand(strandNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if strandNumber != "" {
			return db.Where("strand_number = ?", strandNumber)
		}
		return db
	}
}
