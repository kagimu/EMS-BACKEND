package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	responses7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"
	responses6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/responses"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
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

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
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

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")
	strandNumber := context.Query("strand_number")
	topicNumber := context.Query("topic_number")
	subtopicNumber := context.Query("sub_topic_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber)
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

					topicResp, err := curriculum.FetchTopicCurriculum(lstST.Number, topicNumber, "")
					if err != "" {
						context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
						return
					}

					var topic_ []responses6.TopicDataList
					for _, list := range topicResp {

						subtopicResp, err := curriculum.FetchSubtopicCurriculum(list.Number, subtopicNumber)
						if err != "" {
							context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
							return
						}

						var subTopicInfo_ []responses7.SubTopicList
						for idx, subtopicList := range subtopicResp {

							conceptResp, err := curriculum.FetchConceptCurriculum(subtopicList.Number)
							if err != "" {
								context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
								return
							}

							var concept_ []responses.ConceptList
							for _, list := range conceptResp {
								concept_ = append(concept_, responses.ConceptList{
									Number:         list.Number,
									Name:           list.Name,
									Code:           list.Code,
									UserNumber:     list.UserNumber,
									SubTopicNumber: list.SubTopicNumber,
									SubTopicName:   list.SubTopic.Name,
								})
							}

							subTopicInfo_ = append(subTopicInfo_, responses7.SubTopicList{
								Number:        subtopicList.Number,
								Code:          subtopicList.Code,
								SortCode:      int64(idx + 1),
								Name:          subtopicList.Name,
								UserNumber:    subtopicList.UserNumber,
								EstimatedTime: subtopicList.EstimatedTime,
								TopicNumber:   subtopicList.TopicNumber,
								Concept:       concept_,
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
