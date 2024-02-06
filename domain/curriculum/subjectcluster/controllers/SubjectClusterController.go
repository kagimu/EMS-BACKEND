package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	subjectClusterModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"html"
	"net/http"
	"strconv"
)

func NewSubjectCluster(context *gin.Context) {

	var subjectClusterInfo requests2.SubjectClusterData
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectClusterInfo.Code
	name := subjectClusterInfo.Name
	learningAreaNumber := subjectClusterInfo.LearningAreaNumber
	curriculumNumber := subjectClusterInfo.CurriculumNumber

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if learningAreaNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area is required!"})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
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

	var duplicateCode subjectClusterModel.SubjectCluster
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("learning_area_number = ?", learningAreaNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubjectCluster subjectClusterModel.SubjectCluster
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("learning_area_number = ?", learningAreaNumber).
		Where("curriculum_number = ?", curriculumNumber).
		First(&duplicateSubjectCluster).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveSubjectCluster := subjectClusterModel.SubjectCluster{
		Number:             number,
		Name:               name,
		LearningAreaNumber: learningAreaNumber,
		UserNumber:         userNumber,
		Code:               codeNo,
		CurriculumNumber:   curriculumNumber,
	}
	config.DB.Save(&saveSubjectCluster)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveSubjectCluster})
}

func FetchSubjectCluster(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
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
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(curriculumNumber, lstLA.Number, "")
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses.SubjectClusterInfo
		//if len(subjectClusterResp) > 0 {
		for _, lstSC := range subjectClusterResp {
			subjectCluster_ = append(subjectCluster_, responses.SubjectClusterInfo{
				Number:             lstSC.Number,
				Name:               lstSC.Name,
				Code:               lstSC.Code,
				LearningAreaNumber: lstSC.LearningAreaNumber,
			})
		}

		learningAreas_ = append(learningAreas_, responses2.LearningAreaInfoData{
			Number:         lstLA.Number,
			Name:           lstLA.Name,
			Code:           lstLA.Code,
			UserNumber:     lstLA.UserNumber,
			SubjectCluster: subjectCluster_,
		})
		//}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})

}

func FetchSubjectClusterList(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")
	curriculumNumber := context.Query("curriculum_number")

	//fetch subject cluster
	subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(curriculumNumber, learningAreaNumber, "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var subjectCluster_ []responses.SubjectClusterInfo
	for _, lstSC := range subjectClusterResp {
		subjectCluster_ = append(subjectCluster_, responses.SubjectClusterInfo{
			Number:             lstSC.Number,
			Name:               lstSC.Name,
			Code:               lstSC.Code,
			LearningAreaNumber: lstSC.LearningAreaNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjectCluster_})
}

func DeleteSubjectCluster(context *gin.Context) {
	var subjectClusterInfo requests2.SubjectClusterDelRequest
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	subjectClusterNumber := subjectClusterInfo.Number

	var subjectClusterChecker subjectModel.Subject
	if err := config.DB.
		Where("subject_cluster_number = ?", subjectClusterNumber).
		First(&subjectClusterChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Subject cluster, it has subjects."})
		return
	}

	var subjectClusterData subjectClusterModel.SubjectCluster
	academicYearDataInfo := config.DB.
		Where("number = ?", subjectClusterNumber).
		First(&subjectClusterData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Subject Cluster not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", subjectClusterNumber).
		Delete(&subjectClusterData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete subject cluster"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete subject cluster"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster deleted successfully!"})
}

func SubjectClusterInfo(context *gin.Context) {
	subjectClusterNumber := context.Query("subject_cluster_number")

	if subjectClusterNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject Cluster is required!"})
		return
	}

	var subjectCluster subjectClusterModel.SubjectCluster
	if err := config.DB.
		Where("number = ?", subjectClusterNumber).
		First(&subjectCluster).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No learning area found"})
		return
	}

	subjectCluster_ := responses.SubjectClusterInfo{
		Number:             subjectCluster.Number,
		Name:               html.UnescapeString(subjectCluster.Name),
		Code:               subjectCluster.Code,
		LearningAreaNumber: subjectCluster.LearningAreaNumber,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjectCluster_})
}

func SubjectClusterEdit(context *gin.Context) {

	var subjectClusterInfo requests2.SubjectClusterData
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectClusterInfo.Code
	name := subjectClusterInfo.Name
	learningAreaNumber := subjectClusterInfo.LearningAreaNumber
	curriculumNumber := subjectClusterInfo.CurriculumNumber
	subjectClusterNumber := subjectClusterInfo.SubjectClusterNumber

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if learningAreaNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area is required!"})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
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

	var duplicateCode subjectClusterModel.SubjectCluster
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("learning_area_number = ?", learningAreaNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", subjectClusterNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubjectCluster subjectClusterModel.SubjectCluster
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("learning_area_number = ?", learningAreaNumber).
		Where("curriculum_number = ?", curriculumNumber).
		Where("number != ?", subjectClusterNumber).
		First(&duplicateSubjectCluster).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists"})
		return
	}

	///save user data

	tx := config.DB.Begin()

	var subjectClusterDataInfo subjectClusterModel.SubjectCluster
	updateData := tx.Model(&subjectClusterDataInfo).
		Where("number = ?", subjectClusterNumber).
		Updates(subjectClusterModel.SubjectCluster{Name: name, UserNumber: userNumber, Code: codeNo, LearningAreaNumber: learningAreaNumber})
	if updateData.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create learning area"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Edited Successfully!", "data": subjectClusterDataInfo})
}
