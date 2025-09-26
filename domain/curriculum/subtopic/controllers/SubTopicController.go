package controllers

import (
	"net/http"
	"strconv"

	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	conceptModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/requests"
	responses7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/api/responses"

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
)

func NewSubTopic(context *gin.Context) {

	var topicInfo requests.SubTopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	topicInfo.Prepare()

	if err := topicInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := topicInfo.Code

	name := topicInfo.Name
	topicNumber := topicInfo.TopicNumber
	estimatedTime := topicInfo.EstimatedTime
	curriculumNumber := topicInfo.CurriculumNumber

	userNumber := context.Request.Header["User-Number"][0]

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	teachingSequenceCode, err := strconv.ParseInt(topicInfo.TeachingSequenceCode, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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

	var duplicateCode models2.SubTopic
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("topic_number = ?", topicNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubTopic models2.SubTopic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("topic_number = ?", topicNumber).
		Where("curriculum_number = ?", curriculumNumber).
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
		Number:              number,
		Name:                name,
		EstimatedTime:       estTime,
		TopicNumber:         topicNumber,
		UserNumber:          userNumber,
		Code:                codeNo,
		CurriculumNumber:    curriculumNumber,
		TeacherSequenceCode: teachingSequenceCode,
	}
	config.DB.Save(&saveSubtopic)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subtopic Added Successfully!", "data": saveSubtopic})
}

func FetchSubTopic(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")
	strandNumber := context.Query("strand_number")
	topicNumber := context.Query("topic_number")
	curriculumNumber := context.Query("curriculum_number")
	classNumber := context.Query("class_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber, curriculumNumber)
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses4.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum("", lstLA.Number, subjectClusterNumber)
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses3.SubjectClusterInfo
		for _, lstSC := range subjectClusterResp {

			subjectResp, err := curriculum.FetchSubjectCurriculum("", lstSC.Number, subjectNumber)
			if err != "" {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			}

			var subjects_ []responses5.SubjectInfoList
			for _, lstSB := range subjectResp {

				strandResp, err := curriculum.FetchStrandCurriculum("", lstSB.Number, strandNumber)
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				var strand_ []responses2.StrandDataList
				for _, lstST := range strandResp {

					topicResp, err := curriculum.FetchTopicCurriculum("", lstST.Number, topicNumber, classNumber)
					if err != "" {
						context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
						return
					}

					var topic_ []responses6.TopicDataList
					for _, list := range topicResp {

						subtopicResp, err := curriculum.FetchSubtopicCurriculum("", list.Number, "")
						if err != "" {
							context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
							return
						}

						var subTopicInfo_ []responses.SubTopicList
						for idx, subtopicList := range subtopicResp {

							referenceResp, err := curriculum.FetchSubTopicReferences(subtopicList.Number)
							if err != "" {
								context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
								return
							}

							var referenceInfo_ []responses7.LearningMaterialList
							for _, refList := range referenceResp {
								referenceInfo_ = append(referenceInfo_, responses7.LearningMaterialList{
									Number:                 refList.Number,
									Code:                   refList.Code,
									Description:            refList.Description,
									LearningMaterialNumber: refList.LearningMaterialNumber,
									SubTopicNumber:         refList.SubTopicNumber,
									UserNumber:             refList.UserNumber,
								})
							}

							subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
								Number:        subtopicList.Number,
								Code:          subtopicList.Code,
								SortCode:      int64(idx + 1),
								Name:          subtopicList.Name,
								UserNumber:    subtopicList.UserNumber,
								EstimatedTime: subtopicList.EstimatedTime,
								TopicNumber:   subtopicList.TopicNumber,
								ReferenceInfo: referenceInfo_,
							})
						}

						topic_ = append(topic_, responses6.TopicDataList{
							Number:     list.Number,
							Name:       list.Name,
							Code:       list.Code,
							UserNumber: list.UserNumber,
							//ClassNumber: list.ClassNumber,
							SubTopic: subTopicInfo_,
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

func FetchSubTopiList(context *gin.Context) {

	topicNumber := context.Query("topic_number")
	curriculumNumber := context.Query("curriculum_number")

	subtopicResp, err := curriculum.FetchSubtopicCurriculum(curriculumNumber, topicNumber, "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var subTopicInfo_ []responses.SubTopicList
	for idx, subtopicList := range subtopicResp {
		subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
			Number:               subtopicList.Number,
			Code:                 subtopicList.Code,
			SortCode:             int64(idx + 1),
			Name:                 subtopicList.Name,
			UserNumber:           subtopicList.UserNumber,
			EstimatedTime:        subtopicList.EstimatedTime,
			TopicNumber:          subtopicList.TopicNumber,
			TeachingSequenceCode: subtopicList.TeacherSequenceCode,
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

func FetchSubTopicByClass(context *gin.Context) {

	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	var campusClassInfo models3.CampusClass
	if err := config.DB.
		Where("number = ?", classNumber).
		Preload("CurriculumClass.Class").
		Preload("CurriculumClass").
		Preload("Streams").
		Order("code asc").
		First(&campusClassInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	var subtopic []models2.SubTopic
	var subTopicInfo_ []responses.SubTopicList
	if err := config.DB.
		Scopes(classFilter(campusClassInfo.CurriculumClassNumber)).
		Joins("join topics on topics.number = sub_topics.topic_number").
		Preload("Topic.Strand.Subject").
		Preload("Topic.Strand").
		Preload("Topic").
		Order("created_at asc").
		Find(&subtopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Sub topic found!"})
		return
	}

	for idx, subtopicList := range subtopic {
		if subjectNumber == subtopicList.Topic.Strand.Subject.Number {
			subTopicInfo_ = append(subTopicInfo_, responses.SubTopicList{
				Number:               subtopicList.Number,
				Code:                 subtopicList.Code,
				SortCode:             int64(idx + 1),
				Name:                 subtopicList.Name,
				UserNumber:           subtopicList.UserNumber,
				EstimatedTime:        subtopicList.EstimatedTime,
				TopicNumber:          subtopicList.TopicNumber,
				TopicName:            subtopicList.Topic.Name,
				StrandNumber:         subtopicList.Topic.Strand.Number,
				StrandName:           subtopicList.Topic.Strand.Name,
				TeachingSequenceCode: subtopicList.TeacherSequenceCode,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subTopicInfo_})
}

func classFilter(classNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if classNumber != "" {
			return db.Where("topics.curriculum_class_number = ?", classNumber)
		}
		return db
	}
}

func DeleteSubTopic(context *gin.Context) {
	var subTopicInfo requests.SubTopicDelRequest
	err := context.BindJSON(&subTopicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	subTopicNumber := subTopicInfo.Number

	var subTopicChecker conceptModel.Concept
	if err := config.DB.
		Where("sub_topic_number = ?", subTopicNumber).
		First(&subTopicChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this sub topic, it has concept."})
		return
	}

	var subTopicData models2.SubTopic
	subTopicDataInfo := config.DB.
		Where("number = ?", subTopicNumber).
		First(&subTopicData)
	if subTopicDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Sub topic not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", subTopicNumber).
		Delete(&subTopicData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete sub topic "+err.Error()})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete topic"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Sub topic deleted successfully!"})
}

func SubTopicInfo(context *gin.Context) {
	subTopicNumber := context.Query("sub_topic_number")

	if subTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
		return
	}

	var subTopic models2.SubTopic
	if err := config.DB.
		Where("number = ?", subTopicNumber).
		Preload("Topic.Strand.Subject.SubjectCluster").
		Preload("Topic.Strand.Subject").
		Preload("Topic.Strand").
		Preload("Topic").
		First(&subTopic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No sub topic found"})
		return
	}

	subTopic_ := responses.SubTopicList{
		Number:               subTopic.Number,
		Code:                 subTopic.Code,
		Name:                 subTopic.Name,
		UserNumber:           subTopic.UserNumber,
		EstimatedTime:        subTopic.EstimatedTime,
		TopicNumber:          subTopic.TopicNumber,
		StrandNumber:         subTopic.Topic.StrandNumber,
		SubjectNumber:        subTopic.Topic.Strand.SubjectNumber,
		SubjectClusterNumber: subTopic.Topic.Strand.Subject.SubjectClusterNumber,
		LearningAreaNumber:   subTopic.Topic.Strand.Subject.SubjectCluster.LearningAreaNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subTopic_})
}

func SubTopicEdit(context *gin.Context) {

	var topicInfo requests.SubTopicInfo
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	topicInfo.Prepare()

	if err := topicInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := topicInfo.Code
	name := topicInfo.Name
	topicNumber := topicInfo.TopicNumber
	estimatedTime := topicInfo.EstimatedTime
	curriculumNumber := topicInfo.CurriculumNumber
	subTopicNumber := topicInfo.SubTopicNumber

	userNumber := context.Request.Header["User-Number"][0]

	if subTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub Topic is required!"})
		return
	}

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	teachingSequenceCode, err := strconv.ParseInt(topicInfo.TeachingSequenceCode, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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

	var duplicateCode models2.SubTopic
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("topic_number = ?", topicNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", subTopicNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubTopic models2.SubTopic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("topic_number = ?", topicNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", subTopicNumber).
		First(&duplicateSubTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic already exists"})
		return
	}

	estTime, _ := strconv.ParseInt(estimatedTime, 10, 64)

	tx := config.DB.Begin()

	var subTopicDataInfo models2.SubTopic
	updateData := tx.Model(&subTopicDataInfo).
		Where("number = ?", subTopicNumber).
		Updates(models2.SubTopic{Name: name, UserNumber: userNumber, Code: codeNo, EstimatedTime: estTime, TopicNumber: topicNumber, TeacherSequenceCode: teachingSequenceCode})
	if updateData.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update sub topic"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Sub topic Edited Successfully!", "data": subTopicDataInfo})
}
