package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	topicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewSubject(context *gin.Context) {

	type SubjectInfo struct {
		Code                 string `json:"code"`
		Name                 string `json:"name"`
		SubjectClusterNumber string `json:"subject_cluster_number"`
	}

	var subjectInfo SubjectInfo
	err := context.BindJSON(&subjectInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectInfo.Code
	name := subjectInfo.Name
	subjectClusterNumber := subjectInfo.SubjectClusterNumber

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

	if subjectClusterNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
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

	var duplicateCode subjectModel.Subject
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("subject_cluster_number = ?", subjectClusterNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubject subjectModel.Subject
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_cluster_number = ?", subjectClusterNumber).
		First(&duplicateSubject).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveSubject := subjectModel.Subject{
		Number:               number,
		Name:                 name,
		UserNumber:           userNumber,
		SubjectClusterNumber: subjectClusterNumber,
		Code:                 codeNo,
	}

	config.DB.Save(&saveSubject)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveSubject})
}

func FetchSubject(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber)
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses2.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(lstLA.Number, subjectClusterNumber)
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

			var subjects_ []responses.SubjectInfoList
			for _, item := range subjectResp {
				subjects_ = append(subjects_, responses.SubjectInfoList{
					Number:               item.Number,
					Name:                 item.Name,
					Code:                 item.Code,
					UserNumber:           item.UserNumber,
					SubjectClusterNumber: item.SubjectClusterNumber,
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

		learningAreas_ = append(learningAreas_, responses2.LearningAreaInfoData{
			Number:         lstLA.Number,
			Name:           lstLA.Name,
			Code:           lstLA.Code,
			UserNumber:     lstLA.UserNumber,
			SubjectCluster: subjectCluster_,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})

}

func FetchSubjectList(context *gin.Context) {

	subjectClusterNumber := context.Query("subject_cluster_number")

	subjectResp, err := curriculum.FetchSubjectCurriculum(subjectClusterNumber, "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var subjects_ []responses.SubjectInfoList
	for _, item := range subjectResp {
		subjects_ = append(subjects_, responses.SubjectInfoList{
			Number:               item.Number,
			Name:                 item.Name,
			Code:                 item.Code,
			UserNumber:           item.UserNumber,
			SubjectClusterNumber: item.SubjectClusterNumber,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjects_})
}

type SubjectData struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

func FetchSubjectByClass(context *gin.Context) {

	classNumber := context.Query("class_number")

	var topic []topicModel.Topic
	if err := config.DB.
		Where("topics.class_number = ?", classNumber).
		Joins("join strands on strands.number = topics.strand_number").
		Preload("Strand.Subject").
		Preload("Strand").
		Order("created_at desc").
		Find(&topic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	var subjectData []SubjectData
	for _, item := range topic {
		subjectData = append(subjectData, SubjectData{
			Number: item.Strand.Subject.Number,
			Name:   item.Strand.Subject.Name,
		})
	}

	subjectData = removeDuplicates(subjectData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjectData})
}

func removeDuplicates(data []SubjectData) []SubjectData {

	subjectMap := map[SubjectData]bool{}
	var result []SubjectData

	for _, subject := range data {
		if !subjectMap[subject] {
			subjectMap[subject] = true
			result = append(result, subject)
		}
	}
	return result
}

func FetchSubjectsByClassLevelNumber(ctx *gin.Context) {

	classLevelNumber := ctx.Param("class_level_number")

	var subjects []classModels.ClassLevelSubject
	config.DB.
		Where("class_level_number = ?", classLevelNumber).
		Preload("Subject").
		Find(&subjects)
	ctx.JSON(http.StatusOK, gin.H{"data": subjects})

}
