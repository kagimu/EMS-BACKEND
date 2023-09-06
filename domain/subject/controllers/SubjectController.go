package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models3 "github.com/hamsoftug/wholesomeEd/domain/subject/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func NewSubject(context *gin.Context) {

	type SubjectInfo struct {
		Name                 string `json:"name"`
		SubjectClusterNumber string `json:"subject_cluster_number"`
	}

	var subjectInfo SubjectInfo
	err := context.BindJSON(&subjectInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := subjectInfo.Name
	subjectClusterNumber := subjectInfo.SubjectClusterNumber

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
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

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateSubject models3.Subject
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("level_number = ?", subjectClusterNumber).
		First(&duplicateSubject).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveCampus := models3.Subject{
		Number:               number,
		Name:                 name,
		UserNumber:           userNumber,
		SubjectClusterNumber: subjectClusterNumber,
	}

	config.DB.Save(&saveCampus)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
}

func FetchSubject(context *gin.Context) {

	subjectClusterNumber := context.Query("subject_cluster_number")

	/*	if subjectClusterNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject cluster is required!"})
		return
	}*/

	var subjects []models3.Subject
	if err := config.DB.
		Scopes(subjectCluster(subjectClusterNumber)).
		Order("created_at desc").
		Find(&subjects).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjects})
}

func subjectCluster(subjectClusterNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subjectClusterNumber != "" {
			return db.Where("subject_cluster_number = ?", subjectClusterNumber)
		}
		return db
	}
}
