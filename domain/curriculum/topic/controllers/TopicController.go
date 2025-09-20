package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	subTopicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/responses"
	topicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func NewTopic(context *gin.Context) {

	var topicInfo requests2.TopicInfo
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
	strandNumber := topicInfo.StrandNumber
	classNumber := topicInfo.ClassNumber
	curriculumNumber := topicInfo.CurriculumNumber
	topicNumber := topicInfo.TopicNumber

	userNumber := context.Request.Header["User-Number"][0]

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

	var duplicateCode topicModel.Topic
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("strand_number = ?", strandNumber).
		Where("curriculum_class_number = ?", classNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("topic_number != ?", topicNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateTopic topicModel.Topic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("strand_number = ?", strandNumber).
		Where("curriculum_class_number = ?", classNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("topic_number != ?", topicNumber).
		First(&duplicateTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveTopic := topicModel.Topic{
		Number:                number,
		Name:                  name,
		StrandNumber:          strandNumber,
		UserNumber:            userNumber,
		CurriculumClassNumber: classNumber,
		Code:                  codeNo,
		CurriculumNumber:      curriculumNumber,
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

					topicResp, err := curriculum.FetchTopicCurriculum("", lstST.Number, "", classNumber)
					if err != "" {
						context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
						return
					}

					var topic_ []responses.TopicDataList
					for _, list := range topicResp {
						topic_ = append(topic_, responses.TopicDataList{
							Number:                list.Number,
							Name:                  list.Name,
							Code:                  list.Code,
							UserNumber:            list.UserNumber,
							CurriculumClassNumber: list.CurriculumClassNumber,
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
	curriculumNumber := context.Query("curriculum_number")

	topicResp, err := curriculum.FetchTopicCurriculum(curriculumNumber, strandNumber, "", "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var topic_ []responses.TopicDataList
	for _, list := range topicResp {
		topic_ = append(topic_, responses.TopicDataList{
			Number:     list.Number,
			Name:       list.Name,
			Code:       list.Code,
			UserNumber: list.UserNumber,
			//ClassNumber: list.ClassNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": topic_})
}

func DeleteTopic(context *gin.Context) {
	var topicInfo requests2.TopicDelRequest
	err := context.BindJSON(&topicInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	topicNumber := topicInfo.Number

	var topicChecker subTopicModel.SubTopic
	if err := config.DB.
		Where("topic_number = ?", topicNumber).
		First(&topicChecker).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": http.StatusNotFound, "message": "You can not delete this topic, it has subtopic."})
		return
	}

	var topicData topicModel.Topic
	academicYearDataInfo := config.DB.
		Where("number = ?", topicNumber).
		First(&topicData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Topic not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", topicNumber).
		Delete(&topicData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete topic " + err.Error()})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete topic " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Topic deleted successfully!"})
}

func TopicInfo(context *gin.Context) {
	topicNumber := context.Query("topic_number")

	if topicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic is required!"})
		return
	}

	var topic topicModel.Topic
	if err := config.DB.
		Where("number = ?", topicNumber).
		Preload("Strand.Subject.SubjectCluster").
		Preload("Strand.Subject").
		Preload("Strand").
		First(&topic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No learning area found"})
		return
	}

	topic_ := responses.TopicDataList{
		Number:                topic.Number,
		Name:                  topic.Name,
		Code:                  topic.Code,
		UserNumber:            topic.UserNumber,
		CurriculumClassNumber: topic.CurriculumClassNumber,
		StrandNumber:          topic.StrandNumber,
		SubjectNumber:         topic.Strand.SubjectNumber,
		SubjectClusterNumber:  topic.Strand.Subject.SubjectClusterNumber,
		LearningAreaNumber:    topic.Strand.Subject.SubjectCluster.LearningAreaNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": topic_})
}

func TopicEdit(context *gin.Context) {

	var topicInfo requests2.TopicInfo
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
	strandNumber := topicInfo.StrandNumber
	classNumber := topicInfo.ClassNumber
	curriculumNumber := topicInfo.CurriculumNumber
	topicNumber := topicInfo.TopicNumber

	userNumber := context.Request.Header["User-Number"][0]

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

	var duplicateCode topicModel.Topic
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("strand_number = ?", strandNumber).
		Where("curriculum_class_number = ?", classNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", topicNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateTopic topicModel.Topic
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("strand_number = ?", strandNumber).
		Where("curriculum_class_number = ?", classNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", topicNumber).
		First(&duplicateTopic).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Topic already exists"})
		return
	}

	tx := config.DB.Begin()

	var topicDataInfo topicModel.Topic
	updateData := tx.Model(&topicDataInfo).
		Where("number = ?", topicNumber).
		Updates(topicModel.Topic{Name: name, UserNumber: userNumber, Code: codeNo, StrandNumber: strandNumber, CurriculumClassNumber: classNumber})
	if updateData.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update topic"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Topic Updated Successfully!", "data": topicDataInfo})
}
