package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	responses6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/responses"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewSubTopic(context *gin.Context) {

	type TopicInfo struct {
		Code          string `json:"code"`
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

	code := topicInfo.Code
	name := topicInfo.Name
	topicNumber := topicInfo.TopicNumber
	estimatedTime := topicInfo.EstimatedTime

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
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

	if estimatedTime == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Estimated time is required!"})
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
		Code:          codeNo,
	}
	config.DB.Save(&saveSubtopic)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subtopic Added Successfully!", "data": saveSubtopic})
}

func FetchSubTopic(context *gin.Context) {

	//learningAreaNumber := context.Query("learning_area_number")
	//subjectClusterNumber := context.Query("subject_cluster_number")
	//subjectNumber := context.Query("subject_number")
	//strandNumber := context.Query("strand_number")
	topicNumber := context.Query("topic_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum("")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses4.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(lstLA.Number, "")
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses3.SubjectClusterInfo
		for _, lstSC := range subjectClusterResp {

			subjectResp, err := curriculum.FetchSubjectCurriculum(lstSC.Number, "")
			if err != "" {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			}

			var subjects_ []responses5.SubjectInfoList
			for _, lstSB := range subjectResp {

				strandResp, err := curriculum.FetchStrandCurriculum(lstSB.Number, "")
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				var strand_ []responses2.StrandDataList
				for _, lstST := range strandResp {

					topicResp, err := curriculum.FetchTopicCurriculum(lstST.Number, topicNumber)
					if err != "" {
						context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
						return
					}

					var topic_ []responses6.TopicDataList
					for _, list := range topicResp {

						subtopicResp, err := curriculum.FetchSubtopicCurriculum(list.Number)
						if err != "" {
							context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
							return
						}

						var subTopicInfo_ []responses.SubTopicList
						for idx, subtopicList := range subtopicResp {
							subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
								Number:        subtopicList.Number,
								Code:          subtopicList.Code,
								SortCode:      int64(idx + 1),
								Name:          subtopicList.Name,
								UserNumber:    subtopicList.UserNumber,
								EstimatedTime: subtopicList.EstimatedTime,
								TopicNumber:   subtopicList.TopicNumber,
							})
						}

						topic_ = append(topic_, responses6.TopicDataList{
							Number:      list.Number,
							Name:        list.Name,
							Code:        list.Code,
							UserNumber:  list.UserNumber,
							ClassNumber: list.ClassNumber,
							SubTopic:    subTopicInfo_,
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

func topic(topicNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if topicNumber != "" {
			return db.Where("topic_number = ?", topicNumber)
		}
		return db
	}
}

func FetchSubTopicByClass(context *gin.Context) {

	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	var subtopic []models2.SubTopic
	var subTopicInfo_ []responses.SubTopicList
	if err := config.DB.
		Scopes(classFilter(classNumber)).
		//Scopes(subjectFilter(subjectNumber)).
		Joins("join topics on topics.number = sub_topics.topic_number").
		Preload("Topic.Strand.Subject").
		Preload("Topic.Strand").
		Preload("Topic").
		Order("code asc").
		Find(&subtopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!"})
		return
	}

	for idx, subtopicList := range subtopic {
		if subjectNumber == subtopicList.Topic.Strand.Subject.Number {
			subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
				Number:        subtopicList.Number,
				Code:          subtopicList.Code,
				SortCode:      int64(idx + 1),
				Name:          subtopicList.Name,
				UserNumber:    subtopicList.UserNumber,
				EstimatedTime: subtopicList.EstimatedTime,
				TopicNumber:   subtopicList.TopicNumber,
				TopicName:     subtopicList.Topic.Name,
				StrandNumber:  subtopicList.Topic.Strand.Number,
				StrandName:    subtopicList.Topic.Strand.Name,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subTopicInfo_})
}

func classFilter(classNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if classNumber != "" {
			return db.Where("topics.class_number = ?", classNumber)
		}
		return db
	}
}

/*func subjectFilter(subjectNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subjectNumber != "" {
			//return db.Where("topics.class_number = ?", subjectNumber)
		}
		return db
	}
}*/
