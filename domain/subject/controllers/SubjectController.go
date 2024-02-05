package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	topicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/requests"
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
		CurriculumNumber     string `json:"curriculum_number"`
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
	curriculumNumber := subjectInfo.CurriculumNumber

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

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
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
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubject subjectModel.Subject
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_cluster_number = ?", subjectClusterNumber).
		Where("curriculum_number = ?", curriculumNumber).
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
		CurriculumNumber:     curriculumNumber,
	}

	config.DB.Save(&saveSubject)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveSubject})
}

func FetchSubject(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	curriculumNumber := context.Query("curriculum_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber, curriculumNumber)
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses2.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum("", lstLA.Number, subjectClusterNumber)
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses3.SubjectClusterInfo
		for _, lstSC := range subjectClusterResp {

			subjectResp, err := curriculum.FetchSubjectCurriculum("", lstSC.Number, "")
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
	curriculumNumber := context.Query("curriculum_number")

	subjectResp, err := curriculum.FetchSubjectCurriculum(curriculumNumber, subjectClusterNumber, "")
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
	Code   int64  `json:"code"`
}

func FetchSubjectByClass(context *gin.Context) {

	classNumber := context.Query("class_number")

	var classInfo classModels.CampusClass

	if err := config.DB.
		Where("number = ?", classNumber).
		Preload("CurriculumClass.Class").
		Preload("CurriculumClass").
		Order("code asc").
		Find(&classInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	var topic []topicModel.Topic
	if err := config.DB.
		Where("topics.curriculum_class_number = ?", classInfo.CurriculumClassNumber).
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
			Code:   item.Strand.Subject.Code,
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

func DeleteSubject(context *gin.Context) {

	var subjectInfo requests2.SubjectDelRequest
	err := context.BindJSON(&subjectInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	subjectNumber := subjectInfo.Number

	var subjectChecker models3.Strand
	if err := config.DB.
		Where("subject_number = ?", subjectNumber).
		First(&subjectChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Subject, it has strand."})
		return
	}

	var subjectData subjectModel.Subject
	academicYearDataInfo := config.DB.
		Where("number = ?", subjectNumber).
		First(&subjectData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Subject not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", subjectNumber).
		Delete(&subjectData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete Subject"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete Subject"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject deleted successfully!"})
}
