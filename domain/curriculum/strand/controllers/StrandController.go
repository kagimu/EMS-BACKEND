package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	strandModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	topicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewStrand(context *gin.Context) {

	var strandInfo requests.SubjectClusterInfo
	err := context.BindJSON(&strandInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	strandInfo.Prepare()

	if err := strandInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := strandInfo.Code
	name := strandInfo.Name
	subjectNumber := strandInfo.SubjectNumber
	curriculumNumber := strandInfo.CurriculumNumber

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

	var duplicateCode strandModel.Strand
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("subject_number = ?", subjectNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateStrand strandModel.Strand
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_number = ?", subjectNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateStrand).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStrand := strandModel.Strand{
		Number:           number,
		Name:             name,
		SubjectNumber:    subjectNumber,
		UserNumber:       userNumber,
		Code:             codeNo,
		CurriculumNumber: curriculumNumber,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveStrand})
}

func FetchStrand(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")
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

			var subjects_ []responses.SubjectInfoList
			for _, lstSB := range subjectResp {

				strandResp, err := curriculum.FetchStrandCurriculum("", lstSB.Number, "")
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				var strand_ []responses2.StrandDataList
				for _, lstST := range strandResp {
					strand_ = append(strand_, responses2.StrandDataList{
						Number:        lstST.Number,
						Name:          lstST.Name,
						Code:          lstST.Code,
						UserNumber:    lstST.UserNumber,
						SubjectNumber: lstST.SubjectNumber,
					})
				}

				subjects_ = append(subjects_, responses.SubjectInfoList{
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

func FetchStrandList(context *gin.Context) {

	subjectNumber := context.Query("subject_number")
	curriculumNumber := context.Query("curriculum_number")

	strandResp, err := curriculum.FetchStrandCurriculum(curriculumNumber, subjectNumber, "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var strand_ []responses2.StrandDataList
	for _, lstST := range strandResp {
		strand_ = append(strand_, responses2.StrandDataList{
			Number:        lstST.Number,
			Name:          lstST.Name,
			Code:          lstST.Code,
			UserNumber:    lstST.UserNumber,
			SubjectNumber: lstST.SubjectNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": strand_})
}

func DeleteStrand(context *gin.Context) {

	var strandInfo requests.StrandDelRequest
	err := context.BindJSON(&strandInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	strandNumber := strandInfo.Number

	var strandChecker topicModel.Topic
	if err := config.DB.
		Where("strand_number = ?", strandNumber).
		First(&strandChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this strand, it has topic."})
		return
	}

	var strandData strandModel.Strand
	academicYearDataInfo := config.DB.
		Where("number = ?", strandNumber).
		First(&strandData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Strand not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", strandNumber).
		Delete(&strandData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete strand"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete strand"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Strand deleted successfully!"})
}

func StrandInfo(context *gin.Context) {
	strandNumber := context.Query("strand_number")

	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand is required!"})
		return
	}

	var strand strandModel.Strand
	if err := config.DB.
		Where("number = ?", strandNumber).
		Preload("Subject.SubjectCluster").
		Preload("Subject").
		First(&strand).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No learning area found"})
		return
	}

	strand_ := responses2.StrandDataList{
		Number:                   strand.Number,
		Name:                     strand.Name,
		Code:                     strand.Code,
		UserNumber:               strand.UserNumber,
		SubjectNumber:            strand.SubjectNumber,
		SubjectClusterNumber:     strand.Subject.SubjectClusterNumber,
		LearningAreaNumberNumber: strand.Subject.SubjectCluster.LearningAreaNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": strand_})
}

func StrandEdit(context *gin.Context) {

	var strandInfo requests.SubjectClusterInfo
	err := context.BindJSON(&strandInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	strandInfo.Prepare()

	if err := strandInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	code := strandInfo.Code
	name := strandInfo.Name
	subjectNumber := strandInfo.SubjectNumber
	curriculumNumber := strandInfo.CurriculumNumber
	strandNumber := strandInfo.StrandNumber

	userNumber := context.Request.Header["User-Number"][0]

	if strandNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand is required!"})
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

	var duplicateCode strandModel.Strand
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("subject_number = ?", subjectNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", strandNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateStrand strandModel.Strand
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_number = ?", subjectNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", strandNumber).
		First(&duplicateStrand).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	tx := config.DB.Begin()

	var strandDataInfo strandModel.Strand
	updateData := tx.Model(&strandDataInfo).
		Where("number = ?", strandNumber).
		Updates(strandModel.Strand{Name: name, UserNumber: userNumber, Code: codeNo, SubjectNumber: subjectNumber})
	if updateData.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update strand"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Strand Added Successfully!", "data": strandDataInfo})
}
