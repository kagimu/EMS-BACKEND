package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	var subjectCluster []models3.SubjectCluster
	var subjectCluster_ []models3.SubjectClusterInfo
	if err := config.DB.
		Scopes(learningArea(learningAreaNumber)).
		Preload("LearningArea").
		Order("code asc").
		Find(&subjectCluster).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject Cluster!"})
		return
	}

	for _, lst := range subjectCluster {
		subjectCluster_ = append(subjectCluster_, models3.SubjectClusterInfo{
			Number:             lst.Number,
			Name:               lst.Name,
			Code:               lst.Code,
			LearningAreaNumber: lst.LearningAreaNumber,
			LearningAreaName:   lst.LearningArea.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjectCluster_})
}

func learningArea(learningAreaNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if learningAreaNumber != "" {
			return db.Where("learning_area_number = ?", learningAreaNumber)
		}
		return db
	}
}
