package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewSubjectCluster(context *gin.Context) {

	type SubjectClusterInfo struct {
		Code               string `json:"code"`
		Name               string `json:"name"`
		LearningAreaNumber string `json:"learning_area_number"`
	}

	var subjectClusterInfo SubjectClusterInfo
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectClusterInfo.Code
	name := subjectClusterInfo.Name
	learningAreaNumber := subjectClusterInfo.LearningAreaNumber

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

	var duplicateCode models3.SubjectCluster
	if duplicate := config.DB.
		Where("code = ?", codeNo).
		Where("learning_area_number = ?", learningAreaNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateSubjectCluster models3.SubjectCluster
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("learning_area_number = ?", learningAreaNumber).
		First(&duplicateSubjectCluster).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveSubjectCluster := models3.SubjectCluster{
		Number:             number,
		Name:               name,
		LearningAreaNumber: learningAreaNumber,
		UserNumber:         userNumber,
		Code:               codeNo,
	}
	config.DB.Save(&saveSubjectCluster)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveSubjectCluster})
}

func FetchSubjectCluster(context *gin.Context) {

	learningAreaNumber := context.Query("learning_area_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum(learningAreaNumber, "")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses2.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(lstLA.Number, "")
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

	//fetch subject cluster
	subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(learningAreaNumber, "")
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
